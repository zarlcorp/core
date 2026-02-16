# 105: Burn cascade with confirmation

## Objective
Implement the full burn flow for personas — confirmation dialog followed by best-effort cascading deletion of all associated resources (credentials, email forwarding, phone numbers).

## Context
When a user burns a persona, zburn should clean up everything associated with it: local credentials, Namecheap email forwarding rules, and Twilio phone numbers. The cleanup is best-effort — if an API call fails (network down, service error), report the failure but continue burning everything else.

## Requirements

### Confirmation dialog
- When user presses delete/burn on a persona, show confirmation:
  ```
  Burn [FirstName LastName]?

  This will:
  - Delete all credentials (X)
  - Remove email forwarding for persona@domain.com
  - Release phone number +44XXXXXXXXX

  This cannot be undone. (y/n)
  ```
- The summary should reflect what's actually configured — don't mention email forwarding if Namecheap isn't configured, don't mention phone numbers if none are provisioned
- Must press `y` to confirm, any other key cancels

### Burn cascade
On confirmation, execute in order:
1. **Delete credentials** — delete all credentials where `IdentityID` matches the persona
   - Load all credentials via `List()`, filter by `IdentityID`, delete each
2. **Remove email forwarding** — if Namecheap is configured and the persona has a forwarded email
   - Call `RemoveForwarding(ctx, domain, mailbox)` from spec 099
   - If it fails: record error, continue
3. **Release phone number** — if Twilio is configured and the persona has a provisioned number
   - Call `ReleaseNumber(ctx, numberSID)` from spec 101
   - If it fails: record error, continue
4. **Delete identity** — delete the identity from zstore

### Result reporting
- After burn completes, show a summary:
  ```
  Burned [FirstName LastName]
  - Deleted 3 credentials
  - Removed email forwarding for john.doe@domain.com
  - Released phone number +447123456789
  ```
- If any step failed:
  ```
  Burned [FirstName LastName] (with errors)
  - Deleted 3 credentials
  - Email forwarding removal failed: network timeout
  - Released phone number +447123456789
  ```
- Show for 3 seconds or until key press, then return to list view

### Credential deletion confirmation
- Individual credential deletion (from credential detail view, spec 103) also gets a confirmation:
  ```
  Delete credential [Label]? This cannot be undone. (y/n)
  ```

### Testing
- Test burn cascade with all integrations configured (mock API calls)
- Test burn with partial configuration (no Twilio, no Namecheap)
- Test burn with API failures (verify best-effort continues)
- Test confirmation dialog cancellation
- Test result reporting with mixed success/failure

## Target Repo
zarlcorp/zburn

## Agent Role
backend

## Depends On
- 098 (zstore migration — credential deletion)
- 099 (Namecheap — forwarding removal)
- 101 (Twilio — number release)

## Files to Modify
- internal/tui/burn.go (new — burn confirmation + cascade logic)
- internal/tui/detail.go (modify — wire up burn action)
- internal/tui/tui.go (update routing for burn flow)
- internal/burn/burn.go (new — cascade logic separate from TUI, testable)
- internal/burn/burn_test.go (new — test cascade with mocked clients)

## Notes
Depends on 098, 099, 101. Gmail does NOT need cleanup on burn — the forwarding rule removal (Namecheap) stops new emails, and old emails in the shared inbox are harmless. The cascade logic should be in its own `internal/burn` package so it's testable without the TUI.
