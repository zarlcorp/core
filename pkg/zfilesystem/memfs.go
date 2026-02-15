package zfilesystem

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/zarlcorp/core/pkg/zsync"
)

var _ ReadWriteFileFS = (*MemFS)(nil)

// MemFS provides an in-memory filesystem implementation.
// It implements ReadWriteFileFS and is safe for concurrent use.
type MemFS struct {
	files *zsync.ZMap[string, *memFile]
	dirs  *zsync.ZSet[string]
}

type memFile struct {
	data    []byte
	modTime time.Time
	mode    fs.FileMode
}

// NewMemFS creates a new in-memory filesystem.
func NewMemFS() *MemFS {
	return &MemFS{
		files: zsync.NewZMap[string, *memFile](),
		dirs:  zsync.NewZSet[string](),
	}
}

// cleanPath validates and cleans a path for MemFS. It rejects absolute paths
// and any path that resolves outside the logical root via .. traversal.
func cleanPath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return "", fmt.Errorf("path escapes base directory: %s", path)
	}
	cleaned := filepath.Clean(path)
	if cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path escapes base directory: %s", path)
	}
	return cleaned, nil
}

// ReadFile reads the file from memory.
func (mfs *MemFS) ReadFile(filename string) ([]byte, error) {
	p, err := cleanPath(filename)
	if err != nil {
		return nil, err
	}
	file, ok := mfs.files.Get(p)
	if !ok {
		return nil, os.ErrNotExist
	}

	return bytes.Clone(file.data), nil
}

// WriteFile writes data to memory.
func (mfs *MemFS) WriteFile(filename string, data []byte, perm fs.FileMode) error {
	p, err := cleanPath(filename)
	if err != nil {
		return err
	}
	mfs.files.Set(p, &memFile{
		data:    bytes.Clone(data),
		modTime: time.Now(),
		mode:    perm,
	})
	mfs.addParentDirs(p)

	return nil
}

// addParentDirs records all parent directories for a path.
func (mfs *MemFS) addParentDirs(p string) {
	for {
		parent := filepath.Dir(p)
		if parent == "." || parent == p {
			return
		}
		mfs.dirs.Add(parent)
		p = parent
	}
}

// Remove removes a file from memory.
func (mfs *MemFS) Remove(filename string) error {
	p, err := cleanPath(filename)
	if err != nil {
		return err
	}
	if !mfs.files.Delete(p) {
		return os.ErrNotExist
	}
	return nil
}

// MkdirAll creates a directory and all necessary parents.
func (mfs *MemFS) MkdirAll(path string, perm fs.FileMode) error {
	p, err := cleanPath(path)
	if err != nil {
		return err
	}
	// "." is the root, no need to record
	if p == "." {
		return nil
	}
	mfs.dirs.Add(p)
	mfs.addParentDirs(p)
	return nil
}

// OpenFile opens a file in memory with the given flags.
func (mfs *MemFS) OpenFile(name string, flag int, perm fs.FileMode) (File, error) {
	p, err := cleanPath(name)
	if err != nil {
		return nil, err
	}

	// For write operations, create a new file
	if flag&(os.O_CREATE|os.O_WRONLY|os.O_RDWR) != 0 {
		return &memFileHandle{
			mfs:       mfs,
			filename:  p,
			perm:      perm,
			buffer:    &bytes.Buffer{},
			writeMode: true,
		}, nil
	}

	// For read operations, check if file exists
	file, ok := mfs.files.Get(p)
	if !ok {
		return nil, os.ErrNotExist
	}

	return &memFileHandle{
		mfs:       mfs,
		filename:  p,
		buffer:    bytes.NewBuffer(bytes.Clone(file.data)),
		writeMode: false,
	}, nil
}

// WalkDir walks the file tree rooted at root, calling fn for each file or
// directory in the tree, including root. Entries are yielded in lexical order
// and directory entries are synthesized from file paths.
func (mfs *MemFS) WalkDir(root string, fn fs.WalkDirFunc) error {
	root, err := cleanPath(root)
	if err != nil {
		return err
	}

	// collect all entries under root: files + explicit dirs + synthesized dirs
	dirSet := make(map[string]struct{})
	fileSet := make(map[string]struct{})

	for _, k := range mfs.files.Keys() {
		if !pathMatchesRoot(root, k) {
			continue
		}
		fileSet[k] = struct{}{}
		// synthesize parent dirs between root and this file
		for p := filepath.Dir(k); p != root && pathMatchesRoot(root, p); p = filepath.Dir(p) {
			dirSet[p] = struct{}{}
		}
	}

	// include explicit dirs that match root (but not root itself — added separately)
	for _, d := range mfs.dirs.Values() {
		if d != root && pathMatchesRoot(root, d) {
			dirSet[d] = struct{}{}
		}
	}

	// remove root from dirSet if synthesized above
	delete(dirSet, root)

	// build sorted list: root first, then all dirs and files in lexical order
	entries := make([]string, 0, len(dirSet)+len(fileSet)+1)
	entries = append(entries, root)
	for d := range dirSet {
		entries = append(entries, d)
	}
	for f := range fileSet {
		entries = append(entries, f)
	}
	sort.Strings(entries[1:])

	// track which directories to skip
	var skipPrefix string

	for _, path := range entries {
		// if we're skipping a directory subtree, check prefix
		if skipPrefix != "" {
			if strings.HasPrefix(path, skipPrefix) {
				continue
			}
			skipPrefix = ""
		}

		_, isFile := fileSet[path]
		if isFile {
			file, ok := mfs.files.Get(path)
			if !ok {
				continue
			}
			info := &memFileInfo{
				name:    filepath.Base(path),
				size:    int64(len(file.data)),
				mode:    file.mode,
				modTime: file.modTime,
			}
			if err := fn(path, &memDirEntry{info: info}, nil); err != nil {
				if err == fs.SkipDir || err == fs.SkipAll {
					return nil
				}
				return err
			}
			continue
		}

		// directory entry
		info := &memFileInfo{
			name:    filepath.Base(path),
			mode:    fs.ModeDir | 0o755,
			modTime: time.Now(),
		}
		if err := fn(path, &memDirEntry{info: info}, nil); err != nil {
			if err == fs.SkipAll {
				return nil
			}
			if err == fs.SkipDir {
				skipPrefix = path + string(filepath.Separator)
				continue
			}
			return err
		}
	}

	return nil
}

