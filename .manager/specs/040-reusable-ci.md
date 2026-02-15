# 040: Reusable CI workflows for Go tool repos

## Objective
Create a `zarlcorp/.github` repository containing reusable GitHub Actions workflows that all zarlcorp tool repos (zburn, zvault, zshield) reference instead of duplicating CI configuration.

## Context
The CI pipeline in core (spec 032) is purpose-built for a multi-module workspace. Tool repos are single-module and need a simpler, standardized workflow. Rather than copy-pasting CI config into every tool repo, we create reusable workflows in the org-level `.github` repo.

The manifesto states: "Reusable CI workflows live in `zarlcorp/.github` — individual tool repos reference them, not copy them."

Depends on: spec 032 (CI pipeline — establishes patterns to reuse).
Issue #26.

## Requirements

### Repository: zarlcorp/.github

### Workflow: go-ci.yml
Reusable workflow (`workflow_call`) for Go projects:

**Inputs:**
- `go-version` (string, default: `"1.26"`) — Go version to use.
- `coverage-threshold` (number, default: `80`) — minimum coverage percentage.

**Jobs:**
1. **build** — `go build ./...`
2. **test** — `go test -race -coverprofile=coverage.out -covermode=atomic ./...`
3. **coverage** — parse coverage, fail if below threshold.
4. **lint** — `golangci-lint run` using the repo's `.golangci.yml` (or a sensible default).

### Workflow: go-release.yml
Reusable workflow for releasing Go binaries:

**Inputs:**
- `binary-name` (string, required) — name of the binary.

**Triggers:** called on tag push matching `v*`.

**Jobs:**
1. **goreleaser** — cross-compile (linux/darwin/windows, amd64/arm64) and publish to GitHub Releases.

### Caller example
Tool repos reference these with minimal config:

```yaml
# .github/workflows/ci.yml in zarlcorp/zburn
name: CI
on:
  pull_request:
    branches: [main]

jobs:
  ci:
    uses: zarlcorp/.github/.github/workflows/go-ci.yml@main
    with:
      go-version: "1.26"
      coverage-threshold: 80
```

### Shared lint config (optional)
If tool repos share the same `.golangci.yml`, include a default in `.github/` that tools can copy or symlink. But each repo can override — the workflow uses the repo's own config if present.

## Target Repo
zarlcorp/.github

## Agent Role
devops

## Files to Modify
All files are new — this is a new repository.
- `.github/workflows/go-ci.yml` — reusable CI workflow
- `.github/workflows/go-release.yml` — reusable release workflow
- `README.md` — usage docs for tool repo maintainers

## Notes
- The `.github` repo is special in GitHub — its workflows are automatically available to all repos in the org via `uses: zarlcorp/.github/.github/workflows/<name>@main`.
- GoReleaser config (`.goreleaser.yml`) lives in each tool repo, not here. The reusable workflow just runs it.
- Keep workflows simple. Tool repos are single-module, single-binary — no matrix, no module detection needed.
