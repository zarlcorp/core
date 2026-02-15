# 082: zburn TUI polish pass

## Objective
Improve the overall feel of zburn's TUI — fix layout jank, improve spacing, refine wording, and integrate the zarlcorp logo.

## Context
The TUI works functionally but has rough edges: flash messages cause content to shift vertically, spacing is inconsistent, and some wording could be tighter. This is a polish pass, not a feature add.

## Requirements

### 1. Integrate zarlcorp logo on password screen
Replace the plain `zstyle.Title.Render("zburn")` in the password view with the zarlcorp ASCII logo from zstyle (added in 081). Show the tool name "zburn" below the logo in a smaller style.

### 2. Fix layout jank from flash messages
Flash messages ("saved", "copied!", etc.) currently appear and disappear, causing the content below to shift up/down. Fix by reserving space for the flash — always render a line for it, but make it empty/invisible when there's no flash. This prevents vertical shifting.

Apply this fix in all views that use flash messages:
- `generate.go` — "saved", "copied!", "copy: err"
- `list.go` — "deleted", "load: err"
- `detail.go` — "copied!", "copy: err"

### 3. Consistent spacing
Review all views and ensure consistent vertical spacing:
- Title → 1 blank line → content
- Content → 1 blank line → flash area (always present, even if empty)
- Flash area → help text
- No trailing double newlines

### 4. Refine wording
- Menu: show `zburn v0.x.x` with version below the logo, not as a title
- Help text: ensure all views use the same format (`key action  key action  ...`)
- Lowercase everything per voice guide (should already be done from 069, but verify)

### 5. Review colors
Ensure all text uses appropriate zstyle colors:
- Labels: `MutedText`
- Values: default text color
- Active/selected: uses accent color via `ActiveBorder`
- Flash success: `StatusOK`
- Flash error: `StatusErr`
- Help text: `MutedText`

No new colors — just verify existing usage is consistent.

## Target Repo
zarlcorp/zburn

## Agent Role
backend

## Files to Modify
- internal/tui/password.go (logo integration)
- internal/tui/menu.go (version display, spacing)
- internal/tui/generate.go (flash jank fix, spacing)
- internal/tui/list.go (flash jank fix, spacing)
- internal/tui/detail.go (flash jank fix, spacing)

## Notes
Depends on 081 (ASCII logo in zstyle). Can be developed in parallel if the agent stubs the logo import initially, but ideally 081 merges first.

Read every file before editing. This is a polish pass — don't change functionality, just presentation. Run `go test ./...` after every change.
