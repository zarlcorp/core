# 093: zcache Redis tests with miniredis

## Objective
Replace the `-short` skip pattern in Redis cache tests with `alicebob/miniredis/v2`, an in-memory Redis server that runs in-process. All Redis tests should run in CI and locally without a real Redis.

## Context
The zcache contract tests and `TestRedisCache_ClearPrefixIsolation` currently skip when `-short` is set and fail when no Redis is running. This means Redis cache behavior is never tested in CI. Using miniredis gives us a real Redis protocol implementation in-memory — no mocks, no skips.

## Requirements
- Add `github.com/alicebob/miniredis/v2` as a test dependency in `pkg/zcache/go.mod`
- Create a test helper that starts a miniredis server and returns a `redis.Client` pointed at it:
  ```go
  func newTestRedisClient(t *testing.T) *redis.Client {
      s := miniredis.RunT(t) // auto-closed on test cleanup
      return redis.NewClient(&redis.Options{Addr: s.Addr()})
  }
  ```
- Update `cache_test.go` contract test for RedisCache:
  - Remove the `testing.Short()` skip
  - Use `WithClient` with the miniredis-backed client
- Update `redis_test.go`:
  - `TestRedisCache_ClearPrefixIsolation`: remove the `testing.Short()` skip, use miniredis client for both caches
- All Redis tests must now pass without `-short` and without a real Redis server
- Run `go test ./...` (not `-short`) to verify everything passes
- Run `go mod tidy` after adding the dependency

## Target Repo
zarlcorp/core

## Agent Role
testing

## Files to Modify
- pkg/zcache/go.mod (add miniredis dependency)
- pkg/zcache/cache_test.go (use miniredis in contract test)
- pkg/zcache/redis_test.go (use miniredis in isolation test, remove short skips)

## Notes
Single item, no dependencies. Small scope.
