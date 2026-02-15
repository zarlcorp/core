# 030: Minimal root README

## Objective
Add a concise README.md to the repo root. Public-facing but minimal — what zarlcorp is, what packages are available, how to use them. Not a copy of the manifesto.

## Context
The repo currently has no README. Anyone landing on `github.com/zarlcorp/core` sees a file listing and nothing else. The manifesto is thorough but internal. The README is the public face.

Godoc comments on each package are excellent and sufficient — no per-package READMEs needed. The root README just needs to point people in the right direction.

## Requirements

### Content
The README should include:

1. **Header** — `zarlcorp/core` with the tagline "tools that fight back"
2. **One-liner** — what this repo is (shared Go packages for zarlcorp privacy tools)
3. **Package table** — name, one-line description, status (stable/wip/planned) for all 7 packages
4. **Install example** — `go get github.com/zarlcorp/core/pkg/<name>`
5. **Link to MANIFESTO.md** — for people who want the full story
6. **License** — MIT, one line

### Style
- No badges, no shields.io
- No emojis
- Terse, lowercase tone matching the manifesto
- Under 60 lines total

### What NOT to include
- Per-package API docs (godoc handles this)
- Contributing guide (not needed yet)
- Build instructions (it's a library)
- CI status badges

## Acceptance Criteria
1. README.md exists at repo root
2. Lists all 7 packages with current status
3. Shows install command
4. Links to MANIFESTO.md
5. Under 60 lines
6. No badges, no emojis

## Target Repo
zarlcorp/core

## Agent Role
backend

## Files to Create
- `README.md`

## Dependencies
None — can run in parallel with anything.

## Notes
- Package statuses: zsync (stable), zcache (stable), zoptions (stable), zfilesystem (stable), zapp (stable), zcrypto (wip), zstyle (wip)
- Update statuses as packages ship — this is a living doc
- Module path pattern: `github.com/zarlcorp/core/pkg/<name>`
