// Package zfilesystem provides file system abstractions for different storage backends.
//
// This package offers clean interfaces and implementations for file operations,
// enabling dependency injection and testing through in-memory filesystems while
// supporting production use cases with OS filesystems.
//
// # Interfaces
//
// ReadWriteFileFS: Complete file system interface with read, write, and directory operations
//
// # Implementations
//
// MemFS: In-memory filesystem for testing and development
// OSFileSystem: OS-backed filesystem for production use
//
// # Usage Example
//
//	// In-memory filesystem for testing
//	memfs := zfilesystem.NewMemFS()
//	memfs.WriteFile("test.txt", []byte("data"), 0644)
//
//	// OS filesystem for production
//	osfs := zfilesystem.NewOSFileSystem("/path/to/base/dir")
//	osfs.WriteFile("prod.txt", []byte("data"), 0644)
package zfilesystem

import (
	"io"
	"io/fs"
)

// ReadFileFS defines the interface for reading files.
type ReadFileFS interface {
	// ReadFile reads the file named by filename and returns the contents.
	ReadFile(filename string) ([]byte, error)
}

// WriteFileFS defines the interface for writing files.
type WriteFileFS interface {
	// WriteFile writes data to a file named by filename.
	WriteFile(filename string, data []byte, perm fs.FileMode) error
}

// RemoveFS defines the interface for removing files.
type RemoveFS interface {
	// Remove removes the named file.
	Remove(filename string) error
}

// MkdirFS defines the interface for creating directories.
type MkdirFS interface {
	// MkdirAll creates a directory named path, along with any necessary parents.
	MkdirAll(path string, perm fs.FileMode) error
}

// File is a file that can be read from and written to.
type File interface {
	fs.File
	io.Writer
}

// OpenFileFS defines the interface for opening files with flags.
type OpenFileFS interface {
	// OpenFile opens the named file with specified flag and perm.
	OpenFile(name string, flag int, perm fs.FileMode) (File, error)
}

// WalkDirFS defines the interface for walking directory trees.
type WalkDirFS interface {
	// WalkDir walks the file tree rooted at root, calling fn for each file or
	// directory in the tree, including root.
	WalkDir(root string, fn fs.WalkDirFunc) error
}

// ReadWriteFileFS defines a complete file system interface that combines
// all basic file operations through interface composition.
type ReadWriteFileFS interface {
	ReadFileFS
	WriteFileFS
	RemoveFS
	MkdirFS
	OpenFileFS
	WalkDirFS
}
