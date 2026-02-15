# 092: Deduplicate ErrNotFound sentinels across zsync and zcache

## Objective
Eliminate the duplicate `ErrNotFound` definitions in zsync and zcache so consumers don't accidentally check the wrong sentinel.

## Context
Both `zsync` and `zcache` define `var ErrNotFound = errors.New("key not found")` with identical messages. `errors.Is(err, zcache.ErrNotFound)` won't match `zsync.ErrNotFound` since they're different pointers. This is confusing for consumers who use both packages.

After spec 085 lands, zsync.ErrNotFound is removed (ZMap.Get returns bool). This spec handles any remaining cleanup.

## Requirements
- Verify that after 085 lands, `zsync.ErrNotFound` is fully removed and no consumers reference it
- If `ErrNotFound` is still needed in zsync for queue operations or future use, rename it to something specific (e.g. `ErrQueueNotFound`) to avoid confusion
- Ensure `zcache.ErrNotFound` is the single canonical "key not found" error for cache operations
- Grep across all zarlcorp repos (core, zburn, zvault, zshield) for any imports of `zsync.ErrNotFound` and update them
- All tests pass

## Target Repo
zarlcorp/core

## Agent Role
backend

## Files to Modify
- pkg/zsync/zsync.go (remove or rename ErrNotFound)
- Any consumers across repos that reference zsync.ErrNotFound

## Notes
Depends on 085 (which removes ErrNotFound from ZMap.Get). May be trivial if 085 already removes it entirely.
