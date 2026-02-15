# 079: Fix WalkDir returning absolute paths

## Objective
`OSFileSystem.WalkDir` passes absolute paths to the callback function, but callers expect relative paths (relative to the filesystem's base directory). This breaks any code that checks path prefixes or passes callback paths back to other filesystem methods.

## Context
The store's `List()` method walks `identities/` and checks `strings.HasPrefix(path, "identities/")`. But `WalkDir` resolves the root to an absolute path before calling `filepath.WalkDir`, so the callback receives paths like `/Users/x/.local/share/zburn/identities/abc.enc` instead of `identities/abc.enc`. The prefix check always fails, so List always returns empty — even though Save writes files correctly (using relative paths).

This is in core's `pkg/zfilesystem`, so it affects all tools that use WalkDir, not just zburn.

## Requirements

### Fix in `pkg/zfilesystem/osfs.go`

The `WalkDir` method needs to strip the base directory prefix from paths before passing them to the callback. The callback should receive paths relative to the filesystem's base directory.

Current code:
```go
func (o *OSFileSystem) WalkDir(root string, fn fs.WalkDirFunc) error {
    p, err := o.resolvePath(root)
    if err != nil {
        return err
    }
    return filepath.WalkDir(p, fn)
}
```

Fix: wrap the callback to convert absolute paths to relative:
```go
func (o *OSFileSystem) WalkDir(root string, fn fs.WalkDirFunc) error {
    p, err := o.resolvePath(root)
    if err != nil {
        return err
    }
    return filepath.WalkDir(p, func(path string, d fs.DirEntry, err error) error {
        rel, relErr := filepath.Rel(o.baseDir, path)
        if relErr != nil {
            return relErr
        }
        return fn(rel, d, err)
    })
}
```

### Tests

Update or add tests in `pkg/zfilesystem/osfs_test.go` that verify:
1. WalkDir callback receives relative paths (not absolute)
2. The relative paths can be passed back to ReadFile/WriteFile without error
3. Walking a subdirectory returns paths relative to baseDir (e.g., `identities/foo.enc`, not just `foo.enc`)

Run `go test ./pkg/zfilesystem/...` to verify.

## Target Repo
zarlcorp/core

## Agent Role
backend

## Files to Modify
- pkg/zfilesystem/osfs.go
- pkg/zfilesystem/osfs_test.go

## Notes
This is the critical fix — it's the root cause of zburn's "save works but list is empty" bug. After this is tagged and zburn updates the dependency, the store will work correctly. No changes needed in the store itself — the paths will just be correct.
