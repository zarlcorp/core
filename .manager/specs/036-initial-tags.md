# 036: Tag initial v0.1.0 for all founding packages

## Objective
Create the initial `v0.1.0` version tag for each of the 7 founding packages, establishing the baseline for Go module versioning and enabling downstream tools to depend on tagged releases.

## Context
All 7 packages are built, tested, and reviewed. The CI pipeline (spec 032) validates them. The auto-tagging workflow (spec 033) will handle subsequent version bumps on merge. But the auto-tagger needs an existing tag to bump — the very first tags must be created manually.

Once these tags exist, `go get github.com/zarlcorp/core/pkg/zapp@v0.1.0` works. This unblocks tool repos (zburn, zvault, zshield) from depending on core packages.

Depends on: spec 032 (CI pipeline merged), spec 033 (auto-tagging merged).
Issue #22.

## Requirements

### Tags to create
Per Go multi-module convention, each tag includes the module path prefix:

```
pkg/zapp/v0.1.0
pkg/zcache/v0.1.0
pkg/zcrypto/v0.1.0
pkg/zfilesystem/v0.1.0
pkg/zoptions/v0.1.0
pkg/zstyle/v0.1.0
pkg/zsync/v0.1.0
```

### Process
- Tags are created on `main` at the commit where CI and auto-tagging are both merged and green.
- Push all 7 tags to the remote.
- Verify each module is resolvable: `GOPROXY=proxy.golang.org go list -m github.com/zarlcorp/core/pkg/<name>@v0.1.0`.

### No code changes
- This is a git-only operation. No files are modified.

## Target Repo
zarlcorp/core

## Agent Role
devops

## Files to Modify
None — git tags only.

## Notes
- This is a one-time bootstrap. After this, spec 033's auto-tagger handles all future version bumps.
- If the auto-tagger from spec 033 is already merged and can handle the initial v0.1.0 creation (first run with no existing tag → creates v0.1.0), then this spec reduces to "trigger the auto-tagger for all modules." Verify this before proceeding.
- The manager handles all git operations — this is not delegated to a sub-agent. It's a manager task.
