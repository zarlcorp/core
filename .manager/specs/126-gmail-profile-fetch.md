# 126: Gmail Profile Fetch

## Objective
Fetch the authenticated user's Gmail address after OAuth and store it in GmailSettings so it can be used as the forwarding target for burner domain emails.

## Context
The auto catch-all forwarding feature (spec 127) needs to know which Gmail address to forward to. Currently GmailSettings stores OAuth credentials and tokens but not the email address. The Gmail API provides a `users/me/profile` endpoint that returns the authenticated user's email.

## Requirements
- Add a `GetProfile` method to the gmail client (`internal/gmail/client.go`) that calls `GET https://gmail.googleapis.com/gmail/v1/users/me/profile` and returns the email address
- Add an `Email` field to `GmailSettings` in `internal/tui/config.go`
- After successful OAuth in `settings_gmail.go`, fetch the profile to get the email address and include it in the saved `GmailSettings`
- Update `GmailSettings.Configured()` to also require `Email != ""` (token + email = fully configured)
- Show the connected Gmail address in the Gmail settings view (e.g. "connected as you@gmail.com")
- Write tests for the `GetProfile` method using httptest
- All existing tests must continue to pass

## Target Repo
zarlcorp/zburn

## Agent Role
backend

## Files to Modify
- internal/gmail/client.go — add GetProfile method
- internal/gmail/client_test.go — add tests for GetProfile
- internal/tui/config.go — add Email field to GmailSettings
- internal/tui/settings_gmail.go — fetch profile after OAuth, show email in view

## Notes
- The profile endpoint returns `{ "emailAddress": "...", "messagesTotal": N, "threadsTotal": N, "historyId": "..." }` — we only need emailAddress
- The access token from the OAuth flow is used for the profile fetch
- If the profile fetch fails, treat it as an OAuth failure (don't save partial config)
