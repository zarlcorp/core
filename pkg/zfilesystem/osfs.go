package zfilesystem

import (
	"io/fs"
	"os"
	"path/filepath"
)

var _ ReadWriteFileFS = (*OSFileSystem)(nil)

// OSFileSystem implements ReadWriteFileFS using the standard os package.
type OSFileSystem struct {
	baseDir string
}

// NewOSFileSystem creates a new OS filesystem with the specified base directory.
// The base directory will be used as the root for all file operations.
func NewOSFileSystem(baseDir string) *OSFileSystem {
	return &OSFileSystem{baseDir: baseDir}
}

// ReadFile reads a file from the OS filesystem.
func (osfs *OSFileSystem) ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filepath.Join(osfs.baseDir, filename))
}

// WriteFile writes data to a file in the OS filesystem.
func (osfs *OSFileSystem) WriteFile(filename string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(filepath.Join(osfs.baseDir, filename), data, perm)
}

// Remove removes a file from the OS filesystem.
func (osfs *OSFileSystem) Remove(filename string) error {
	return os.Remove(filepath.Join(osfs.baseDir, filename))
}

// MkdirAll creates a directory and all necessary parents.
func (osfs *OSFileSystem) MkdirAll(path string, perm fs.FileMode) error {
	return os.MkdirAll(filepath.Join(osfs.baseDir, path), perm)
}

// OpenFile opens the named file with specified flag and perm.
func (osfs *OSFileSystem) OpenFile(name string, flag int, perm fs.FileMode) (File, error) {
	return os.OpenFile(filepath.Join(osfs.baseDir, name), flag, perm)
}

// WalkDir walks the directory tree in the OS filesystem.
func (osfs *OSFileSystem) WalkDir(root string, fn fs.WalkDirFunc) error {
	return filepath.WalkDir(filepath.Join(osfs.baseDir, root), fn)
}

// BaseDir returns the base directory for this filesystem.
func (osfs *OSFileSystem) BaseDir() string {
	return osfs.baseDir
}
