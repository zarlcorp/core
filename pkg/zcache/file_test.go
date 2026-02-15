package zcache_test

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/zarlcorp/core/pkg/zcache"
	"github.com/zarlcorp/core/pkg/zfilesystem"
)

func TestFileCache_Constructor(t *testing.T) {
	tests := []struct {
		name  string
		setup func() *zcache.FileCache[string, int]
	}{
		{
			name: "with default temp directory",
			setup: func() *zcache.FileCache[string, int] {
				return zcache.NewFileCache[string, int]()
			},
		},
		{
			name: "with custom directory",
			setup: func() *zcache.FileCache[string, int] {
				return zcache.NewFileCache[string, int](zcache.WithOSFileSystem[string, int](t.TempDir()))
			},
		},
		{
			name: "with in-memory filesystem",
			setup: func() *zcache.FileCache[string, int] {
				return zcache.NewFileCache[string, int](zcache.WithFileSystem[string, int](zfilesystem.NewMemFS()))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.setup()
			if c == nil {
				t.Error("NewFileCache() returned nil")
			}
		})
	}
}

func TestFileCache_Persistence(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value int
		check func(t *testing.T, c1, c2 *zcache.FileCache[string, int])
	}{
		{
			name:  "data persists across instances",
			key:   "persistent",
			value: 123,
			check: func(t *testing.T, c1, c2 *zcache.FileCache[string, int]) {
				ctx := t.Context()
				c1.Set(ctx, "persistent", 123)

				if got, _ := c1.Len(ctx); got != 1 {
					t.Errorf("first instance Len() = %v, want 1", got)
				}

				value, err := c2.Get(ctx, "persistent")
				if err != nil {
					t.Errorf("second instance Get() error = %v, want nil", err)
					return
				}

				if value != 123 {
					t.Errorf("second instance Get() = %v, want 123", value)
				}

				if got, _ := c2.Len(ctx); got != 1 {
					t.Errorf("second instance Len() = %v, want 1", got)
				}
			},
		},
		{
			name:  "multiple values persist",
			key:   "multiple",
			value: 456,
			check: func(t *testing.T, c1, c2 *zcache.FileCache[string, int]) {
				ctx := t.Context()
				c1.Set(ctx, "key1", 1)
				c1.Set(ctx, "key2", 2)
				c1.Set(ctx, "key3", 3)

				if got, _ := c2.Len(ctx); got != 3 {
					t.Errorf("second instance Len() = %v, want 3", got)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			c1 := zcache.NewFileCache[string, int](zcache.WithOSFileSystem[string, int](tmpDir))
			c2 := zcache.NewFileCache[string, int](zcache.WithOSFileSystem[string, int](tmpDir))

			tt.check(t, c1, c2)
		})
	}
}

func TestFileCache_FileHandling(t *testing.T) {
	tests := []struct {
		name  string
		setup func(c *zcache.FileCache[string, string], dir string)
		check func(t *testing.T, c *zcache.FileCache[string, string], dir string)
	}{
		{
			name: "files created with .cache extension",
			setup: func(c *zcache.FileCache[string, string], dir string) {
				c.Set(t.Context(), "test", "value")
			},
			check: func(t *testing.T, c *zcache.FileCache[string, string], dir string) {
				files, err := filepath.Glob(filepath.Join(dir, "*.cache"))
				if err != nil {
					t.Errorf("glob cache files: %v", err)
					return
				}

				if len(files) != 1 {
					t.Errorf("found %d cache files, want 1", len(files))
				}
			},
		},
		{
			name: "special characters sanitized",
			setup: func(c *zcache.FileCache[string, string], dir string) {
				c.Set(t.Context(), "key/with\\special:chars*?\"<>|", "value")
			},
			check: func(t *testing.T, c *zcache.FileCache[string, string], dir string) {
				value, err := c.Get(t.Context(), "key/with\\special:chars*?\"<>|")
				if err != nil {
					t.Errorf("Get() with special chars error = %v, want nil", err)
					return
				}

				if value != "value" {
					t.Errorf("Get() with special chars = %v, want value", value)
				}
			},
		},
		{
			name: "multiple files created",
			setup: func(c *zcache.FileCache[string, string], dir string) {
				ctx := t.Context()
				c.Set(ctx, "file1", "value1")
				c.Set(ctx, "file2", "value2")
				c.Set(ctx, "file3", "value3")
			},
			check: func(t *testing.T, c *zcache.FileCache[string, string], dir string) {
				files, err := filepath.Glob(filepath.Join(dir, "*.cache"))
				if err != nil {
					t.Errorf("glob cache files: %v", err)
					return
				}

				if len(files) != 3 {
					t.Errorf("found %d cache files, want 3", len(files))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			c := zcache.NewFileCache[string, string](zcache.WithOSFileSystem[string, string](tmpDir))

			tt.setup(c, tmpDir)
			tt.check(t, c, tmpDir)
		})
	}
}

func TestFileCache_NoKeyCollision(t *testing.T) {
	// keys that would have collided with the old char-replacement sanitization
	tests := []struct {
		name string
		keyA string
		keyB string
	}{
		{
			name: "slash vs colon",
			keyA: "a/b",
			keyB: "a:b",
		},
		{
			name: "backslash vs pipe",
			keyA: "x\\y",
			keyB: "x|y",
		},
		{
			name: "star vs question mark",
			keyA: "foo*bar",
			keyB: "foo?bar",
		},
		{
			name: "angle brackets",
			keyA: "a<b",
			keyB: "a>b",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := zcache.NewFileCache[string, string](
				zcache.WithFileSystem[string, string](zfilesystem.NewMemFS()),
			)
			ctx := t.Context()

			if err := c.Set(ctx, tt.keyA, "valueA"); err != nil {
				t.Fatalf("Set(%q) error = %v", tt.keyA, err)
			}
			if err := c.Set(ctx, tt.keyB, "valueB"); err != nil {
				t.Fatalf("Set(%q) error = %v", tt.keyB, err)
			}

			n, err := c.Len(ctx)
			if err != nil {
				t.Fatalf("Len() error = %v", err)
			}
			if n != 2 {
				t.Errorf("Len() = %d, want 2 (keys collided)", n)
			}

			gotA, err := c.Get(ctx, tt.keyA)
			if err != nil {
				t.Fatalf("Get(%q) error = %v", tt.keyA, err)
			}
			if gotA != "valueA" {
				t.Errorf("Get(%q) = %q, want %q", tt.keyA, gotA, "valueA")
			}

			gotB, err := c.Get(ctx, tt.keyB)
			if err != nil {
				t.Fatalf("Get(%q) error = %v", tt.keyB, err)
			}
			if gotB != "valueB" {
				t.Errorf("Get(%q) = %q, want %q", tt.keyB, gotB, "valueB")
			}

			// delete one key, the other should remain
			deleted, err := c.Delete(ctx, tt.keyA)
			if err != nil {
				t.Fatalf("Delete(%q) error = %v", tt.keyA, err)
			}
			if !deleted {
				t.Errorf("Delete(%q) = false, want true", tt.keyA)
			}

			_, err = c.Get(ctx, tt.keyA)
			if !errors.Is(err, zcache.ErrNotFound) {
				t.Errorf("Get(%q) after delete error = %v, want ErrNotFound", tt.keyA, err)
			}

			gotB, err = c.Get(ctx, tt.keyB)
			if err != nil {
				t.Errorf("Get(%q) after deleting other key error = %v", tt.keyB, err)
			}
			if gotB != "valueB" {
				t.Errorf("Get(%q) after deleting other key = %q, want %q", tt.keyB, gotB, "valueB")
			}
		})
	}
}
