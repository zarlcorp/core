package zcache_test

import (
	"testing"
	"time"

	"github.com/zarlcorp/core/pkg/zcache"
)

func TestMemoryCache_Constructor(t *testing.T) {
	tests := []struct {
		name  string
		build func() any
		check func(t *testing.T, c any)
	}{
		{
			name: "string to int",
			build: func() any {
				return zcache.NewMemoryCache[string, int]()
			},
			check: func(t *testing.T, c any) {
				mc := c.(*zcache.MemoryCache[string, int])
				if mc == nil {
					t.Error("NewMemoryCache() returned nil")
				}
				if got, _ := mc.Len(t.Context()); got != 0 {
					t.Errorf("Len() = %v, want 0", got)
				}
			},
		},
		{
			name: "string to string",
			build: func() any {
				return zcache.NewMemoryCache[string, string]()
			},
			check: func(t *testing.T, c any) {
				if c == nil {
					t.Error("NewMemoryCache[string, string]() returned nil")
				}
			},
		},
		{
			name: "int to bool",
			build: func() any {
				return zcache.NewMemoryCache[int, bool]()
			},
			check: func(t *testing.T, c any) {
				if c == nil {
					t.Error("NewMemoryCache[int, bool]() returned nil")
				}
			},
		},
		{
			name: "custom struct types",
			build: func() any {
				type customKey struct{ ID int }
				type customValue struct{ Data string }
				return zcache.NewMemoryCache[customKey, customValue]()
			},
			check: func(t *testing.T, c any) {
				if c == nil {
					t.Error("NewMemoryCache[customKey, customValue]() returned nil")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.build()
			tt.check(t, c)
		})
	}
}

func TestMemoryCache_TTL(t *testing.T) {
	t.Run("entry expires after TTL", func(t *testing.T) {
		c := zcache.NewMemoryCache[string, int](
			zcache.WithMemoryTTL[string, int](50 * time.Millisecond),
		)
		ctx := t.Context()

		c.Set(ctx, "key", 42)

		// should be accessible immediately
		got, err := c.Get(ctx, "key")
		if err != nil {
			t.Fatalf("Get() before expiry: %v", err)
		}
		if got != 42 {
			t.Errorf("Get() = %v, want 42", got)
		}

		// wait for TTL to expire
		time.Sleep(60 * time.Millisecond)

		_, err = c.Get(ctx, "key")
		if err != zcache.ErrNotFound {
			t.Errorf("Get() after expiry: got %v, want ErrNotFound", err)
		}
	})

	t.Run("entry accessible before TTL", func(t *testing.T) {
		c := zcache.NewMemoryCache[string, int](
			zcache.WithMemoryTTL[string, int](1 * time.Second),
		)
		ctx := t.Context()

		c.Set(ctx, "key", 99)

		got, err := c.Get(ctx, "key")
		if err != nil {
			t.Fatalf("Get() before expiry: %v", err)
		}
		if got != 99 {
			t.Errorf("Get() = %v, want 99", got)
		}
	})

	t.Run("zero TTL means no expiry", func(t *testing.T) {
		c := zcache.NewMemoryCache[string, int]()
		ctx := t.Context()

		c.Set(ctx, "key", 1)

		got, err := c.Get(ctx, "key")
		if err != nil {
			t.Fatalf("Get() error = %v", err)
		}
		if got != 1 {
			t.Errorf("Get() = %v, want 1", got)
		}
	})

	t.Run("Len excludes expired entries", func(t *testing.T) {
		c := zcache.NewMemoryCache[string, int](
			zcache.WithMemoryTTL[string, int](50 * time.Millisecond),
		)
		ctx := t.Context()

		c.Set(ctx, "a", 1)
		c.Set(ctx, "b", 2)

		got, _ := c.Len(ctx)
		if got != 2 {
			t.Errorf("Len() before expiry = %v, want 2", got)
		}

		time.Sleep(60 * time.Millisecond)

		// add one fresh entry after the old ones expire
		c.Set(ctx, "c", 3)

		got, _ = c.Len(ctx)
		if got != 1 {
			t.Errorf("Len() after partial expiry = %v, want 1", got)
		}
	})

	t.Run("Set refreshes TTL on existing key", func(t *testing.T) {
		c := zcache.NewMemoryCache[string, int](
			zcache.WithMemoryTTL[string, int](80 * time.Millisecond),
		)
		ctx := t.Context()

		c.Set(ctx, "key", 1)

		// wait 50ms, then re-set to refresh
		time.Sleep(50 * time.Millisecond)
		c.Set(ctx, "key", 2)

		// wait another 50ms -- would have expired under the original insertion time
		time.Sleep(50 * time.Millisecond)

		got, err := c.Get(ctx, "key")
		if err != nil {
			t.Fatalf("Get() after refresh: %v", err)
		}
		if got != 2 {
			t.Errorf("Get() = %v, want 2", got)
		}
	})
}
