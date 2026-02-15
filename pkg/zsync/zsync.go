// Package zsync provides thread-safe concurrent data structures for Go applications.
//
// This package offers high-performance, generic implementations of common data structures
// with built-in synchronization, eliminating the need for external locking when accessing
// shared data across multiple goroutines.
//
// # Data Structures
//
// ZMap: A thread-safe generic map implementation using read-write locks for optimal
// performance. Supports all standard map operations with concurrent safety.
//
// ZSet: A thread-safe generic set implementation built on top of ZMap, providing
// efficient operations for managing collections of unique values.
//
// ZQueue: A thread-safe generic FIFO queue implementation with blocking operations
// for producer/consumer scenarios. Supports context cancellation and graceful shutdown.
//
// # Error Handling
//
// All operations follow consistent error handling patterns:
//   - ErrNotFound: Returned when attempting to access non-existent keys/values
//   - ErrQueueClosed: Returned when attempting to operate on a closed queue
//   - ErrQueueEmpty: Returned when attempting to pop from an empty queue
//   - ErrCanceled: Returned when operations are canceled via context
//   - Boolean returns: Indicate success/failure for operations like Delete/Remove
//   - No panics: All error conditions are handled gracefully
//
// # Performance Characteristics
//
// - Read operations use RLock for concurrent access
// - Write operations use exclusive Lock for data integrity
// - Memory efficient: ZSet uses struct{} values to minimize overhead
// - Pre-allocated slices in bulk operations to reduce allocations
//
// # Usage Examples
//
// Basic ZMap usage:
//
//	m := zsync.NewZMap[string, int]()
//	m.Set("key", 42)
//
//	value, err := m.Get("key")
//	if err != nil {
//	    // handle ErrNotFound
//	}
//
//	existed := m.Delete("key")
//
// Basic ZSet usage:
//
//	s := zsync.NewZSet[string]()
//	s.Add("value")
//
//	if s.Contains("value") {
//	    s.Remove("value")
//	}
//
//	values := s.Values() // []string with all unique values
//
// Basic ZQueue usage:
//
//	q := zsync.NewZQueue[string]()
//	q.Push("item1")
//
//	item, err := q.Pop() // blocks until item available
//	if err != nil {
//	    // handle ErrQueueClosed
//	}
//
//	q.Close() // signal shutdown to waiting consumers
//
// # Thread Safety
//
// All data structures in this package are fully thread-safe and can be used
// concurrently from multiple goroutines without additional synchronization.
// The implementations use efficient read-write locks to maximize concurrent
// read performance while ensuring data consistency during writes.
package zsync

import "errors"

// Package-level errors used across multiple data structures
var (
	ErrNotFound    = errors.New("key not found")
	ErrQueueClosed = errors.New("queue closed")
	ErrQueueEmpty  = errors.New("queue empty")
	ErrCanceled    = errors.New("operation canceled")
)
