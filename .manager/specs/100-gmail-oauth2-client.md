# 100: Gmail OAuth2 client

## Objective
Create an internal Gmail client in zburn that handles OAuth2 authentication and reads emails from a shared Gmail inbox.

## Context
zburn personas have email addresses that forward to a shared Gmail account. We need to read incoming emails to extract 2FA codes and other verification messages. This requires OAuth2 access to the Gmail API.

The OAuth2 flow for a TUI app: start a temporary localhost HTTP server, open the user's browser to Google's consent screen, capture the callback with the auth code, exchange for tokens, store them. Subsequent calls use the refresh token silently.

## Requirements

### OAuth2 flow
- Create `internal/gmail/auth.go`:
  ```go
  type OAuthConfig struct {
      ClientID     string
      ClientSecret string
  }

  type Token struct {
      AccessToken  string    `json:"access_token"`
      RefreshToken string    `json:"refresh_token"`
      Expiry       time.Time `json:"expiry"`
  }
  ```
- `Authenticate(ctx context.Context, cfg OAuthConfig) (*Token, error)`
  - Starts a temporary HTTP server on a random localhost port
  - Builds the Google OAuth2 authorization URL with:
    - Scope: `https://www.googleapis.com/auth/gmail.readonly`
    - Redirect URI: `http://localhost:{port}/callback`
    - Access type: `offline` (to get refresh token)
    - Prompt: `consent` (force consent screen to ensure refresh token)
  - Opens the URL in the default browser (use `exec.Command("open", url)` on macOS)
  - Waits for the callback with the auth code
  - Exchanges the code for tokens via Google's token endpoint
  - Returns the Token (caller stores it in zstore)
  - Shuts down the HTTP server

### Token refresh
- `RefreshToken(ctx context.Context, cfg OAuthConfig, refreshToken string) (*Token, error)`
  - Exchanges refresh token for a new access token
  - Returns updated Token

### Gmail API client
- Create `internal/gmail/client.go`:
  ```go
  type Client struct { ... }

  func NewClient(accessToken string) *Client
  ```
- `ListMessages(ctx context.Context, query string, maxResults int) ([]Message, error)`
  - Calls `GET https://gmail.googleapis.com/gmail/v1/users/me/messages`
  - `query` is a Gmail search query (e.g. `to:john.doe@userdomain.com`)
  - Returns message IDs and snippet previews
- `GetMessage(ctx context.Context, messageID string) (*Message, error)`
  - Calls `GET https://gmail.googleapis.com/gmail/v1/users/me/messages/{id}`
  - Returns full message with headers (From, To, Subject, Date) and body text
  - Handles both plain text and MIME multipart messages — extract the text/plain part

### Message type
```go
type Message struct {
    ID      string
    From    string
    To      string
    Subject string
    Date    time.Time
    Body    string // plain text content
}
```

### Testing
- Test OAuth2 URL construction (correct scopes, redirect URI)
- Test token exchange with an httptest.Server returning canned JSON
- Test token refresh with an httptest.Server
- Test message listing and parsing with canned Gmail API JSON responses
- Test MIME multipart body extraction
- Do NOT call real Google APIs in tests

### Dependencies
- Standard library only for HTTP and JSON
- `golang.org/x/oauth2` is NOT required — the OAuth2 flow is simple enough to implement with `net/http` and direct token exchange. Keeps the dependency count down.

## Target Repo
zarlcorp/zburn

## Agent Role
backend

## Files to Modify
- internal/gmail/auth.go
- internal/gmail/client.go
- internal/gmail/auth_test.go
- internal/gmail/client_test.go

## Notes
Independent — no dependency on other specs. The client takes tokens as parameters; token storage in zstore and the initial OAuth2 setup flow are handled by the settings TUI (spec 104). The gmail.readonly scope is sufficient — we only read, never send.
