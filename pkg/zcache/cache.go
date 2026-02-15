// Package zcache provides thread-safe concurrent data structures for caching.
//
// This package offers high-performance, generic implementations of cache operations
// with built-in synchronization, eliminating the need for external locking when accessing
// shared data across multiple goroutines.
//
// # Interfaces
//
// Reader: Basic read operations (Get, Len)
// Writer: Basic write operations (Set, Delete, Clear)
// ReadWriter: Combines Reader and Writer
// Cache: Complete cache interface
//
// # Error Handling
//
// All operations follow consistent error handling patterns:
//   - ErrNotFound: Returned when attempting to access non-existent keys
//   - Boolean returns: Indicate success/failure for operations like Delete
//   - No panics: All error conditions are handled gracefully
//
// # Usage Example
//
//	ctx := context.Background()
//	c := zcache.NewMemoryCache[string, int]()
//	c.Set(ctx, "key", 42)
//
//	value, err := c.Get(ctx, "key")
//	if err != nil {
//	    // handle ErrNotFound or context.Canceled
//	}
//
//	existed, err := c.Delete(ctx, "key")
package zcache

import (
	"context"
	"errors"
)

var (
	ErrNotFound = errors.New("key not found")
)

// Reader defines basic read operations for cache implementations.
type Reader[K comparable, V any] interface {
	// Get retrieves the value associated with the given key.
	// Returns ErrNotFound if the key does not exist.
	// Returns context.Canceled if the context is canceled.
	Get(ctx context.Context, key K) (V, error)

	// Len returns the number of entries in the cache.
	// Returns context.Canceled if the context is canceled.
	Len(ctx context.Context) (int, error)
}

// Writer defines basic write operations for cache implementations.
type Writer[K comparable, V any] interface {
	// Set stores a key-value pair in the cache.
	// If the key already exists, its value is updated.
	// Returns context.Canceled if the context is canceled.
	Set(ctx context.Context, key K, value V) error

	// Delete removes a key-value pair from the cache.
	// Returns true if the key existed and was deleted, false otherwise.
	// Returns context.Canceled if the context is canceled.
	Delete(ctx context.Context, key K) (bool, error)

	// Clear removes all entries from the cache.
	// Returns context.Canceled if the context is canceled.
	Clear(ctx context.Context) error
}

// ReadWriter combines basic read and write operations.
type ReadWriter[K comparable, V any] interface {
	Reader[K, V]
	Writer[K, V]
}

// Cache combines all common cache operations.
type Cache[K comparable, V any] interface {
	ReadWriter[K, V]
	Healthy() error
}
