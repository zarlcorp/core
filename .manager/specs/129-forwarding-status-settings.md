# 129: Forwarding Status in Settings

## Objective
Show per-domain email forwarding status in the settings view so the user can see which domains have catch-all forwarding configured and which don't.

## Context
Specs 126-127 set up automatic catch-all forwarding. This spec adds visibility into the forwarding state. Users need to see:
- Which domains have forwarding active
- What Gmail address they forward to
- If auth is missing for either service

Depends on spec 127 (auto catch-all forwarding).

## Requirements

### Forwarding status view
Add a forwarding status section to the settings view (or as a sub-view accessible from settings). Show each burner domain with its forwarding state:

```
  forwarding

  moxmail.site        * → you@gmail.com
  jotmail.xyz         * → you@gmail.com
  snapmail.icu        * → you@gmail.com
  calvera.run         * → you@gmail.com
  ...
  zarlcorp.com        excluded
  zarl.dev            excluded
```

Use zstyle for consistent formatting:
- Domain name left-aligned
- Forwarding target in muted text
- "excluded" for org domains in muted text
- If forwarding is not set up on a domain, show a warning style indicator

### Auth warnings
- If Namecheap is configured but Gmail is not: show "gmail not connected — forwarding inactive"
- If Gmail is configured but Namecheap is not: show "namecheap not connected — no domains"
- If neither: show "configure namecheap and gmail to enable forwarding"

### Data source
- Read forwarding status by calling `GetForwarding` on each domain via the Namecheap client
- Cache the results to avoid API calls on every render — fetch once when entering settings, refresh on demand
- The Gmail address comes from `GmailSettings.Email` (added in spec 126)

### Navigation
- Accessible from the main settings menu (add a "forwarding" menu item)
- Esc returns to settings menu

## Target Repo
zarlcorp/zburn

## Agent Role
backend

## Files to Modify
- internal/tui/settings_forwarding.go — new file for forwarding status view
- internal/tui/settings.go — add "forwarding" menu item
- internal/tui/tui.go — add viewForwarding to view enum, wire up navigation

## Notes
- Keep the view read-only for now — no editing from this screen
- The forwarding status fetch should be async (tea.Cmd) to avoid blocking the UI
- Consider batching GetForwarding calls or running them in parallel for faster load
- Excluded domains list should match the one defined in spec 127 (share the constant)
