# 064: zburn TUI Catppuccin update

## Objective
Update zburn's TUI to use the new Catppuccin Mocha palette from zstyle. This is the testbed to prove the new palette works in a real terminal application.

## Context
zstyle is being updated to Catppuccin Mocha (spec 058). zburn is the first tool to adopt it. The TUI should look like a beautifully riced terminal app — the same aesthetic as the web pages but native in the terminal.

## Requirements

### Update zstyle dependency
- Update go.mod to use the new zstyle with Catppuccin palette
- Run `go mod tidy` to pull the latest

### Update all TUI styles
- Replace any hardcoded colors with zstyle references
- Use zstyle's tool accent for zburn-specific highlights
- Ensure all lipgloss styles use the new palette
- Rounded borders, consistent spacing, the riced look

### Verify all TUI views
- Password prompt screen
- Main menu
- Generate view
- Browse/list view
- Detail view
- All should look cohesive with the new palette

### Test
- Build and verify the TUI renders correctly
- All existing tests must pass

## Target Repo
zarlcorp/zburn

## Agent Role
backend

## Files to Modify
- go.mod — updated zstyle dependency
- internal/tui/*.go — any files with hardcoded colors or lipgloss styles
- Any other files that reference zstyle colors

## Dependencies
- 058 (zstyle Catppuccin palette)

## Notes
- Read the current zburn TUI code to understand how it uses zstyle
- The goal is a beautiful terminal experience that matches the web aesthetic
- Don't change functionality — only visuals
- If zburn has hardcoded colors that should come from zstyle, fix that
- The agent should build and manually verify the TUI looks right
