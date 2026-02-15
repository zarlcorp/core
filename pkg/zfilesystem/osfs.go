package zfilesystem

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var _ ReadWriteFileFS = (*OSFileSystem)(nil)

// OSFileSystem implements ReadWriteFileFS using the standard os package.
type OSFileSystem struct {
	baseDir string
}

// NewOSFileSystem creates a new OS filesystem with the specified base directory.
// The base directory will be used as the root for all file operations.
func NewOSFileSystem(baseDir string) *OSFileSystem {
	abs, err := filepath.Abs(baseDir)
	if err != nil {
		abs = filepath.Clean(baseDir)
	}
	return &OSFileSystem{baseDir: abs}
}

// resolvePath validates that path stays within baseDir and returns the resolved
// absolute path. It rejects absolute paths and any path that resolves outside
// the base directory via .. traversal.
func (o *OSFileSystem) resolvePath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return "", fmt.Errorf("path escapes base directory: %s", path)
	}
	joined := filepath.Join(o.baseDir, path)
	resolved, err := filepath.Abs(joined)
	if err != nil {
		return "", fmt.Errorf("resolve path: %w", err)
	}
	// ensure the resolved path is the base dir itself or a child of it
	if resolved != o.baseDir && !strings.HasPrefix(resolved, o.baseDir+string(filepath.Separator)) {
		return "", fmt.Errorf("path escapes base directory: %s", path)
	}
	return resolved, nil
}

// ReadFile reads a file from the OS filesystem.
func (o *OSFileSystem) ReadFile(filename string) ([]byte, error) {
	p, err := o.resolvePath(filename)
	if err != nil {
		return nil, err
	}
	return os.ReadFile(p)
}

// WriteFile writes data to a file in the OS filesystem.
func (o *OSFileSystem) WriteFile(filename string, data []byte, perm fs.FileMode) error {
	p, err := o.resolvePath(filename)
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, perm)
}

// Remove removes a file from the OS filesystem.
func (o *OSFileSystem) Remove(filename string) error {
	p, err := o.resolvePath(filename)
	if err != nil {
		return err
	}
	return os.Remove(p)
}

// MkdirAll creates a directory and all necessary parents.
func (o *OSFileSystem) MkdirAll(path string, perm fs.FileMode) error {
	p, err := o.resolvePath(path)
	if err != nil {
		return err
	}
	return os.MkdirAll(p, perm)
}

// OpenFile opens the named file with specified flag and perm.
func (o *OSFileSystem) OpenFile(name string, flag int, perm fs.FileMode) (File, error) {
	p, err := o.resolvePath(name)
	if err != nil {
		return nil, err
	}
	return os.OpenFile(p, flag, perm)
}

// WalkDir walks the directory tree in the OS filesystem.
// Paths passed to the callback are relative to the filesystem's base directory.
func (o *OSFileSystem) WalkDir(root string, fn fs.WalkDirFunc) error {
	p, err := o.resolvePath(root)
	if err != nil {
		return err
	}
	return filepath.WalkDir(p, func(path string, d fs.DirEntry, err error) error {
		rel, relErr := filepath.Rel(o.baseDir, path)
		if relErr != nil {
			return relErr
		}
		return fn(rel, d, err)
	})
}

// BaseDir returns the base directory for this filesystem.
func (o *OSFileSystem) BaseDir() string {
	return o.baseDir
}
