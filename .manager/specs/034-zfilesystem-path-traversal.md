# 034: zfilesystem — path traversal validation on OS implementation

## Objective
Add path traversal protection to `OSFileSystem` so callers cannot escape the configured `baseDir` using `../` sequences or absolute paths.

## Context
The `OSFileSystem` wraps `os` package calls with `filepath.Join(baseDir, filename)`. Currently no validation ensures the resolved path stays within `baseDir`. A caller passing `../../../etc/passwd` would escape the sandbox. This is a security prerequisite before any zarlcorp tool (zburn, zvault) stores user data through this interface.

Related: spec 031 (package review) identified this as a known gap. Issue #20.

## Requirements

### Path validation
- Every method that accepts a path (`ReadFile`, `WriteFile`, `Remove`, `MkdirAll`, `OpenFile`, `WalkDir`) must validate the resolved path stays within `baseDir`.
- Resolve the joined path with `filepath.Abs` (or `filepath.Clean` + prefix check) and confirm it has `baseDir` as a prefix.
- Reject absolute paths in the input (e.g. `/etc/passwd`) — all paths must be relative to `baseDir`.
- Return a clear error on rejection: `"path escapes base directory"` or similar.
- Extract the validation into a single unexported helper (`resolvePath` or `cleanPath`) called by every method — no duplication.

### MemFS parity
- `MemFS` doesn't have the same security concern (no real filesystem access) but should reject `..` components for API consistency. Clean paths with `filepath.Clean` and reject anything that resolves outside the logical root.

### Contract tests
- Add path traversal test cases to the existing contract test suite so both implementations are verified:
  - `../escape` — rejected
  - `foo/../../escape` — rejected
  - `/absolute/path` — rejected
  - `foo/../bar` (resolves within base) — allowed
  - `./normal/path` — allowed
  - Empty string — defined behavior (either reject or treat as root)

### No API changes
- No new exported types or functions. The validation is internal to existing methods.

## Target Repo
zarlcorp/core

## Agent Role
backend

## Files to Modify
- `pkg/zfilesystem/osfs.go` — add path validation helper, call from all methods
- `pkg/zfilesystem/memfs.go` — add path cleaning/validation for consistency
- `pkg/zfilesystem/filesystem_test.go` — add path traversal contract tests

## Notes
- `filepath.Join` already resolves `..` but doesn't prevent escaping. The fix is checking the result, not changing the join.
- The `BaseDir()` accessor on `OSFileSystem` remains unchanged.
- Keep the error message direct per coding standards: `"path escapes base directory: %s"`, not `"failed to validate path"`.
