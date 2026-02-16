# 099: Namecheap email forwarding client

## Objective
Create an internal Namecheap API client in zburn for managing email forwarding rules programmatically.

## Context
zburn personas need real email addresses. Users have domains on Namecheap with API access and whitelisted IPs. The Namecheap API supports `domains.dns.setEmailForwarding` and `domains.dns.getEmailForwarding` for managing forwarding rules. Emails forward to a shared Gmail address.

API docs:
- https://www.namecheap.com/support/api/methods/domains-dns/set-email-forwarding/
- https://www.namecheap.com/support/api/methods/domains-dns/get-email-forwarding/

## Requirements

### Client
- Create `internal/namecheap/client.go`:
  ```go
  type Config struct {
      APIUser    string
      APIKey     string
      Username   string
      ClientIP   string
  }

  type Client struct { ... }

  func NewClient(cfg Config) *Client
  ```
- All API calls go to `https://api.namecheap.com/xml.response`
- Namecheap API uses GET requests with query parameters and returns XML

### Methods
- `SetForwarding(ctx context.Context, domain string, rules []ForwardingRule) error`
  - Calls `namecheap.domains.dns.setEmailForwarding`
  - A rule maps a mailbox (e.g. `john.doe`) to a forwarding address (e.g. `shared@gmail.com`)
  - NOTE: this API replaces ALL forwarding rules for the domain on each call — it's not additive. Must GET existing rules first, merge, then SET.
- `GetForwarding(ctx context.Context, domain string) ([]ForwardingRule, error)`
  - Calls `namecheap.domains.dns.getEmailForwarding`
  - Returns current forwarding rules for the domain
- `AddForwarding(ctx context.Context, domain, mailbox, forwardTo string) error`
  - Convenience: gets existing rules, appends the new one, sets all rules
- `RemoveForwarding(ctx context.Context, domain, mailbox string) error`
  - Convenience: gets existing rules, removes the matching one, sets remaining rules

### ForwardingRule type
```go
type ForwardingRule struct {
    Mailbox   string // e.g. "john.doe" (the part before @)
    ForwardTo string // e.g. "shared@gmail.com"
}
```

### XML parsing
- Namecheap returns XML responses with status, errors, and command results
- Parse the XML response, check for API errors, return Go errors with context
- Use `encoding/xml` — no external XML library

### Testing
- Test XML response parsing (success and error cases) with hardcoded XML fixtures
- Test AddForwarding merge logic (existing rules preserved)
- Test RemoveForwarding (correct rule removed, others preserved)
- Tests should NOT call the real Namecheap API — use an `httptest.Server` returning canned XML
- Test error handling (API errors, network errors, malformed responses)

## Target Repo
zarlcorp/zburn

## Agent Role
backend

## Files to Modify
- internal/namecheap/client.go
- internal/namecheap/client_test.go

## Notes
Independent — no dependency on other specs. The client takes Config as a parameter; credential storage in zstore is handled by the settings TUI (spec 104). Domains must use Namecheap's BasicDNS (not custom nameservers) for email forwarding to work — up to 100 forwarding addresses per domain.
