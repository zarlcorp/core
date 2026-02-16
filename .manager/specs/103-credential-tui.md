# 103: Credential TUI views + clipboard

## Objective
Add credential management views to zburn's TUI — list, add, edit, delete credentials for a persona, with TOTP live display and clipboard copy.

## Context
Spec 098 adds the credential model and zstore collection. This spec surfaces credentials in the TUI. Users should be able to manage multiple credentials per persona (websites, logins, TOTP secrets) from the persona detail view.

## Requirements

### Navigation
- From the identity detail view, add a "Credentials" section below the identity info
  - Shows a list of credentials by label (e.g. "GitHub", "Netflix")
  - Count indicator: "Credentials (3)"
- Navigate into a credential to see full details
- "Add credential" action from the detail view

### Credential detail view
- Display: Label, URL, Username
- Password: masked by default (`••••••••`), press `r` to reveal/hide
- TOTP: if a TOTP secret is set, show the live 6-digit code with a 30-second countdown timer
  - Code refreshes automatically when the window rolls over
  - Use `zcrypto.TOTPCode` for generation
- Notes: displayed if present
- Timestamps: CreatedAt, UpdatedAt

### Actions (keyboard shortcuts)
- `c` — copy password to clipboard via `pbcopy` (`exec.Command("pbcopy")` with password piped to stdin)
- `t` — copy TOTP code to clipboard (only shown if TOTP secret is set)
- `e` — edit credential (opens edit form)
- `d` — delete credential (with confirmation: "Delete credential [Label]? y/n")
- `q` or `esc` — back to identity detail

### Add/Edit credential form
- Fields: Label, URL, Username, Password, TOTP Secret, Notes
- Password field: option to auto-generate (`g` to generate) or type manually
  - Auto-generate uses `zcrypto.GeneratePassword(20)` as default
- TOTP Secret: paste the base32 secret from the service's QR setup
  - Validate it's valid base32 on input
- On save: set CreatedAt (add) or UpdatedAt (edit), generate ID (add), write to zstore

### Confirmation dialogs
- Delete credential: "Delete credential [Label]? This cannot be undone. (y/n)"
- Must press `y` to confirm, any other key cancels

### Clipboard
- `pbcopy` via `os/exec` — macOS only for now
- Brief status message after copy: "Password copied" / "TOTP code copied" (2 second display)

## Target Repo
zarlcorp/zburn

## Agent Role
backend

## Depends On
- 098 (credential model + zstore migration)
- 096 (zcrypto TOTP for live code display)

## Files to Modify
- internal/tui/credential_list.go (new — credential list within detail view)
- internal/tui/credential_detail.go (new — single credential view)
- internal/tui/credential_form.go (new — add/edit form)
- internal/tui/detail.go (modify — add credentials section)
- internal/tui/clipboard.go (new — pbcopy helper)
- internal/tui/tui.go (update routing/model for new views)

## Notes
Depends on 098 and 096. The TOTP countdown timer uses Bubble Tea's tick mechanism (tea.Tick) to refresh every second. Clipboard is macOS-only (`pbcopy`) — cross-platform can come later.
