# 078: Automate Homebrew tap update on release

## Objective
When a tool release is tagged and the GitHub Release is created, automatically update the Homebrew formula in `zarlcorp/homebrew-tap` with the new version, URLs, and SHA256 checksums.

## Context
Currently the Homebrew tap is updated manually — downloading binaries, computing checksums, editing the formula, committing, and creating a PR. This should happen automatically as part of the reusable release workflow in `zarlcorp/.github`.

The release workflow (`go-release.yml`) already builds binaries and creates a GitHub Release. We need to add a step (or a new job) that updates the tap after the release is published.

`GITHUB_TOKEN` only has access to the calling repo, so pushing to `homebrew-tap` requires a PAT. This should be passed as a secret from the calling workflow.

## Requirements

### Update the reusable workflow (`go-release.yml`)

Add a new step after "Create GitHub Release" (or a new job that depends on the release job):

1. Compute SHA256 checksums for all 4 binaries in `dist/`
2. Clone `zarlcorp/homebrew-tap` using the TAP_TOKEN
3. Update `Formula/<binary-name>.rb`:
   - Set `version` to the new version (strip the `v` prefix)
   - Update all 4 `url` lines to point to the new release
   - Update all 4 `sha256` lines with the computed checksums
   - Update the `test` block version string
4. Commit and push directly to main (no PR needed for automated formula bumps)

The step should use `sed` or a simple script to update the formula. The formula structure is standardised — all tool formulas follow the same pattern (see `homebrew-tap/Formula/zburn.rb`).

### Add a secret input to the reusable workflow

The reusable workflow needs a `secrets` input for the PAT:

```yaml
on:
  workflow_call:
    inputs:
      binary-name:
        description: "Name of the binary to release"
        type: string
        required: true
    secrets:
      TAP_TOKEN:
        description: "PAT with push access to zarlcorp/homebrew-tap"
        required: true
```

### Update calling workflows

Each tool's `release.yml` needs to pass the secret:

```yaml
jobs:
  release:
    uses: zarlcorp/.github/.github/workflows/go-release.yml@main
    with:
      binary-name: zburn
    secrets:
      TAP_TOKEN: ${{ secrets.TAP_TOKEN }}
```

### Set up the secret

The `TAP_TOKEN` secret needs to be configured. Options (in order of preference):
1. **Org-level secret** in zarlcorp — available to all tool repos automatically
2. **Per-repo secret** in each tool repo — more work to maintain

The PAT needs `repo` scope (or fine-grained: `contents: write` on `homebrew-tap`).

**NOTE:** The agent should update the workflow files. The PM will handle setting up the actual secret via `gh secret set` since that requires the PAT value.

## Target Repo
zarlcorp/dot-github

## Agent Role
backend

## Files to Modify
- .github/workflows/go-release.yml (add tap update step + secret input)

## Notes
The agent should also update `zarlcorp/zburn/.github/workflows/release.yml` to pass the secret. The PM will handle:
1. Creating the PAT
2. Setting the org/repo secret via `gh secret set`

After this is merged, every `git tag v*` + `git push --tags` on any tool repo will automatically: build binaries → create release → update Homebrew tap.
