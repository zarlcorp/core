# 052: First zburn release

## Objective
Tag and release zburn v0.1.0 — the first downloadable binary.

## Context
zburn has all core features: identity generator, encrypted store, CLI subcommands, and interactive TUI. The release workflow (spec 049) is in place on zarlcorp/.github. Time to ship.

## Requirements

### Pre-release checklist
- Verify `go build ./...` succeeds
- Verify `go test ./...` passes
- Verify the release workflow exists in zburn's `.github/workflows/release.yml` calling `zarlcorp/.github/.github/workflows/go-release.yml@main`
- If the release workflow doesn't exist, create it

### Release workflow file (if missing)
Create `.github/workflows/release.yml`:
```yaml
name: Release
on:
  push:
    tags: ["v*"]

jobs:
  release:
    uses: zarlcorp/.github/.github/workflows/go-release.yml@main
    with:
      binary-name: zburn
```

### Tag and release
- Tag `v0.1.0` on main
- Push the tag to trigger the release workflow
- Verify the GitHub Release is created with binaries for linux/darwin amd64/arm64

### Homebrew formula update
- Update the formula in `zarlcorp/homebrew-tap` to point to the v0.1.0 release
- Verify `brew install zarlcorp/tap/zburn` works (or at least that the formula is valid)

## Target Repo
zarlcorp/zburn (primary), zarlcorp/homebrew-tap (formula update)

## Agent Role
devops

## Dependencies
- 050 (README update should be merged before release)

## Notes
- This is a manager-driven task — tagging, pushing, and verifying are git operations
- The release workflow does the actual build
- The homebrew formula needs the release URL and SHA256 checksums from the GitHub Release
- This may need to be done partially manually since it involves waiting for CI
