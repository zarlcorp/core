# 090: zapp Track-after-Close error + remove dead name field

## Objective
Make App.Track return an error when called after Close, and remove the unused `name` field and `WithName` option.

## Context
Review of pkg/ shared libraries identified:
1. Calling `Track(c)` after `Close()` silently appends the closer but it never gets called — the resource leaks with no indication.
2. `App.name` is set via `WithName` and defaults to `filepath.Base(os.Args[0])`, but nothing in the package reads or exposes it. Dead code.

## Requirements
- Change `Track(c io.Closer)` to `Track(c io.Closer) error`
- Track returns an error after Close has been called (define a sentinel: `var ErrClosed = errors.New("app closed")`)
- Remove the `name` field from `App` struct
- Remove `WithName` from options.go
- Update package doc example to not use `WithName`
- All existing tests must pass
- Add tests:
  - Track after Close returns ErrClosed
  - The tracked closer is NOT added to the cleanup list

## Target Repo
zarlcorp/core

## Agent Role
backend

## Files to Modify
- pkg/zapp/zapp.go (Track signature, remove name field, add ErrClosed)
- pkg/zapp/options.go (remove WithName)
- pkg/zapp/zapp_test.go

## Notes
Independent — no dependency on other specs. Small scope.
