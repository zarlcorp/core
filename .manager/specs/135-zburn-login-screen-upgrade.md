# 135: Upgrade zburn login screen to match zvault

## Objective
Replace zburn's minimal single-field password screen with zvault's richer two-field design: title, description, labeled fields, tab switching, and bullet echo character.

## Context
zvault's password screen has a polished UX with a title ("create new vault" / "unlock vault"), description text, labeled "password" and "confirm" fields visible simultaneously on first run, and tab to switch between them. zburn's screen is a bare single-field sequential prompt. This spec brings zburn's login in line with zvault's quality.

## Requirements

### 1. Two-field password model
Replace the single `input` field with two fields:
- `password textinput.Model` — main password field
- `confirm textinput.Model` — confirmation field (first run only)

Add a `focused` field (enum: `fieldPassword`, `fieldConfirm`) to track which input has focus.

Remove the `confirming bool` and `firstPass string` fields — the two-field approach makes these unnecessary.

### 2. Field configuration
Both fields should use:
- `EchoMode: textinput.EchoPassword`
- `EchoCharacter: '•'` (bullet, not asterisk)
- `PromptStyle` with `zstyle.ZburnAccent` foreground
- `TextStyle` with `zstyle.Text` foreground

### 3. Submit logic
- If `firstRun` and focused on password field: tab/enter moves to confirm field
- If `firstRun` and focused on confirm field: compare both values
  - Mismatch: show error "passwords do not match", clear confirm field
  - Match: submit the password
- If not `firstRun`: submit immediately from password field
- Empty password: show error "password cannot be empty"

### 4. Tab switching
- Tab key switches focus between password and confirm fields (first run only)
- `password.Blur()` / `confirm.Focus()` and vice versa

### 5. View rendering
Match zvault's layout:
```
\n
  <logo>
  zburn

  <title>                    // "create new store" or "unlock store"

  <description>              // "choose a master password..." or "enter your master password."

  password
  <password input>

  confirm                    // first run only
  <confirm input>

  <error if any>
```

All text lowercase:
- First run title: "create new store"
- Unlock title: "unlock store"
- First run description: "choose a master password to protect your store."
- Unlock description: "enter your master password."
- Field labels: "password", "confirm"

Title styled with `ZburnAccent` + bold. Description with `MutedText`. Labels with `Subtext1`.

### 6. Error handling
- Clear error on any key press (like zvault does)
- Display error with `zstyle.StatusErr`

### 7. Keep existing message types
Keep `passwordSubmitMsg` and `passwordErrMsg` as they are — the parent model in tui.go consumes these and they work fine. Only the internal password model changes.

### 8. Update `newPasswordModel`
The function currently takes `firstRun bool`. Keep this signature — it's called from tui.go which determines first run from the store directory existence check.

### 9. Update tests
Update password-related tests in `tui_test.go` to verify:
- Two-field layout on first run (check for "confirm" in View output)
- Single-field layout on unlock (no "confirm")
- Tab switching between fields
- Password mismatch error
- Successful submission

## Target Repo
zarlcorp/zburn

## Agent Role
backend

## Files to Modify
- internal/tui/password.go (rewrite)
- internal/tui/tui_test.go (update password tests if any)

## Notes
- Reference zvault's `internal/tui/password.go` for the pattern — but use `zstyle.ZburnAccent` not `ZvaultAccent`
- The parent model in `tui.go` should need zero changes — it already handles `passwordSubmitMsg` and `passwordErrMsg`
- Keep Ctrl+C → quit behavior from the current zburn implementation
