# 110: Pre-release smoke test for zburn v0.5.0

## Objective
Build the zburn binary and verify all v0.5.0 features work before tagging the release. Exercise every new user flow added in specs 098-105.

## Context
v0.5.0 adds: encrypted zstore backend, credential vault, settings TUI, burn cascade, Namecheap/Gmail/Twilio clients, 2FA code extraction. Unit and integration tests pass, but nobody has run the actual binary through the new flows. External integrations (Namecheap, Gmail, Twilio) cannot be end-to-end tested without real API accounts — focus on local features.

## Requirements

### Automated CLI tests (shell script)
Create `test/smoke.sh` that builds the binary and exercises all non-interactive and scriptable CLI commands:

1. `zburn version` — verify output matches expected version string
2. `zburn email` — verify output is a valid email matching `*@zburn.id` pattern
3. `zburn identity --json` — verify valid JSON with all expected fields (id, firstName, lastName, email, phone, street, city, state, zip, dob, password)
4. `zburn identity --save` — pipe password via stdin, verify exit 0, verify store files created in temp dir
5. `zburn list --json` — pipe password, verify JSON array with the saved identity
6. `zburn forget <id>` — pipe password + confirmation, verify identity removed
7. `zburn list --json` — pipe password, verify empty array

The script should:
- Use a temp directory for the data store (`XDG_DATA_HOME` or pass data dir)
- Clean up after itself
- Exit non-zero on any failure
- Print PASS/FAIL for each test

### Manual TUI checklist
Create `test/SMOKE_TEST.md` documenting each manual verification step. The user runs the binary and follows the checklist:

#### Store lifecycle
- [ ] Fresh start: password prompt appears, create new password
- [ ] Store created: menu appears after password
- [ ] Quit and relaunch: password prompt, enter correct password → menu
- [ ] Wrong password: enter wrong password → error message → retry prompt

#### Generate identity
- [ ] Menu → "generate identity" → shows 10 fields with random data
- [ ] Arrow keys navigate between fields
- [ ] `n` generates a new identity (all fields change)
- [ ] Enter copies selected field to clipboard
- [ ] `c` copies all fields
- [ ] `s` saves identity → flash message confirms save
- [ ] Esc returns to menu

#### Quick email
- [ ] Menu → "generate email" → flash message with email → back to menu
- [ ] Email is in clipboard

#### Browse identities
- [ ] Menu → "browse saved identities" → list shows saved identity
- [ ] Enter on identity → detail view with all 10 fields + credential count (0)
- [ ] Esc from detail → back to list
- [ ] Esc from list → back to menu

#### Credential vault
- [ ] Detail view → `v` → credential list (empty)
- [ ] `a` → credential form with fields: label, URL, username, password, TOTP secret, notes
- [ ] Tab/Shift+Tab navigates fields
- [ ] Ctrl+G on password field generates a password
- [ ] Fill in label + URL + username, save → back to credential list with 1 credential
- [ ] Select credential → Enter → detail view with all fields
- [ ] Edit credential (change label) → verify change persisted
- [ ] Add second credential → verify list shows 2
- [ ] Delete one credential → verify list shows 1 with correct remaining credential

#### TOTP
- [ ] Add credential with TOTP secret: `JBSWY3DPEHPK3PXP` (standard test vector)
- [ ] Credential detail shows rotating 6-digit TOTP code
- [ ] Code changes approximately every 30 seconds

#### Settings
- [ ] Menu → "settings" → three options, all showing "not configured"
- [ ] Namecheap: enter dummy values in all fields → Ctrl+S → back to settings → shows "configured"
- [ ] Gmail: enter dummy client ID + secret → verify form accepts input → Esc (don't attempt OAuth)
- [ ] Twilio: enter dummy SID + token → Ctrl+S → back to settings → shows "configured"
- [ ] Quit and relaunch → settings → Namecheap and Twilio still show "configured"

#### Burn cascade
- [ ] Generate and save a new identity
- [ ] Add 2 credentials to that identity
- [ ] Detail view → `d` → confirmation dialog shows plan: "delete all credentials (2)"
- [ ] Press `n` (or any non-y key) → cancels, back to detail
- [ ] Press `d` again → confirmation → press `y` → burn executes
- [ ] Result screen shows success: credentials deleted, identity deleted
- [ ] Any key → back to list → identity is gone
- [ ] Browse list → original identity still exists (only the burned one removed)

#### CLI commands (manual verification of interactive ones)
- [ ] `zburn identity --save` → prompts for password → saves → exit
- [ ] `zburn list --json` → prompts for password → shows JSON array
- [ ] `zburn forget <id>` → prompts for password → deletes identity

### Acceptance criteria
- `test/smoke.sh` runs green on the built binary
- User completes the manual checklist with no failures
- Any bugs found are filed as issues before release

## Target Repo
zarlcorp/zburn

## Agent Role
testing

## Files to Create
- test/smoke.sh (automated CLI smoke test script)
- test/SMOKE_TEST.md (manual TUI verification checklist)

## Depends On
- 107 (integration tests merged — done)

## Notes
The automated script exercises what it can. The manual checklist covers TUI-specific interactions that need a real terminal. External service integrations (Namecheap API calls, Gmail OAuth, Twilio provisioning) are out of scope — they're unit-tested but can't be end-to-end verified without real accounts. The data directory can be controlled via `XDG_DATA_HOME` environment variable pointing to a temp dir.
