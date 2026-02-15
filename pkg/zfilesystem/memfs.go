package zfilesystem

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/zarlcorp/core/pkg/zsync"
)

var _ ReadWriteFileFS = (*MemFS)(nil)

// MemFS provides an in-memory filesystem implementation.
// It implements ReadWriteFileFS and is safe for concurrent use.
type MemFS struct {
	files *zsync.ZMap[string, *memFile]
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
	}
}

// ReadFile reads the file from memory.
func (mfs *MemFS) ReadFile(filename string) ([]byte, error) {
	file, err := mfs.files.Get(filename)
	if err != nil {
		return nil, os.ErrNotExist
	}

	return bytes.Clone(file.data), nil
}

// WriteFile writes data to memory.
func (mfs *MemFS) WriteFile(filename string, data []byte, perm fs.FileMode) error {
	mfs.files.Set(filename, &memFile{
		data:    bytes.Clone(data),
		modTime: time.Now(),
		mode:    perm,
	})

	return nil
}

// Remove removes a file from memory.
func (mfs *MemFS) Remove(filename string) error {
	if !mfs.files.Delete(filename) {
		return os.ErrNotExist
	}
	return nil
}

// MkdirAll is a no-op for in-memory filesystem as directories are implicit.
func (mfs *MemFS) MkdirAll(path string, perm fs.FileMode) error {
	return nil
}

// OpenFile opens a file in memory with the given flags.
func (mfs *MemFS) OpenFile(name string, flag int, perm fs.FileMode) (File, error) {
	// For write operations, create a new file
	if flag&(os.O_CREATE|os.O_WRONLY|os.O_RDWR) != 0 {
		return &memFileHandle{
			mfs:       mfs,
			filename:  name,
			perm:      perm,
			buffer:    &bytes.Buffer{},
			writeMode: true,
		}, nil
	}

	// For read operations, check if file exists
	file, err := mfs.files.Get(name)
	if err != nil {
		return nil, os.ErrNotExist
	}

	return &memFileHandle{
		mfs:       mfs,
		filename:  name,
		buffer:    bytes.NewBuffer(bytes.Clone(file.data)),
		writeMode: false,
	}, nil
}

// WalkDir walks through all files in memory.
func (mfs *MemFS) WalkDir(root string, fn fs.WalkDirFunc) error {
	keys := mfs.files.Keys()

	for _, filename := range keys {
		file, err := mfs.files.Get(filename)
		if err != nil {
			continue
		}

		info := &memFileInfo{
			name:    filepath.Base(filename),
			size:    int64(len(file.data)),
			mode:    file.mode,
			modTime: file.modTime,
		}

		entry := &memDirEntry{info: info}
		if err := fn(filename, entry, nil); err != nil {
			return err
		}
	}

	return nil
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

	file, err := fh.mfs.files.Get(fh.filename)
	if err != nil {
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
