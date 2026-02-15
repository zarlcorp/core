# 083: Redesign zarlcorp logo as lowercase

## Objective
The ASCII art logo in zstyle uses uppercase box-drawing letterforms. Redesign it as lowercase to match the zarlcorp voice guide (everything lowercase).

## Context
The current logo (added in 081) spells ZARLCORP using uppercase-style box-drawing characters. The user wants all text lowercase, including the branding logo. The logo is used on password/splash screens of every TUI tool.

## Requirements

### Redesign `Logo` in `pkg/zstyle/logo.go`

Replace the current uppercase logo with a lowercase version using box-drawing characters. The logo must:
- Spell "zarlcorp" in lowercase letterforms
- Be 3 lines tall (same as current)
- Fit comfortably in 80 columns
- Use box-drawing characters (┌ ─ ┐ └ ┘ │ ├ ┤ ┬ ┴ etc.)
- Be readable — each letter should be clearly distinguishable
- Letters should NOT bleed into each other — add 1-space gaps between letters if needed for clarity

Lowercase letterform suggestions (3 lines tall, using box-drawing):
- z: ──┐ / ┌─┘ / └──
- a: ┌─┐ / ├─┤ / ┴ ┴  (same as uppercase A since there's no good lowercase distinction at this size)
- r: ┬─┐ / ├┬┘ / ┴└  (no descender)
- l: │ / │ / └─
- c: ┌─ / │  / └─
- o: ┌─┐ / │ │ / └─┘
- r: (same as above)
- p: ┌─┐ / ├─┘ / │  (descender continues down)

These are suggestions — the agent should pick letterforms that look clean and readable together. Test by printing the logo in a terminal with a monospace font. Prioritize readability over cleverness.

### Update test

Update `pkg/zstyle/logo_test.go` to verify the new logo is non-empty and has 3 lines.

Run `go test ./pkg/zstyle/...` to verify.

## Target Repo
zarlcorp/core

## Agent Role
backend

## Files to Modify
- pkg/zstyle/logo.go
- pkg/zstyle/logo_test.go

## Notes
Keep the same exports: `Logo` string constant and `StyledLogo(s lipgloss.Style) string` function. No API changes. After this merges, tag `pkg/zstyle/v0.4.0` so zburn can pick it up.
