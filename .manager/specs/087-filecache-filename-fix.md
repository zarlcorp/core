# 087: FileCache filename collision fix + remove MemFS type assertions

## Objective
Fix FileCache.makeFilename to prevent key collisions, handle marshal errors, and remove the concrete MemFS type assertions that break the FileSystem abstraction.

## Context
Review of pkg/ shared libraries identified:
1. `makeFilename` replaces multiple chars (`/`, `:`, `*`, etc.) with `_`, causing different keys to map to the same file. Example: keys `"a/b"` and `"a:b"` both become `"_a_b_.cache"`.
2. `json.Marshal(key)` error is silently discarded with `_`.
3. `Clear` and `Len` type-assert `c.fs.(*zfilesystem.MemFS)` for fast paths, breaking the abstraction.

## Requirements
- Replace char-replacement sanitization with hex encoding of the marshaled key bytes
  - Use `hex.EncodeToString(keyBytes)` + `.cache` suffix
  - This guarantees unique filenames for unique keys
- `makeFilename` must return an error if `json.Marshal` fails; propagate to callers
- Remove the `*zfilesystem.MemFS` type assertions from `Clear` and `Len`
  - Use the WalkDir-based implementation for all FileSystem backends
  - The MemFS-specific `ClearCacheFiles` and `CountCacheFiles` methods can stay on MemFS but FileCache must not know about them
- All existing tests must pass
- Add test: two keys that previously collided now produce different filenames

## Target Repo
zarlcorp/core

## Agent Role
backend

## Files to Modify
- pkg/zcache/file.go (makeFilename, Clear, Len, Set, Get, Delete)
- pkg/zcache/file_test.go

## Notes
Depends on 085 (zsync API changes) and 086 (MemFS WalkDir fixes, since FileCache uses WalkDir).