// pathMatchesRoot returns true if path is under root or equal to root.
func pathMatchesRoot(root, path string) bool {
	if root == "." {
		return true
	}
	return path == root || strings.HasPrefix(path, root+string(filepath.Separator))
}

// ClearCacheFiles removes all files with .cache extension from memory.
// This is a specialized method for cache implementations.
func (mfs *MemFS) ClearCacheFiles() {
	keys := mfs.files.Keys()
	for _, filename := range keys {
		if strings.HasSuffix(filename, ".cache") {
			mfs.files.Delete(filename)
		}
	}
}

// CountCacheFiles returns the number of files with .cache extension.
// This is a specialized method for cache implementations.
func (mfs *MemFS) CountCacheFiles() int {
	keys := mfs.files.Keys()
	count := 0
	for _, filename := range keys {
		if strings.HasSuffix(filename, ".cache") {
			count++
		}
	}
	return count
}

type memFileInfo struct {
	name    string
	size    int64
	mode    fs.FileMode
	modTime time.Time
}

// Name returns the file name.
func (fi *memFileInfo) Name() string { return fi.name }

// Size returns the file size in bytes.
func (fi *memFileInfo) Size() int64 { return fi.size }

// Mode returns the file mode.
func (fi *memFileInfo) Mode() fs.FileMode { return fi.mode }

// ModTime returns the modification time.
func (fi *memFileInfo) ModTime() time.Time { return fi.modTime }

// IsDir reports whether the entry is a directory.
func (fi *memFileInfo) IsDir() bool { return fi.mode.IsDir() }

// Sys returns nil.
func (fi *memFileInfo) Sys() any { return nil }

type memDirEntry struct {
	info *memFileInfo
}

// Name returns the entry name.
func (de *memDirEntry) Name() string { return de.info.Name() }

// IsDir reports whether the entry is a directory.
func (de *memDirEntry) IsDir() bool { return de.info.IsDir() }

// Type returns the file mode type bits.
func (de *memDirEntry) Type() fs.FileMode { return de.info.Mode().Type() }

// Info returns the file info.
func (de *memDirEntry) Info() (fs.FileInfo, error) { return de.info, nil }

// memFileHandle implements the File interface for in-memory files
type memFileHandle struct {
	mfs       *MemFS
	filename  string
	buffer    *bytes.Buffer
	perm      fs.FileMode
	writeMode bool
	closed    bool
}

// Stat returns the file info.
func (fh *memFileHandle) Stat() (fs.FileInfo, error) {
	if fh.closed {
		return nil, fs.ErrClosed
	}

	file, ok := fh.mfs.files.Get(fh.filename)
	if !ok {
		// If file doesn't exist yet (new file), return basic info
		return &memFileInfo{
			name:    filepath.Base(fh.filename),
			size:    int64(fh.buffer.Len()),
			mode:    fh.perm,
			modTime: time.Now(),
		}, nil
	}

	return &memFileInfo{
		name:    filepath.Base(fh.filename),
		size:    int64(len(file.data)),
		mode:    file.mode,
		modTime: file.modTime,
	}, nil
}

func (fh *memFileHandle) Read(p []byte) (int, error) {
	if fh.closed {
		return 0, fs.ErrClosed
	}
	return fh.buffer.Read(p)
}

func (fh *memFileHandle) Write(p []byte) (int, error) {
	if fh.closed {
		return 0, fs.ErrClosed
	}
	if !fh.writeMode {
		return 0, fs.ErrPermission
	}
	return fh.buffer.Write(p)
}

// Close flushes writes and closes the handle.
func (fh *memFileHandle) Close() error {
	if fh.closed {
		return fs.ErrClosed
	}
	fh.closed = true

	// If in write mode, save the buffer contents to the filesystem
	if fh.writeMode {
		fh.mfs.files.Set(fh.filename, &memFile{
			data:    fh.buffer.Bytes(),
			modTime: time.Now(),
			mode:    fh.perm,
		})
	}

	return nil
}
