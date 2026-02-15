# 029: Upgrade to Go 1.26 and clean up modules

## Objective
Upgrade the entire workspace to Go 1.26.0, fix the go.work file, and ensure all modules build cleanly. This unblocks everything else.

## Context
The workspace currently declares `go 1.24.4` across go.work and all go.mod files. Go 1.26.0 was released on 2026-02-10. The IDE/gopls is running Go 1.24.3, which causes diagnostic errors because go.work requires >= 1.24.4. Upgrading to 1.26 fixes the version mismatch and gives access to Go 1.25/1.26 features.

Two stub packages (zcrypto, zstyle) have go.mod files but no .go source files, making them invalid modules. These need a minimal `doc.go` placeholder so the workspace compiles.

## Requirements

### Install Go 1.26
- Install Go 1.26.0 via `go install golang.org/dl/go1.26.0@latest && go1.26.0 download` or homebrew
- Verify with `go version`

### Update go.work
```
go 1.26.0

use (
    ./pkg/zapp
    ./pkg/zcache
    ./pkg/zcrypto
    ./pkg/zfilesystem
    ./pkg/zoptions
    ./pkg/zstyle
    ./pkg/zsync
)
```

### Update all go.mod files
Set `go 1.26.0` in every module:
- `pkg/zapp/go.mod`
- `pkg/zcache/go.mod`
- `pkg/zcrypto/go.mod`
- `pkg/zfilesystem/go.mod`
- `pkg/zoptions/go.mod`
- `pkg/zstyle/go.mod`
- `pkg/zsync/go.mod`

### Fix stub modules
Add a minimal `doc.go` to each stub so they compile:

**pkg/zcrypto/doc.go:**
```go
// Package zcrypto provides encryption primitives for zarlcorp privacy tools.
package zcrypto
```

**pkg/zstyle/doc.go:**
```go
// Package zstyle provides the zarlcorp visual identity for TUIs.
package zstyle
```

### Run go mod tidy on all modules
Run `go mod tidy` in each package directory to update checksums and clean up dependencies.

### Verify
- `go build ./...` succeeds from workspace root
- `go test ./...` passes for all modules with existing tests
- `go vet ./...` reports no issues

## Acceptance Criteria
1. `go version` outputs go1.26.0
2. go.work declares `go 1.26.0`
3. All 7 go.mod files declare `go 1.26.0`
4. Stub modules have doc.go and compile
5. `go build ./...` succeeds at workspace root
6. All existing tests pass
7. No go vet warnings

## Target Repo
zarlcorp/core

## Agent Role
backend

## Files to Modify
- `go.work`
- `pkg/zapp/go.mod`
- `pkg/zcache/go.mod`
- `pkg/zcrypto/go.mod` + `pkg/zcrypto/doc.go` (create)
- `pkg/zfilesystem/go.mod`
- `pkg/zoptions/go.mod`
- `pkg/zstyle/go.mod`
- `pkg/zstyle/doc.go` (create if not exists)
- `pkg/zsync/go.mod`

## Dependencies
None — this is the first item. Everything else depends on this.

## Notes
- The go.work file does NOT have a module declaration — it's just the `go` directive and `use` blocks
- Stub doc.go files are temporary — they'll be replaced when the real implementation lands (028-zstyle, 031-zcrypto)
- Run `go mod tidy` in dependency order: zoptions, zsync first (no deps), then zfilesystem, then zcache, zapp, zcrypto, zstyle
