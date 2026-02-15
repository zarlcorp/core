package zcache_test

import (
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/zarlcorp/core/pkg/zcache"
)

func TestRedisCache_Constructor(t *testing.T) {
	tests := []struct {
		name  string
		opts  []zcache.RedisOption[string, int]
		check func(t *testing.T, c any)
	}{
		{
			name: "with default client",
			check: func(t *testing.T, c any) {
				if c == nil {
					t.Error("NewRedisCache() returned nil")
				}
			},
		},
		{
			name: "with custom client",
			opts: []zcache.RedisOption[string, int]{
				zcache.WithPrefix[string, int]("test"),
				zcache.WithClient[string, int](&redis.Client{}),
			},
			check: func(t *testing.T, c any) {
				if c == nil {
					t.Error("NewRedisCache() with custom client returned nil")
				}
			},
		},
		{
			name: "with prefix",
			opts: []zcache.RedisOption[string, int]{
				zcache.WithPrefix[string, int]("test"),
			},
			check: func(t *testing.T, c any) {
				if c == nil {
					t.Error("NewRedisCache() with prefix returned nil")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := zcache.NewRedisCache[string, int](tt.opts...)
			tt.check(t, c)
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
