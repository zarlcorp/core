package zfilesystem_test

import (
	"io/fs"
	"path/filepath"
	"sort"
	"testing"

	"github.com/zarlcorp/core/pkg/zfilesystem"
)

func TestOSFileSystem_WalkDir_RelativePaths(t *testing.T) {
	tmpDir := t.TempDir()
	osfs := zfilesystem.NewOSFileSystem(tmpDir)

	// create files in a subdirectory
	if err := osfs.MkdirAll("identities", 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	files := map[string][]byte{
		"identities/alice.enc": []byte("alice-data"),
		"identities/bob.enc":   []byte("bob-data"),
		"root.txt":             []byte("root-data"),
	}
	for name, data := range files {
		if err := osfs.WriteFile(name, data, 0o644); err != nil {
			t.Fatalf("WriteFile(%s): %v", name, err)
		}
	}

	t.Run("callback receives relative paths", func(t *testing.T) {
		var paths []string
		err := osfs.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			paths = append(paths, path)
			return nil
		})
		if err != nil {
			t.Fatalf("WalkDir: %v", err)
		}

		for _, p := range paths {
			if filepath.IsAbs(p) {
				t.Errorf("callback received absolute path: %s", p)
			}
		}
	})

	t.Run("subdirectory walk returns paths relative to baseDir", func(t *testing.T) {
		var paths []string
		err := osfs.WalkDir("identities", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				paths = append(paths, path)
			}
			return nil
		})
		if err != nil {
			t.Fatalf("WalkDir: %v", err)
		}

		sort.Strings(paths)
		want := []string{"identities/alice.enc", "identities/bob.enc"}
		if len(paths) != len(want) {
			t.Fatalf("got %v, want %v", paths, want)
		}
		for i := range want {
			if paths[i] != want[i] {
				t.Errorf("paths[%d] = %q, want %q", i, paths[i], want[i])
			}
		}
	})

	t.Run("walked paths can be passed back to ReadFile", func(t *testing.T) {
		err := osfs.WalkDir("identities", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			data, readErr := osfs.ReadFile(path)
			if readErr != nil {
				t.Errorf("ReadFile(%q): %v", path, readErr)
				return nil
			}
			expected, ok := files[path]
			if !ok {
				t.Errorf("unexpected file: %s", path)
				return nil
			}
			if string(data) != string(expected) {
				t.Errorf("ReadFile(%q) = %q, want %q", path, data, expected)
			}
			return nil
		})
		if err != nil {
			t.Fatalf("WalkDir: %v", err)
		}
	})

	t.Run("walked paths can be passed back to WriteFile", func(t *testing.T) {
		err := osfs.WalkDir("identities", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			// overwrite each file via its walked path
			return osfs.WriteFile(path, []byte("updated"), 0o644)
		})
		if err != nil {
			t.Fatalf("WalkDir: %v", err)
		}

		// verify writes succeeded
		data, err := osfs.ReadFile("identities/alice.enc")
		if err != nil {
			t.Fatalf("ReadFile: %v", err)
		}
		if string(data) != "updated" {
			t.Errorf("ReadFile after walk-write = %q, want %q", data, "updated")
		}
	})

	t.Run("root entry is relative", func(t *testing.T) {
		var rootPath string
		err := osfs.WalkDir("identities", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				rootPath = path
				return fs.SkipAll
			}
			return nil
		})
		if err != nil {
			t.Fatalf("WalkDir: %v", err)
		}
		if rootPath != "identities" {
			t.Errorf("root dir path = %q, want %q", rootPath, "identities")
		}
	})
}

