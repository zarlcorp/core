# 089: MemoryCache TTL support via functional option

## Objective
Add optional TTL (time-to-live) support to MemoryCache via a variadic functional option, so entries expire automatically.

## Context
Review of pkg/ shared libraries noted that MemoryCache grows unbounded with no eviction. For long-running TUI apps this could become a problem. The user wants a `WithTTL` option following the existing zoptions pattern.

## Requirements
- Add `type MemoryOption[K, V] = zoptions.Option[MemoryCache[K, V]]`
- Add `func WithMemoryTTL[K comparable, V any](ttl time.Duration) MemoryOption[K, V]`
- Update `NewMemoryCache` to accept variadic options: `func NewMemoryCache[K comparable, V any](opts ...MemoryOption[K, V]) *MemoryCache[K, V]`
- When TTL is set:
  - Each entry stores its insertion time
  - `Get` checks if the entry has expired; if so, deletes it and returns `ErrNotFound`
  - Expired entries are lazily cleaned up on access (no background goroutine needed)
- When TTL is zero (default), behavior is unchanged — no expiry
- Do NOT add max-size eviction — keep it simple, TTL only
- All existing tests must pass (zero-arg NewMemoryCache still works)
- Add tests:
  - Entry expires after TTL
  - Entry accessible before TTL
  - Zero TTL means no expiry
  - Len excludes expired entries

## Target Repo
zarlcorp/core

## Agent Role
backend

## Files to Modify
- pkg/zcache/memory.go (add TTL field, wrap entries, lazy expiry on Get)
- pkg/zcache/memory_test.go
- pkg/zcache/cache_test.go (contract tests should still pass)

## Notes
Independent — no dependency on other specs.
