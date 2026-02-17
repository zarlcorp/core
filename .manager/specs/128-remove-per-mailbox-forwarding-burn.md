# 128: Remove Per-Mailbox Forwarding from Burn

## Objective
Remove the per-mailbox email forwarding teardown from the burn cascade. With catch-all wildcard forwarding (spec 127), individual forwarding rules are no longer created per identity, so there's nothing to tear down on delete.

## Context
The burn cascade (`internal/burn/burn.go`) currently:
1. Deletes credentials
2. Removes email forwarding for the identity's mailbox
3. Releases phone number
4. Deletes the identity

Step 2 is no longer needed because forwarding is now a domain-level catch-all (`*@domain → gmail`), not per-identity. Removing it simplifies the burn flow and eliminates an API call during deletion.

This spec is independent of specs 126/127 — it can be done in parallel.

## Requirements

### burn package (`internal/burn/burn.go`)
- Remove the `EmailForwarder` interface
- Remove the `EmailConfig` struct
- Remove the `Email` and `Forwarder` fields from `burn.Request`
- Remove the `removeEmail` method and its call from `Execute`
- Remove the email forwarding step from `Plan`
- Update tests in `burn_test.go` to remove email forwarding test cases

### TUI (`internal/tui/tui.go`)
- Remove the `EmailForwarder` field from `ExternalServices` (or the entire Forwarder reference)
- Remove the `EmailDomain` field from `ExternalServices`
- Stop populating `req.Email` and `req.Forwarder` in the burn confirmation handler
- Remove the `splitEmail` helper if it's only used for burn forwarding

### Tests
- Update `internal/burn/burn_test.go` — remove email forwarding test scenarios
- Update `internal/tui/tui_test.go` — remove email forwarding references from burn tests (e.g. `TestBurnConfirmViewShowsPlan` references "remove email forwarding")
- Update `internal/tui/integration_test.go` — remove forwarding references if present
- All remaining tests must pass

## Target Repo
zarlcorp/zburn

## Agent Role
backend

## Files to Modify
- internal/burn/burn.go — remove EmailForwarder, EmailConfig, removeEmail
- internal/burn/burn_test.go — remove forwarding test cases
- internal/tui/tui.go — remove Forwarder/EmailDomain from ExternalServices, stop populating burn request with email config
- internal/tui/tui_test.go — update burn-related tests
- internal/tui/integration_test.go — update if forwarding references exist

## Notes
- The `ExternalServices.Forwarder` field may still be needed by spec 127 for setting up catch-all forwarding. Check if Forwarder is used elsewhere before removing entirely. If it's only used in the burn flow, remove it. If 127 needs it, leave it but remove burn-specific usage.
- Actually, spec 127 creates a new forwarding setup flow that's independent of ExternalServices.Forwarder — it uses the Namecheap client directly from settings. So Forwarder can be fully removed from ExternalServices.
