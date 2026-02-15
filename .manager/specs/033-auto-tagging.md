# 033: Automated version tagging on merge

## Objective
Create a GitHub Actions workflow that automatically tags affected modules with semantic version tags when a PR merges to main. Per-module tags following Go multi-module convention: `pkg/<name>/v0.x.y`.

## Context
zarlcorp/core uses per-module versioning — each package under `pkg/` is an independent Go module with its own go.mod. Go modules in a multi-module repo use path-prefixed tags: `pkg/zsync/v0.1.0`, `pkg/zcache/v0.2.0`, etc.

Currently there are no git tags at all. The first run should establish v0.1.0 for all existing modules.

## Requirements

### Workflow file — `.github/workflows/tag.yml`

Triggers:
- `push` to `main` (runs after PR merge)

### Logic

1. **Detect changed modules** — same approach as CI pipeline (git diff against previous commit or merge base)
2. **For each changed module:**
   - Find the latest tag matching `pkg/<name>/v*`
   - If no tag exists, use `v0.0.0` as baseline
   - Bump the patch version (v0.1.0 → v0.1.1)
   - Create and push the new tag
3. **Tag format**: `pkg/<name>/v<major>.<minor>.<patch>`

### Version bump strategy

- **Default**: patch bump (v0.1.0 → v0.1.1)
- **Minor bump**: if PR title or commit message contains `[minor]` or the PR has a `minor` label
- **Major bump**: if PR title or commit message contains `[major]` or the PR has a `major` label
- **Skip tagging**: if PR title or commit message contains `[skip-tag]` or `[no-tag]`

### Permissions
- The workflow needs `contents: write` permission to push tags
- Use the default `GITHUB_TOKEN` — no PAT needed for tag pushing

### What this workflow does NOT do
- No release notes generation
- No binary builds (no goreleaser)
- No changelog management
- No homebrew tap updates
- No notification

## Acceptance Criteria
1. `.github/workflows/tag.yml` exists and is valid
2. Merging a PR that touches `pkg/zsync/` creates a `pkg/zsync/vX.Y.Z` tag
3. Tags increment correctly from the previous tag
4. Multiple modules changed in one PR get individual tags
5. `[minor]` and `[major]` markers work for version bumps
6. `[skip-tag]` prevents tagging
7. First run creates v0.1.0 tags for modules with no existing tags

## Target Repo
zarlcorp/core

## Agent Role
backend

## Files to Create
- `.github/workflows/tag.yml`

## Dependencies
Depends on 032 (CI pipeline) — should merge CI first so tags only land on tested code.

## Notes
- Go module tag convention for nested modules: the tag must include the module path prefix. `go get github.com/zarlcorp/core/pkg/zsync@v0.1.0` requires a tag named `pkg/zsync/v0.1.0`.
- The workflow should be idempotent — if the tag already exists (e.g. re-run), skip gracefully
- Start all modules at v0.1.0, not v1.0.0 — we're pre-stable
- The GITHUB_TOKEN has permission to push tags by default, but the workflow must explicitly set `permissions: contents: write`
- Consider using `peter-evans/create-tag` action or raw `git tag` + `git push` — the latter is simpler
