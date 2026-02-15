# 088: RedisCache safety — prefix guard, batch Clear, marshal error handling

## Objective
Prevent RedisCache.Clear from nuking all Redis keys when no prefix is set, batch key deletion during Clear, and handle marshal errors in makeKey.

## Context
Review of pkg/ shared libraries identified:
1. `Clear` with empty prefix produces pattern `"*"` which matches ALL keys in Redis — deleting unrelated data.
2. `Clear` collects all matching keys into a single slice before deleting — memory bomb for large key sets.
3. `makeKey` silently ignores `json.Marshal` errors.

## Requirements
- Require prefix at construction time: `NewRedisCache` must require `WithPrefix` or set a sensible default prefix (e.g. `"zcache:"`)
  - If no prefix is provided, use a default like `"zcache:"` so bare `"*"` scans are impossible
  - Alternatively, `Clear` can return an error if prefix is empty — either approach is acceptable
- `Clear` must batch-delete keys as SCAN yields them (e.g. every 100 keys) instead of collecting all into memory
- `makeKey` must return `(string, error)` and propagate marshal failures to callers (Set, Get, Delete)
- `Len` uses the same SCAN pattern — document that it's approximate under concurrent access
- All existing tests must pass
- Add test: verify Clear with default prefix doesn't match keys outside the prefix

## Target Repo
zarlcorp/core

## Agent Role
backend

## Files to Modify
- pkg/zcache/redis.go (NewRedisCache, Clear, Len, makeKey, Set, Get, Delete)
- pkg/zcache/redis_test.go

## Notes
Independent of 085/086/087 — can run in parallel with other zcache work if the shared cache.go interfaces don't change.
