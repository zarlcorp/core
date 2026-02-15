# 048: Interactive TUI

## Objective
Replace the placeholder TUI with a full interactive Bubble Tea experience for generating, browsing, and managing disposable identities.

## Context
zburn's TUI currently shows a title and "press q to quit." This spec builds the real interface with multiple views, keyboard navigation, and clipboard support.

The Manifesto says: "Interactive — beautiful TUI with keyboard navigation, menus, real-time feedback. A user who knows one zarlcorp tool knows them all. Same keybindings. Same visual language."

## Requirements

### Views

**1. Password prompt (first view)**
- On launch, prompt for master password using a text input (masked)
- On first use (no store exists), prompt twice to confirm
- On wrong password, show error and re-prompt
- Use `zstyle.Title` for the prompt header

**2. Main menu**
- After password, show menu with options:
  - Generate identity
  - Generate email (quick)
  - Browse saved identities
  - Quit
- Use `zstyle.KeyUp`/`zstyle.KeyDown` for navigation, `zstyle.KeyEnter` to select
- Use `zstyle.KeyQuit` to quit

**3. Generate view**
- Generate a new identity and display all fields
- Actions:
  - `s` — save to store
  - `c` — copy all fields to clipboard
  - `enter` on a field — copy that field to clipboard
  - `n` — generate another (discard current)
  - `esc`/`backspace` — back to menu
- Show confirmation flash message on save/copy

**4. List view**
- Show all saved identities in a scrollable list
- Columns: ID, name, email, created date
- `enter` — view identity details
- `d` — delete selected identity (with confirmation)
- `esc`/`backspace` — back to menu
- Show "no saved identities" if empty
- Use `zstyle.KeyUp`/`zstyle.KeyDown` for scrolling

**5. Detail view**
- Show all fields of a saved identity
- `enter` on a field — copy to clipboard
- `c` — copy all fields
- `d` — delete this identity
- `esc`/`backspace` — back to list

### Clipboard
- Use `golang.design/x/clipboard` or shell out to `pbcopy`/`xclip`/`xsel`
- Prefer a pure Go approach if available, fall back to exec
- Show a brief flash message "copied!" that fades after 1 second

### Styling
- All styling via `zstyle` — Title, Subtitle, MutedText, Highlight, StatusOK, StatusErr, Border
- Use `zstyle.ActiveBorder` for the focused/selected item
- Keep the layout clean: left-aligned, consistent padding, no box-drawing unless zstyle provides it

### Architecture
- Each view is a separate Bubble Tea model in its own file under `internal/tui/`
- Root model delegates to the active view
- Views communicate via messages (tea.Msg), not direct calls

### Testing
- Test each view's Update function: key handling, state transitions
- Test that quit keybinding works from every view
- Test navigation: menu → generate → menu, menu → list → detail → list → menu
- Test password prompt handles correct/incorrect input

## Target Repo
zarlcorp/zburn

## Agent Role
backend

## Files to Modify
- internal/tui/tui.go — replace placeholder with root model that delegates to views

## Files to Create
- internal/tui/password.go — password prompt view
- internal/tui/menu.go — main menu view
- internal/tui/generate.go — generate identity view
- internal/tui/list.go — identity list view
- internal/tui/detail.go — identity detail view
- internal/tui/clipboard.go — clipboard helper
- internal/tui/tui_test.go — update existing tests for new structure

## Dependencies
- `internal/identity` (spec 045) — Generator, Identity type
- `internal/store` (spec 046) — Store
- `github.com/charmbracelet/bubbles` — text input (for password), list, table
- `github.com/charmbracelet/bubbletea` — framework
- `github.com/zarlcorp/core/pkg/zstyle` — styling and keybindings

## Notes
- The TUI receives the store and generator as dependencies from main.go — don't construct them inside the TUI package.
- For clipboard, start simple: `exec.Command("pbcopy")` on macOS, `exec.Command("xclip")` on Linux. Wrap in a helper that returns an error if neither is available. Don't add a dependency for this.
- The password prompt is the first thing users see — make it clean. Just the title, a masked input, and a status line.
- Flash messages (e.g., "copied!", "saved!") should auto-dismiss after ~1 second using `tea.Tick`.
- The `bubbles` package provides `textinput` (for password) and `list` (for identity list) — use those rather than building from scratch.
