# 086: MemFS WalkDir fixes — root filtering, lexical order, directory entries

## Objective
Fix MemFS.WalkDir to match OSFileSystem behavior: filter by root prefix, return entries in lexical order, and emit directory entries.

## Context
Review of pkg/ shared libraries identified three behavioral differences between MemFS and OSFileSystem WalkDir:
1. MemFS ignores the `root` parameter — walks all files regardless of root
2. MemFS returns entries in arbitrary map order; OSFileSystem returns lexical order
3. MemFS has no directory concept — MkdirAll is a no-op, WalkDir never yields directory entries

These differences can hide bugs in tests that use MemFS as a stand-in for OSFileSystem.

## Requirements
- `WalkDir(root, fn)` must only yield entries whose paths have `root` as a prefix
- Entries must be yielded in lexical (sorted) order, matching `filepath.WalkDir` behavior
- When files exist at paths like `"a/b/c.txt"`, WalkDir must synthesize and yield directory entries for `"a"` and `"a/b"` before the file entry
- The root directory itself should be yielded as the first entry (matching filepath.WalkDir)
- `MkdirAll` should record directories so they appear in WalkDir even if empty
- All existing contract tests must pass
- Add tests that verify:
  - Root filtering works (files outside root are excluded)
  - Lexical ordering matches OSFileSystem behavior
  - Directory entries appear in walk results

## Target Repo
zarlcorp/core

## Agent Role
backend

## Files to Modify
- pkg/zfilesystem/memfs.go (WalkDir, MkdirAll, internal directory tracking)
- pkg/zfilesystem/filesystem_test.go (contract tests for new behavior)

## Notes
Depends on 085 (zsync API changes) since MemFS uses ZMap internally.
