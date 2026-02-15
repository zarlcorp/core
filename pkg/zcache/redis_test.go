package zcache_test

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/zarlcorp/core/pkg/zcache"
)

// newTestRedisClient starts an in-memory redis and returns a client pointed at it.
// the server is auto-closed on test cleanup.
func newTestRedisClient(t *testing.T) *redis.Client {
	s := miniredis.RunT(t)
	return redis.NewClient(&redis.Options{Addr: s.Addr()})
}

func TestRedisCache_Constructor(t *testing.T) {
	tests := []struct {
		name       string
		opts       []zcache.RedisOption[string, int]
		wantPrefix string
	}{
		{
			name:       "default prefix applied when none configured",
			wantPrefix: "zcache:",
		},
		{
			name: "with custom client",
			opts: []zcache.RedisOption[string, int]{
				zcache.WithPrefix[string, int]("test"),
				zcache.WithClient[string, int](&redis.Client{}),
			},
			wantPrefix: "test",
		},
		{
			name: "with custom prefix",
			opts: []zcache.RedisOption[string, int]{
				zcache.WithPrefix[string, int]("myapp:"),
			},
			wantPrefix: "myapp:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := zcache.NewRedisCache[string, int](tt.opts...)
			if c == nil {
				t.Fatal("NewRedisCache() returned nil")
			}
			if got := c.Prefix(); got != tt.wantPrefix {
				t.Errorf("Prefix() = %q, want %q", got, tt.wantPrefix)
			}
		})
	}
}

func TestRedisCache_Types(t *testing.T) {
	tests := []struct {
		name  string
		build func() any
	}{
		{
			name: "int to string",
			build: func() any {
				return zcache.NewRedisCache[int, string](zcache.WithPrefix[int, string]("int-string"))
			},
		},
		{
			name: "string to struct",
			build: func() any {
				type customValue struct{ Data string }
				return zcache.NewRedisCache[string, customValue](zcache.WithPrefix[string, customValue]("string-struct"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.build()
			if c == nil {
				t.Errorf("NewRedisCache() for %s returned nil", tt.name)
			}
		})
	}
}

func TestRedisCache_ClearPrefixIsolation(t *testing.T) {
	ctx := t.Context()
	client := newTestRedisClient(t)

	// two caches with different prefixes sharing the same redis
	cacheA := zcache.NewRedisCache[string, int](
		zcache.WithPrefix[string, int]("isolation-a:"),
		zcache.WithClient[string, int](client),
	)
	cacheB := zcache.NewRedisCache[string, int](
		zcache.WithPrefix[string, int]("isolation-b:"),
		zcache.WithClient[string, int](client),
	)

	// seed both caches
	if err := cacheA.Set(ctx, "x", 1); err != nil {
		t.Fatalf("Set cacheA: %v", err)
	}
	if err := cacheB.Set(ctx, "y", 2); err != nil {
		t.Fatalf("Set cacheB: %v", err)
	}

	// clear only cacheA
	if err := cacheA.Clear(ctx); err != nil {
		t.Fatalf("Clear cacheA: %v", err)
	}

	// cacheA should be empty
	lenA, err := cacheA.Len(ctx)
	if err != nil {
		t.Fatalf("Len cacheA: %v", err)
	}
	if lenA != 0 {
		t.Errorf("cacheA Len() = %d after Clear, want 0", lenA)
	}

	// cacheB should still have its entry
	val, err := cacheB.Get(ctx, "y")
	if err != nil {
		t.Fatalf("Get cacheB: %v", err)
	}
	if val != 2 {
		t.Errorf("cacheB Get(y) = %d, want 2", val)
	}

	// cleanup
	cacheB.Clear(ctx)
}
