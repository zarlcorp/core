# 050: Update zburn README

## Objective
Update the zburn README to reflect the current state of the tool — all CLI commands are implemented, TUI is fully built out.

## Context
The README was written during scaffold (spec 042) when CLI commands were planned but not built. Specs 045-048 have since been merged, delivering the full identity generator, encrypted store, CLI subcommands, and interactive TUI.

## Requirements
- Remove "Planned commands (not yet implemented)" — all commands work
- Document all CLI commands with examples and flags:
  - `zburn` — launch interactive TUI
  - `zburn version` — print version
  - `zburn email` — generate a burner email
  - `zburn identity` — generate a complete identity (supports `--json`, `--save`)
  - `zburn list` — list saved identities (supports `--json`)
  - `zburn forget <id>` — delete a saved identity
- Add a brief TUI section describing the interactive experience (password prompt, menu, generate, browse, clipboard)
- Keep it concise — this is a README, not docs

## Target Repo
zarlcorp/zburn

## Agent Role
backend

## Files to Modify
- README.md

## Notes
- Read the current main.go and cli.go to confirm exact command signatures
- Don't add screenshots — just describe the TUI briefly
