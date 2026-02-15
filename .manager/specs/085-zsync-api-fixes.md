# 085: zsync API fixes — ZMap Get pattern + queue memory leak

## Objective
Fix the ZMap.Get API to use idiomatic `(value, ok)` pattern instead of error return, and fix the memory leak in ZQueue's Pop/TryPop/PopContext where removed elements are never zeroed.

## Context
Review of pkg/ shared libraries identified two issues in zsync:
1. `ZMap.Get` returns `(V, error)` with `ErrNotFound` — unusual for a map. Standard Go uses `(V, bool)`.
2. `ZQueue.Pop`, `TryPop`, and `PopContext` do `q.items = q.items[1:]` without zeroing the removed element, causing the backing array to hold references indefinitely.

This is a breaking change to ZMap's API. All consumers (ZSet, MemFS, FileCache) must be updated.

## Requirements
- Change `ZMap.Get(key K) (V, error)` to `ZMap.Get(key K) (V, bool)`
- Remove `ErrNotFound` from zsync package (consumers use zcache's if needed)
- Update `ZSet.Contains` to use the new `(V, bool)` return
- Update all internal consumers in zfilesystem and zcache that call `ZMap.Get`
- In `ZQueue.Pop`, `TryPop`, and `PopContext`: zero the element at index 0 before re-slicing
  ```go
  item := q.items[0]
  var zero T
  q.items[0] = zero
  q.items = q.items[1:]
  ```
- All existing tests must pass (update assertions for new Get signature)
- Add a test that demonstrates the queue doesn't hold references after Pop

## Target Repo
zarlcorp/core

## Agent Role
backend

## Files to Modify
- pkg/zsync/zsync.go (remove ErrNotFound)
- pkg/zsync/map.go (change Get signature)
- pkg/zsync/set.go (update Contains)
- pkg/zsync/queue.go (zero elements in Pop/TryPop/PopContext)
- pkg/zsync/map_test.go
- pkg/zsync/set_test.go
- pkg/zsync/queue_test.go
- pkg/zfilesystem/memfs.go (update Get callers)
- pkg/zcache/file.go (update if any Get calls exist)

## Notes
This must land before 086, 087, and 092 since those depend on the new API.
