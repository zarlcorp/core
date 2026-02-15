# 049: Remove GoReleaser dependency

## Objective
Remove the GoReleaser-based release workflow from zarlcorp/.github and replace it with a simple `go build` + GitHub Releases approach.

## Context
The `go-release.yml` reusable workflow in `zarlcorp/.github` uses GoReleaser to build and publish binaries. GoReleaser is an unnecessary dependency — we can achieve the same result with `go build` and `gh release create`. The Manifesto says: "Every dependency has a cost."

## Requirements
- Remove `.github/workflows/go-release.yml` from zarlcorp/.github
- Replace with a simple workflow that:
  - Triggers on tag push (e.g. `v*`)
  - Takes `binary-name` as input
  - Builds binaries for linux/amd64, linux/arm64, darwin/amd64, darwin/arm64
  - Uses `go build -ldflags "-s -w -X main.version=$TAG"` for each platform
  - Creates a GitHub Release with the binaries attached using `gh release create`
- Update the README in zarlcorp/.github to reflect the new workflow
- Remove the placeholder Homebrew formula in zarlcorp/homebrew-tap that references GoReleaser checksums format

## Target Repo
zarlcorp/.github

## Agent Role
devops

## Files to Modify
- .github/workflows/go-release.yml (rewrite)
- README.md (update usage example)

## Notes
- No goreleaser.yml files exist in any tool repo yet — this is cleanup before first release
- The Homebrew formula in homebrew-tap may need updating too, but that's a separate concern once the release workflow is finalized
