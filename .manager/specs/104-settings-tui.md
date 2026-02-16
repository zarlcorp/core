# 104: Settings TUI

## Objective
Add a settings view to zburn's TUI for configuring external service integrations (Namecheap, Gmail, Twilio). Each integration is optional — features are hidden when unconfigured.

## Context
Specs 099-101 create API clients for Namecheap, Gmail, and Twilio. This spec provides the TUI for users to enter their API credentials and trigger OAuth flows. Credentials are stored encrypted in zstore.

## Requirements

### Settings menu
- Add "Settings" option to the main menu (alongside Generate, List, etc.)
- Settings view shows three integration sections:
  - Namecheap (email forwarding)
  - Gmail (email reading)
  - Twilio (phone numbers)
- Each section shows status: "Configured" / "Not configured"
- Navigate into each to configure

### Namecheap settings
- Fields: API User, API Key, Username, Client IP
- Stored in zstore as config collection: `zstore.Collection[namecheap.Config](store, "config")` with key `"namecheap"`
- After saving, test the connection by calling `GetForwarding` on a domain (optional validation)
- Also configure: list of domains available for email forwarding
  - User enters domains comma-separated or one per line
  - Stored alongside the Namecheap config

### Gmail settings
- Show current status: "Connected as shared@gmail.com" / "Not connected"
- Fields: Client ID, Client Secret (from Google Cloud Console)
- "Connect Gmail" button triggers the OAuth2 flow (spec 100's `Authenticate`)
  - Opens browser, waits for callback, stores tokens
- Tokens stored in zstore config collection with key `"gmail"`
- "Disconnect" option removes stored tokens

### Twilio settings
- Fields: Account SID, Auth Token
- Stored in zstore config collection with key `"twilio"`
- Preferred countries: checkboxes/options for UK, US, or both
- After saving, test the connection by listing available numbers (optional validation)

### Config collection
- All settings stored in: `zstore.Collection[json.RawMessage](store, "config")` or use typed structs per provider
- Config keys: `"namecheap"`, `"gmail"`, `"twilio"`

### Feature gating
- When Namecheap is not configured: email forwarding options don't appear in persona generation
- When Gmail is not configured: email reading features don't appear
- When Twilio is not configured: phone number provisioning doesn't appear
- The app works fine with zero integrations configured — just generates local personas like today

### Testing
- Test that settings round-trip through zstore (save + load)
- Test feature gating logic (configured vs unconfigured states)

## Target Repo
zarlcorp/zburn

## Agent Role
backend

## Depends On
- 098 (zstore migration — for config collection)
- 099 (Namecheap client — for connection testing)
- 100 (Gmail client — for OAuth2 flow)
- 101 (Twilio client — for connection testing)

## Files to Modify
- internal/tui/settings.go (new — settings menu)
- internal/tui/settings_namecheap.go (new — Namecheap config form)
- internal/tui/settings_gmail.go (new — Gmail OAuth2 setup)
- internal/tui/settings_twilio.go (new — Twilio config form)
- internal/tui/menu.go (modify — add Settings option)
- internal/tui/tui.go (update routing for settings views)

## Notes
Depends on 098, 099, 100, 101. This is the integration point where API clients meet encrypted storage. The OAuth2 flow for Gmail is the most complex part — it involves starting a temp HTTP server and opening the browser from within the TUI (the TUI should show a "Waiting for browser authorization..." state).
