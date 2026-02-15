package zsync

import (
	"sync"
)

// ZMap is a thread-safe generic map implementation.
// It provides concurrent access to a map with read-write locking for performance.
type ZMap[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V
}

// NewZMap creates a new thread-safe map with the specified key and value types.
func NewZMap[K comparable, V any]() *ZMap[K, V] {
	return &ZMap[K, V]{
		data: make(map[K]V),
	}
}

// Set stores a key-value pair in the map.
// If the key already exists, its value is updated.
func (m *ZMap[K, V]) Set(key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
}

// Get retrieves the value associated with the given key.
// Returns the value and true if found, or the zero value and false if not.
func (m *ZMap[K, V]) Get(key K) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	value, ok := m.data[key]
	return value, ok
}

// Delete removes a key-value pair from the map.
// Returns true if the key existed and was deleted, false otherwise.
func (m *ZMap[K, V]) Delete(key K) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, exists := m.data[key]
	if exists {
		delete(m.data, key)
	}
	return exists
}

// Len returns the number of key-value pairs in the map.
func (m *ZMap[K, V]) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.data)
}

// Keys returns a slice containing all keys in the map.
// The order of keys is not guaranteed.
func (m *ZMap[K, V]) Keys() []K {
	m.mu.RLock()
	defer m.mu.RUnlock()

	keys := make([]K, 0, len(m.data))
	for k := range m.data {
		keys = append(keys, k)
	}
	return keys
}

// Clear removes all key-value pairs from the map.
func (m *ZMap[K, V]) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data = make(map[K]V)
}
