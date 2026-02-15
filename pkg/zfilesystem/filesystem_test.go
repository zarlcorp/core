package zfilesystem_test

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/zarlcorp/core/pkg/zfilesystem"
)

func TestReadWriteFileFS_Contract(t *testing.T) {
	tests := []struct {
		name  string
		newFS func(t *testing.T) (zfilesystem.ReadWriteFileFS, func())
	}{
		{
			name: "MemFS",
			newFS: func(t *testing.T) (zfilesystem.ReadWriteFileFS, func()) {
				return zfilesystem.NewMemFS(), func() {}
			},
		},
		{
			name: "OSFileSystem",
			newFS: func(t *testing.T) (zfilesystem.ReadWriteFileFS, func()) {
				tmpDir := t.TempDir()
				return zfilesystem.NewOSFileSystem(tmpDir), func() {}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs, cleanup := tt.newFS(t)
			defer cleanup()

			testReadWriteFileFS(t, fs)
		})
	}
}

func testReadWriteFileFS(t *testing.T, fs zfilesystem.ReadWriteFileFS) {
	t.Helper()

	t.Run("WriteFile and ReadFile operations", func(t *testing.T) {
		tests := []struct {
			name     string
			filename string
			data     []byte
			perm     os.FileMode
		}{
			{
				name:     "simple text file",
				filename: "test.txt",
				data:     []byte("hello world"),
				perm:     0o644,
			},
			{
				name:     "binary data",
				filename: "binary.dat",
				data:     []byte{0x00, 0x01, 0x02, 0xFF, 0xFE},
				perm:     0o644,
			},
			{
				name:     "empty file",
				filename: "empty.txt",
				data:     []byte{},
				perm:     0o644,
			},
			{
				name:     "large content",
				filename: "large.txt",
				data:     []byte(strings.Repeat("x", 10000)),
				perm:     0o644,
			},
			{
				name:     "special characters in content",
				filename: "special.txt",
				data:     []byte("hello\nworld\t\r\n日本語"),
				perm:     0o644,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := fs.WriteFile(tt.filename, tt.data, tt.perm)
				if err != nil {
					t.Fatalf("WriteFile() error = %v", err)
				}

				data, err := fs.ReadFile(tt.filename)
				if err != nil {
					t.Fatalf("ReadFile() error = %v", err)
				}

				if !bytes.Equal(data, tt.data) {
					t.Errorf("ReadFile() data = %q, want %q", string(data), string(tt.data))
				}

				fs.Remove(tt.filename)
			})
		}
	})

	t.Run("ReadFile non-existent file", func(t *testing.T) {
		_, err := fs.ReadFile("non-existent.txt")
		if err == nil {
			t.Error("ReadFile() expected error for non-existent file")
		}
		if !os.IsNotExist(err) {
			t.Errorf("ReadFile() error should be os.ErrNotExist-compatible, got %v", err)
		}
	})

	t.Run("Remove operations", func(t *testing.T) {
		tests := []struct {
			name  string
			setup func() string
			error error
		}{
			{
				name: "remove existing file",
				setup: func() string {
					filename := "to-remove.txt"
					fs.WriteFile(filename, []byte("test"), 0o644)
					return filename
				},
				error: nil,
			},
			{
				name: "remove non-existent file",
				setup: func() string {
					return "non-existent.txt"
				},
				error: os.ErrNotExist,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				filename := tt.setup()
				err := fs.Remove(filename)

				if tt.error != nil {
					if err == nil {
						t.Error("Remove() expected error")
						return
					}
					if !os.IsNotExist(err) {
						t.Errorf("Remove() wrong error type: %v", err)
					}
					return
				}

				if err != nil {
					t.Errorf("Remove() unexpected error = %v", err)
					return
				}

				_, err = fs.ReadFile(filename)
				if err == nil {
					t.Error("ReadFile() should fail after Remove()")
				}
			})
		}
	})

	t.Run("File overwrite", func(t *testing.T) {
		filename := "overwrite-test.txt"
		defer fs.Remove(filename)

		err := fs.WriteFile(filename, []byte("original"), 0o644)
		if err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}

		err = fs.WriteFile(filename, []byte("overwritten"), 0o644)
		if err != nil {
			t.Fatalf("WriteFile() overwrite error = %v", err)
		}

		data, err := fs.ReadFile(filename)
		if err != nil {
			t.Fatalf("ReadFile() error = %v", err)
		}

		if string(data) != "overwritten" {
			t.Errorf("ReadFile() after overwrite = %q, want %q", string(data), "overwritten")
		}
	})

	t.Run("WalkDir operations", func(t *testing.T) {
		files := map[string][]byte{
			"walk1.txt": []byte("content1"),
			"walk2.txt": []byte("content2"),
			"data.json": []byte(`{"key": "value"}`),
			"readme.md": []byte("# README"),
		}

		for filename, content := range files {
			err := fs.WriteFile(filename, content, 0o644)
			if err != nil {
				t.Fatalf("WriteFile(%s) error = %v", filename, err)
			}
		}
		defer func() {
			for filename := range files {
				fs.Remove(filename)
			}
		}()

		var walked []string
		err := fs.WalkDir(".", func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				walked = append(walked, path)
			}
			return nil
		})
		if err != nil {
			t.Fatalf("WalkDir() error = %v", err)
		}

		if len(walked) < len(files) {
			t.Errorf("WalkDir() found %d files, want at least %d", len(walked), len(files))
		}

		for filename := range files {
			found := false
			for _, walkedFile := range walked {
				if strings.HasSuffix(walkedFile, filename) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("WalkDir() missing file: %s", filename)
			}
		}
	})

	t.Run("WalkDir with filter function", func(t *testing.T) {
		files := []string{"filter-test.txt", "filter-data.json", "filter-config.yaml", "filter-script.sh"}
		for _, filename := range files {
			fs.WriteFile(filename, []byte("content"), 0o644)
		}
		defer func() {
			for _, filename := range files {
				fs.Remove(filename)
			}
		}()

		var txtFiles []string
		err := fs.WalkDir(".", func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() && strings.HasSuffix(path, ".txt") {
				txtFiles = append(txtFiles, path)
			}
			return nil
		})
		if err != nil {
			t.Fatalf("WalkDir() error = %v", err)
		}

		found := false
		for _, txtFile := range txtFiles {
			if strings.HasSuffix(txtFile, "filter-test.txt") {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("WalkDir() with filter should contain filter-test.txt, found %v", txtFiles)
		}
	})

	t.Run("WalkDir early termination", func(t *testing.T) {
		files := []string{"term1.txt", "term2.txt", "term3.txt"}
		for _, filename := range files {
			fs.WriteFile(filename, []byte("content"), 0o644)
		}
		defer func() {
			for _, filename := range files {
				fs.Remove(filename)
			}
		}()

		testErr := errors.New("stop walking")
		count := 0
		err := fs.WalkDir(".", func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() && strings.Contains(path, "term") {
				count++
				if count == 2 {
					return testErr
				}
			}
			return nil
		})

		if !errors.Is(err, testErr) {
			t.Errorf("WalkDir() error = %v, want %v", err, testErr)
		}

		if count != 2 {
			t.Errorf("WalkDir() visited %d files before stopping, want 2", count)
		}
	})

	t.Run("path traversal protection", func(t *testing.T) {
		rejected := []struct {
			name string
			path string
		}{
			{name: "parent escape", path: "../escape"},
			{name: "nested parent escape", path: "foo/../../escape"},
			{name: "absolute path", path: "/etc/passwd"},
		}

		for _, tt := range rejected {
			t.Run(tt.name+" ReadFile", func(t *testing.T) {
				_, err := fs.ReadFile(tt.path)
				if err == nil {
					t.Errorf("ReadFile(%q) should be rejected", tt.path)
				}
				if err != nil && !strings.Contains(err.Error(), "path escapes base directory") {
					t.Errorf("ReadFile(%q) error = %v, want path escapes error", tt.path, err)
				}
			})

			t.Run(tt.name+" WriteFile", func(t *testing.T) {
				err := fs.WriteFile(tt.path, []byte("bad"), 0o644)
				if err == nil {
					t.Errorf("WriteFile(%q) should be rejected", tt.path)
				}
				if err != nil && !strings.Contains(err.Error(), "path escapes base directory") {
					t.Errorf("WriteFile(%q) error = %v, want path escapes error", tt.path, err)
				}
			})

			t.Run(tt.name+" Remove", func(t *testing.T) {
				err := fs.Remove(tt.path)
				if err == nil {
					t.Errorf("Remove(%q) should be rejected", tt.path)
				}
				if err != nil && !strings.Contains(err.Error(), "path escapes base directory") {
					t.Errorf("Remove(%q) error = %v, want path escapes error", tt.path, err)
				}
			})

			t.Run(tt.name+" MkdirAll", func(t *testing.T) {
				err := fs.MkdirAll(tt.path, 0o755)
				if err == nil {
					t.Errorf("MkdirAll(%q) should be rejected", tt.path)
				}
				if err != nil && !strings.Contains(err.Error(), "path escapes base directory") {
					t.Errorf("MkdirAll(%q) error = %v, want path escapes error", tt.path, err)
				}
			})

			t.Run(tt.name+" OpenFile", func(t *testing.T) {
				_, err := fs.OpenFile(tt.path, os.O_RDONLY, 0o644)
				if err == nil {
					t.Errorf("OpenFile(%q) should be rejected", tt.path)
				}
				if err != nil && !strings.Contains(err.Error(), "path escapes base directory") {
					t.Errorf("OpenFile(%q) error = %v, want path escapes error", tt.path, err)
				}
			})

			t.Run(tt.name+" WalkDir", func(t *testing.T) {
				err := fs.WalkDir(tt.path, func(path string, d os.DirEntry, err error) error {
					return err
				})
				if err == nil {
					t.Errorf("WalkDir(%q) should be rejected", tt.path)
				}
				if err != nil && !strings.Contains(err.Error(), "path escapes base directory") {
					t.Errorf("WalkDir(%q) error = %v, want path escapes error", tt.path, err)
				}
			})
		}

		// paths that resolve within base should be allowed
		allowed := []struct {
			name string
			path string
		}{
			{name: "safe traversal", path: "foo/../bar"},
			{name: "dot prefix", path: "./normal/path"},
			{name: "simple path", path: "simple.txt"},
		}

		for _, tt := range allowed {
			t.Run(tt.name+" WriteFile allowed", func(t *testing.T) {
				err := fs.WriteFile(tt.path, []byte("ok"), 0o644)
				if err != nil && strings.Contains(err.Error(), "path escapes base directory") {
					t.Errorf("WriteFile(%q) should be allowed: %v", tt.path, err)
				}
			})
		}

		// empty string should be treated as current dir (not rejected)
		t.Run("empty path MkdirAll", func(t *testing.T) {
			err := fs.MkdirAll("", 0o755)
			if err != nil && strings.Contains(err.Error(), "path escapes base directory") {
				t.Errorf("MkdirAll(\"\") should not escape: %v", err)
			}
		})
	})
}
