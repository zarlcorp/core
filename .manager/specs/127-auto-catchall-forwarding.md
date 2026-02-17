# 127: Auto Catch-All Forwarding

## Objective
Automatically set up wildcard (`*`) email forwarding on all burner domains so that any email sent to any address at those domains is forwarded to the authenticated Gmail account.

## Context
Users buy burner domains via Namecheap and authenticate Gmail for reading forwarded emails. Currently, email forwarding must be set up manually. This spec automates it: when both Namecheap and Gmail are authenticated, set `* → gmail_address` on every burner domain.

Depends on spec 126 (Gmail profile fetch) which stores the Gmail email address in GmailSettings.

## Requirements

### Trigger points
Set up catch-all forwarding when:
1. **Namecheap auth completes** and Gmail is already configured — forward on all CachedDomains
2. **Gmail auth completes** and Namecheap is already configured — forward on all CachedDomains
3. **New domains appear** (Namecheap re-auth refreshes CachedDomains) — forward on any domains that don't have it yet

### Domain exclusion
Skip these domains (org domains, not burners):
- `zarlcorp.com`
- `zarl.dev`

Store the exclusion list as a constant slice in the TUI or config, not hardcoded inline.

### Forwarding rule
For each burner domain, call `namecheap.SetForwarding(ctx, domain, []ForwardingRule{{Mailbox: "*", ForwardTo: gmailAddress}})`.

Use `SetForwarding` (not `AddForwarding`) to be idempotent — calling it multiple times just overwrites with the same catch-all rule.

### Error handling
- If forwarding setup fails on a domain, log the error but continue with remaining domains
- Show failures in the TUI flash/status area
- Do not block the auth flow on forwarding failures

### Implementation
- Add a function like `setupCatchAllForwarding(ctx, ncClient, domains, gmailAddress, excludedDomains)` that iterates domains and sets the wildcard rule
- Call this from the Namecheap save handler (after CachedDomains are refreshed) and the Gmail save handler (after profile is fetched)
- Run the forwarding setup as a tea.Cmd (background, non-blocking)
- Report results via a message (e.g. `forwardingResultMsg{successes int, failures []error}`)

## Target Repo
zarlcorp/zburn

## Agent Role
backend

## Files to Modify
- internal/tui/tui.go — add forwarding setup logic, call from both auth handlers
- internal/tui/settings_namecheap.go — trigger forwarding after save (if Gmail ready)
- internal/tui/settings_gmail.go — trigger forwarding after save (if Namecheap ready)
- internal/namecheap/client.go — no changes needed (SetForwarding already exists)

## Notes
- `SetForwarding` is destructive (replaces all rules) but since we're setting a single catch-all, this is fine — there should be no other rules on burner domains
- The wildcard `*` as a mailbox is supported by Namecheap's email forwarding
- Run forwarding setup asynchronously so it doesn't block the UI
