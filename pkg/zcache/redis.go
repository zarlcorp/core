package zcache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zarlcorp/core/pkg/zoptions"
)

var (
	_ Reader[string, any]     = (*RedisCache[string, any])(nil)
	_ Writer[string, any]     = (*RedisCache[string, any])(nil)
	_ ReadWriter[string, any] = (*RedisCache[string, any])(nil)
	_ Cache[string, any]      = (*RedisCache[string, any])(nil)
)

// RedisCache is a thread-safe cache implementation using Redis as the backend.
// It provides distributed caching capabilities across multiple application instances.
type RedisCache[K comparable, V any] struct {
	client redis.UniversalClient
	prefix string
	ttl    time.Duration
}

// RedisOption configures a RedisCache.
type RedisOption[K comparable, V any] = zoptions.Option[RedisCache[K, V]]

// WithClient sets the Redis client for the cache.
func WithClient[K comparable, V any](c redis.UniversalClient) RedisOption[K, V] {
	return func(rc *RedisCache[K, V]) {
		rc.client = c
	}
}

// WithPrefix sets the key prefix for all cache entries.
func WithPrefix[K comparable, V any](pre string) RedisOption[K, V] {
	return func(rc *RedisCache[K, V]) {
		rc.prefix = pre
	}
}

// WithTTL sets the time-to-live for cache entries.
func WithTTL[K comparable, V any](ttl time.Duration) RedisOption[K, V] {
	return func(rc *RedisCache[K, V]) {
		rc.ttl = ttl
	}
}

// NewRedisCache creates a new Redis-backed cache with the specified configuration.
func NewRedisCache[K comparable, V any](opts ...RedisOption[K, V]) *RedisCache[K, V] {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	rc := RedisCache[K, V]{
		client: client,
	}
	for _, opt := range opts {
		opt(&rc)
	}
	return &rc
}

// Set stores a key-value pair in Redis.
// If the key already exists, its value is updated.
func (c *RedisCache[K, V]) Set(ctx context.Context, key K, value V) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	redisKey := c.makeKey(key)

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, redisKey, data, c.ttl).Err()
}

// Get retrieves the value associated with the given key from Redis.
// Returns ErrNotFound if the key does not exist.
func (c *RedisCache[K, V]) Get(ctx context.Context, key K) (V, error) {
	select {
	case <-ctx.Done():
		var zero V
		return zero, ctx.Err()
	default:
	}

	redisKey := c.makeKey(key)

	result, err := c.client.Get(ctx, redisKey).Result()
	if err != nil {
		var zero V
		if errors.Is(err, redis.Nil) {
			return zero, ErrNotFound
		}
		return zero, err
	}

	var value V
	if err := json.Unmarshal([]byte(result), &value); err != nil {
		var zero V
		return zero, err
	}

	return value, nil
}

// Delete removes a key-value pair from Redis.
// Returns true if the key existed and was deleted, false otherwise.
func (c *RedisCache[K, V]) Delete(ctx context.Context, key K) (bool, error) {
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
	}

	redisKey := c.makeKey(key)

	result, err := c.client.Del(ctx, redisKey).Result()
	if err != nil {
		return false, err
	}

	return result > 0, nil
}

// Clear removes all entries with the configured prefix from Redis.
func (c *RedisCache[K, V]) Clear(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	pattern := c.prefix + "*"

	iter := c.client.Scan(ctx, 0, pattern, 0).Iterator()
	keys := make([]string, 0)

	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return err
	}

	if len(keys) > 0 {
		return c.client.Del(ctx, keys...).Err()
	}

	return nil
}

// Len returns the approximate number of entries with the configured prefix in Redis.
func (c *RedisCache[K, V]) Len(ctx context.Context) (int, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}

	pattern := c.prefix + "*"

	iter := c.client.Scan(ctx, 0, pattern, 0).Iterator()
	count := 0

	for iter.Next(ctx) {
		count++
	}

	if err := iter.Err(); err != nil {
		return 0, err
	}

	return count, nil
}

func (c *RedisCache[K, V]) makeKey(key K) string {
	keyBytes, _ := json.Marshal(key)
	return c.prefix + string(keyBytes)
}

// Healthy checks if Redis is accessible by pinging it.
func (c *RedisCache[K, V]) Healthy() error {
	ctx := context.Background()
	if err := c.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("ping redis: %w", err)
	}
	return nil
}
