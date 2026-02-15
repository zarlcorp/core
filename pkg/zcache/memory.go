package zcache

import (
	"context"
	"sync"
	"time"

	"github.com/zarlcorp/core/pkg/zoptions"
)

var (
	_ Reader[string, any]     = (*MemoryCache[string, any])(nil)
	_ Writer[string, any]     = (*MemoryCache[string, any])(nil)
	_ ReadWriter[string, any] = (*MemoryCache[string, any])(nil)
	_ Cache[string, any]      = (*MemoryCache[string, any])(nil)
)

// ttlEntry wraps a value with its insertion time for expiry tracking.
type ttlEntry[V any] struct {
	value     V
	createdAt time.Time
}

// MemoryOption configures a MemoryCache.
type MemoryOption[K comparable, V any] = zoptions.Option[MemoryCache[K, V]]

// WithMemoryTTL sets the time-to-live for cache entries.
// entries are lazily evicted on access when expired.
// zero (default) means no expiry.
func WithMemoryTTL[K comparable, V any](ttl time.Duration) MemoryOption[K, V] {
	return func(mc *MemoryCache[K, V]) {
		mc.ttl = ttl
	}
}

// MemoryCache is a thread-safe generic cache implementation.
// It provides concurrent access to a map with read-write locking for performance.
type MemoryCache[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V

	// ttl fields, only used when ttl > 0
	ttl     time.Duration
	ttlData map[K]ttlEntry[V]
	now     func() time.Time // injectable for testing
}

// NewMemoryCache creates a new thread-safe memory cache with the specified key and value types.
func NewMemoryCache[K comparable, V any](opts ...MemoryOption[K, V]) *MemoryCache[K, V] {
	mc := &MemoryCache[K, V]{
		now: time.Now,
	}
	for _, opt := range opts {
		opt(mc)
	}

	if mc.ttl > 0 {
		mc.ttlData = make(map[K]ttlEntry[V])
	} else {
		mc.data = make(map[K]V)
	}

	return mc
}

// Set stores a key-value pair in the cache.
// If the key already exists, its value is updated.
func (c *MemoryCache[K, V]) Set(ctx context.Context, key K, value V) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.ttl > 0 {
		c.ttlData[key] = ttlEntry[V]{value: value, createdAt: c.now()}
		return nil
	}

	c.data[key] = value
	return nil
}

// Get retrieves the value associated with the given key.
// Returns ErrNotFound if the key does not exist or has expired.
func (c *MemoryCache[K, V]) Get(ctx context.Context, key K) (V, error) {
	select {
	case <-ctx.Done():
		var zero V
		return zero, ctx.Err()
	default:
	}

	if c.ttl > 0 {
		return c.getTTL(key)
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	value, exists := c.data[key]
	if !exists {
		var zero V
		return zero, ErrNotFound
	}
	return value, nil
}

// getTTL handles Get with expiry checking.
// uses a write lock because it may need to delete expired entries.
func (c *MemoryCache[K, V]) getTTL(key K) (V, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, exists := c.ttlData[key]
	if !exists {
		var zero V
		return zero, ErrNotFound
	}

	if c.now().Sub(entry.createdAt) >= c.ttl {
		delete(c.ttlData, key)
		var zero V
		return zero, ErrNotFound
	}

	return entry.value, nil
}

// Delete removes a key-value pair from the cache.
// Returns true if the key existed and was deleted, false otherwise.
func (c *MemoryCache[K, V]) Delete(ctx context.Context, key K) (bool, error) {
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.ttl > 0 {
		_, exists := c.ttlData[key]
		if exists {
			delete(c.ttlData, key)
		}
		return exists, nil
	}

	_, exists := c.data[key]
	if exists {
		delete(c.data, key)
	}
	return exists, nil
}

// Len returns the number of entries in the cache.
// When TTL is set, expired entries are excluded from the count.
func (c *MemoryCache[K, V]) Len(ctx context.Context) (int, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}

	if c.ttl > 0 {
		return c.lenTTL(), nil
	}

	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.data), nil
}

// lenTTL counts non-expired entries and cleans up expired ones.
func (c *MemoryCache[K, V]) lenTTL() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := c.now()
	for k, entry := range c.ttlData {
		if now.Sub(entry.createdAt) >= c.ttl {
			delete(c.ttlData, k)
		}
	}
	return len(c.ttlData)
}

// Clear removes all entries from the cache.
func (c *MemoryCache[K, V]) Clear(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.ttl > 0 {
		c.ttlData = make(map[K]ttlEntry[V])
		return nil
	}

	c.data = make(map[K]V)
	return nil
}

// Healthy returns nil as memory cache is always healthy.
func (c *MemoryCache[K, V]) Healthy() error {
	return nil
}
