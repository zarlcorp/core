# 032: CI pipeline — dynamic module detection, lint, coverage

## Objective
Create a GitHub Actions CI pipeline that dynamically detects which Go modules changed in a PR, runs build/test/lint/coverage only for those modules, and enforces minimum code coverage. Single workflow, smart about what to run.

## Context
zarlcorp/core is a Go workspace with 7 modules under `pkg/`. There's currently no CI at all — no .github directory. The pipeline needs to be workspace-aware: a PR touching `pkg/zsync/` should only test zsync, not all 7 modules.

## Requirements

### Workflow file — `.github/workflows/ci.yml`

Triggers:
- `pull_request` targeting `main`
- `push` to `main` (post-merge validation)

### Job 1: Detect changed modules

- Use `dorny/paths-filter` or git diff to determine which `pkg/*` directories have changes
- Output a JSON matrix of changed module paths
- If go.work or .github/ changed, run all modules
- If no modules changed (e.g. docs-only PR), skip the test job

### Job 2: Test (matrix over changed modules)

For each changed module:
1. **Setup Go** — use `actions/setup-go@v5` with Go 1.26
2. **Build** — `go build ./...` from module directory
3. **Vet** — `go vet ./...`
4. **Test with coverage** — `go test -race -coverprofile=coverage.out -covermode=atomic ./...`
5. **Check coverage** — parse coverage.out and fail if below threshold (80%)
6. **Upload coverage** — upload as artifact for visibility

### Job 3: Lint (runs once, whole workspace)

- Use `golangci/golangci-lint-action`
- Run from workspace root
- Use a `.golangci.yml` config file with sensible defaults

### golangci-lint config — `.golangci.yml`

Linters to enable beyond defaults:
- `govet` — correctness
- `errcheck` — unchecked errors
- `staticcheck` — bugs and simplifications
- `unused` — unused code
- `gosimple` — simplifications
- `ineffassign` — ineffectual assignments
- `gocritic` — opinionated style
- `revive` — flexible linter (replaces golint)
- `misspell` — typos in comments/strings
- `gofumpt` — stricter formatting

Linters to explicitly disable:
- `wsl` — whitespace linting (too noisy)
- `nlreturn` — newline before return (too opinionated)

Settings:
- Line length: 120
- Test files included in lint

### Coverage enforcement

The coverage check should:
- Parse the coverage.out file
- Extract the total coverage percentage
- Fail the job if coverage < 80%
- Print the coverage percentage in job output
- Use `go tool cover -func=coverage.out` and parse the total line

### What the pipeline does NOT do
- No release automation (separate spec)
- No deployment
- No Docker builds
- No caching beyond what setup-go provides (keep it simple)

## Acceptance Criteria
1. `.github/workflows/ci.yml` exists and is valid
2. `.golangci.yml` exists with configured linters
3. Changed-module detection works (only tests affected packages)
4. go.work or .github changes trigger all modules
5. Docs-only PRs skip testing
6. Tests run with race detector and coverage
7. Coverage below 80% fails the job
8. golangci-lint runs on the full workspace
9. Pipeline passes on current codebase

## Target Repo
zarlcorp/core

## Agent Role
backend

## Files to Create
- `.github/workflows/ci.yml`
- `.golangci.yml`

## Dependencies
Depends on 029 (Go upgrade) — workflow should target Go 1.26.

## Notes
- Go workspace mode: the lint job should run at workspace root to lint all code consistently
- The test job runs per-module because each module has its own go.mod and test suite
- Coverage threshold of 80% is a starting point — can be adjusted per-module later
- `dorny/paths-filter` is well-maintained and handles the matrix generation cleanly
- For the coverage check, a simple shell script parsing `go tool cover -func` output is sufficient — no need for external coverage tools
- The pipeline should be fast — parallel matrix jobs for independent modules
