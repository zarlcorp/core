# 031: Systematic review of all shared packages

## Objective
Review all 7 packages in zarlcorp/core for coding standard compliance, consistency, and quality. Four packages were migrated from a monorepo (zsync, zcache, zoptions, zfilesystem), three were built new by sub-agents (zapp, zstyle, zcrypto). They need to look and feel like they were written by the same team.

## Context
Each package was built or migrated independently. This review ensures they all follow the zarlcorp coding standards documented in `.claude/CLAUDE.md` and `.claude/CLAUDE_GO.md`. The review should produce a concrete list of fixes, not just observations.

## Requirements

### Review checklist — apply to every package

**1. Error handling**
- No "failed to", "unable to", "could not", "error" prefixes
- Direct context wrapping: `"open file: %w"`, `"create cipher: %w"`
- Errors wrapped at every boundary
- No logging of errors (library code — caller decides)

**2. Package documentation**
- Package comment on the first file (alphabetically or main file)
- Brief description of what the package does
- Usage example in the package doc
- Exported types/functions have doc comments

**3. Naming**
- Scope-based naming (short names for small scopes)
- No unnecessary abbreviations
- Consistent naming patterns across packages

**4. Code quality**
- Early returns over if/else chains
- No duplicated code in branches
- No over-engineering (unused abstractions, unnecessary interfaces)
- No dead code or unused exports

**5. Test quality**
- Table-driven tests where appropriate
- External test package (`_test` suffix) — tests the public API
- Adequate coverage of edge cases
- No test helper functions that swallow errors

**6. Module hygiene**
- `go.mod` clean — no unnecessary dependencies
- `go vet ./...` passes
- `go build ./...` passes
- No deprecated API usage

**7. Cross-package consistency**
- Similar patterns used for similar things (e.g. options, constructors)
- Consistent file organization
- Consistent comment style (lowercase, terse)

### Deliverable
Create a file `pkg/REVIEW.md` with:

1. **Per-package findings** — issues found in each package, with file:line references
2. **Cross-package issues** — inconsistencies between packages
3. **Recommended fixes** — concrete changes, ordered by priority (critical → nice-to-have)

Then **apply all fixes** directly. The review is not advisory — fix everything you find.

After fixing, run `go test ./...` from each package directory to verify nothing broke.

### Packages to review (in dependency order)
1. `pkg/zoptions` — foundation (no deps)
2. `pkg/zsync` — foundation (no deps)
3. `pkg/zcache` — data layer
4. `pkg/zfilesystem` — data layer
5. `pkg/zstyle` — presentation (lipgloss, bubbles)
6. `pkg/zcrypto` — security (x/crypto)
7. `pkg/zapp` — top layer (depends on zoptions)

## Acceptance Criteria
1. Every package passes `go vet ./...`
2. Every package passes `go test -race ./...`
3. No "failed to" / "unable to" / "could not" error prefixes remain
4. Every exported function has a doc comment
5. Package docs exist with usage examples
6. `pkg/REVIEW.md` documents findings and fixes applied
7. Consistent patterns across all packages

## Target Repo
zarlcorp/core

## Agent Role
backend

## Files to Review
- `pkg/zoptions/*.go`
- `pkg/zsync/*.go`
- `pkg/zcache/*.go`
- `pkg/zfilesystem/*.go`
- `pkg/zstyle/*.go`
- `pkg/zcrypto/*.go`
- `pkg/zapp/*.go`

## Files to Create
- `pkg/REVIEW.md` — findings and fixes applied

## Notes
- The migrated packages (zsync, zcache, zoptions, zfilesystem) may have legacy patterns from before the zarlcorp standards were defined. These need the most attention.
- The new packages (zapp, zstyle, zcrypto) were built to spec but by different agents — check for subtle inconsistencies.
- Read `.claude/CLAUDE.md` and `.claude/CLAUDE_GO.md` first to understand the standards.
- This is library code — it must be fast and bulletproof. Prevent panics at all costs.
- Do NOT add features or change public APIs. This is a quality pass, not a redesign.
- Do NOT add comments or docs to code you didn't change for other reasons. Only add docs where they're missing on exported symbols.
