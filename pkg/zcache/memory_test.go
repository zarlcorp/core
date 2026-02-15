package zcache_test

import (
	"testing"

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
