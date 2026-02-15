package zcache

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"sync"

	"github.com/zarlcorp/core/pkg/zfilesystem"

	"github.com/zarlcorp/core/pkg/zoptions"
)

var (
	_ Reader[string, any]     = (*FileCache[string, any])(nil)
	_ Writer[string, any]     = (*FileCache[string, any])(nil)
	_ ReadWriter[string, any] = (*FileCache[string, any])(nil)
	_ Cache[string, any]      = (*FileCache[string, any])(nil)
)

// FileSystem defines the interface that FileCache needs from filesystem implementations.
// This follows the principle that consumers should define interfaces they depend on.
type FileSystem interface {
	// ReadFile reads the file named by filename and returns the contents.
	ReadFile(filename string) ([]byte, error)

	// WriteFile writes data to a file named by filename.
	WriteFile(filename string, data []byte, perm fs.FileMode) error

	// Remove removes the named file.
	Remove(filename string) error

	// WalkDir walks the file tree rooted at root, calling fn for each file or
	// directory in the tree, including root.
	WalkDir(root string, fn fs.WalkDirFunc) error
}

// FileCache is a thread-safe cache implementation using the file system as storage.
// It provides persistent caching capabilities across application restarts when using
// OS filesystem, or in-memory caching for testing when using MemFS (default).
type FileCache[K comparable, V any] struct {
	mu sync.RWMutex
	fs FileSystem
}

// WithFileSystem sets the filesystem implementation for the cache.
// If not provided, an in-memory filesystem (MemFS) is used by default.
func WithFileSystem[K comparable, V any](fsys FileSystem) zoptions.Option[FileCache[K, V]] {
	return func(fc *FileCache[K, V]) {
		fc.fs = fsys
	}
}

// WithOSFileSystem configures the cache to use the OS filesystem with the specified base directory.
// The directory will be created if it doesn't exist.
func WithOSFileSystem[K comparable, V any](baseDir string) zoptions.Option[FileCache[K, V]] {
	return func(fc *FileCache[K, V]) {
		if baseDir == "" {
			tmpDir, err := os.MkdirTemp("", "filecache-*")
			if err != nil {
				// fallback to memory when temp dir fails
				fc.fs = zfilesystem.NewMemFS()
				return
			}
			baseDir = tmpDir
		}

		if err := os.MkdirAll(baseDir, 0o750); err != nil {
			// fallback to memory when mkdir fails
			fc.fs = zfilesystem.NewMemFS()
			return
		}

		fc.fs = zfilesystem.NewOSFileSystem(baseDir)
	}
}

// NewFileCache creates a new file-based cache with optional configuration.
// By default, it uses the OS filesystem with a temporary directory for persistent storage.
// Use WithFileSystem(zfilesystem.NewMemFS()) option for in-memory testing.
func NewFileCache[K comparable, V any](opts ...zoptions.Option[FileCache[K, V]]) *FileCache[K, V] {
	tmpDir, err := os.MkdirTemp("", "filecache-*")
	if err != nil {
		// fallback to memory when temp dir fails
		fc := &FileCache[K, V]{fs: zfilesystem.NewMemFS()}
		for _, opt := range opts {
			opt(fc)
		}
		return fc
	}

	fc := &FileCache[K, V]{
		fs: zfilesystem.NewOSFileSystem(tmpDir),
	}

	for _, opt := range opts {
		opt(fc)
	}

	return fc
}

// Set stores a key-value pair as a file on disk.
// If the key already exists, its value is updated.
func (c *FileCache[K, V]) Set(ctx context.Context, key K, value V) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	filename, err := c.makeFilename(key)
	if err != nil {
		return err
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.fs.WriteFile(filename, data, 0o644)
}

// Get retrieves the value associated with the given key from disk.
// Returns ErrNotFound if the key does not exist.
func (c *FileCache[K, V]) Get(ctx context.Context, key K) (V, error) {
	select {
	case <-ctx.Done():
		var zero V
		return zero, ctx.Err()
	default:
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	filename, err := c.makeFilename(key)
	if err != nil {
		var zero V
		return zero, err
	}

	data, err := c.fs.ReadFile(filename)
	if err != nil {
		var zero V
		if os.IsNotExist(err) {
			return zero, ErrNotFound
		}
		return zero, err
	}

	var value V
	if err := json.Unmarshal(data, &value); err != nil {
		var zero V
		return zero, err
	}

	return value, nil
}

// Delete removes a key-value pair from disk.
// Returns true if the key existed and was deleted, false otherwise.
func (c *FileCache[K, V]) Delete(ctx context.Context, key K) (bool, error) {
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	filename, err := c.makeFilename(key)
	if err != nil {
		return false, err
	}

	err = c.fs.Remove(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// Clear removes all cache files from the base directory.
func (c *FileCache[K, V]) Clear(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	return c.fs.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(d.Name(), ".cache") {
			err = c.fs.Remove(path)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// Len returns the number of cache files in the base directory.
func (c *FileCache[K, V]) Len(ctx context.Context) (int, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	count := 0
	err := c.fs.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(d.Name(), ".cache") {
			count++
		}

		return nil
	})

	return count, err
}

func (c *FileCache[K, V]) makeFilename(key K) (string, error) {
	keyBytes, err := json.Marshal(key)
	if err != nil {
		return "", fmt.Errorf("marshal cache key: %w", err)
	}
	return hex.EncodeToString(keyBytes) + ".cache", nil
}

// Healthy checks if the filesystem is accessible by testing read/write operations.
func (c *FileCache[K, V]) Healthy() error {
	testFile := ".health_check"
	testData := []byte("health_check")

	if err := c.fs.WriteFile(testFile, testData, 0o644); err != nil {
		return fmt.Errorf("write health check: %w", err)
	}

	data, err := c.fs.ReadFile(testFile)
	if err != nil {
		return fmt.Errorf("read health check: %w", err)
	}

	if !bytes.Equal(data, testData) {
		return fmt.Errorf("health check data mismatch")
	}

	_ = c.fs.Remove(testFile)

	return nil
}
