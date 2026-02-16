# 116: Simplify Namecheap settings

## Objective
Reduce the Namecheap settings form from 5 fields to 2 (username + API key), auto-detect client IP via DNS, and fetch domains from the API instead of manual entry.

## Context
The current settings form asks for: api_user, api_key, username, client_ip, domains. In practice, api_user and username are always the same value. Client IP can be auto-detected. Domains should be fetched from the Namecheap API, not entered manually.

Discussion decided:
- Collapse `api_user` + `username` into one `username` field
- Auto-detect client IP via DNS lookup (`myip.opendns.com @resolver1.opendns.com`)
- Fetch domains via `domains.getList` API on save (validation) and on demand
- Cache fetched domains in settings JSON

## Requirements

### 1. DNS-based IP detection
- Add a function to resolve public IP via DNS: query `myip.opendns.com` against `resolver1.opendns.com` (208.67.222.222)
- Use `net.Resolver` with a custom dialer pointing at the OpenDNS resolver
- No HTTP dependency — pure DNS
- Put this in `internal/namecheap/ip.go`
- Return the IP as a string, error if lookup fails

### 2. Simplify namecheap.Config
- Change `Config` to only require `Username` and `APIKey`
- Remove `APIUser` and `ClientIP` fields
- `baseParams` sends `Username` as both `ApiUser` and `UserName`
- `baseParams` calls the IP detection function to populate `ClientIp`
- Cache the detected IP on the Client struct so we don't do a DNS lookup on every call — detect once on first use

### 3. Add ListDomains method
- Add `ListDomains(ctx context.Context) ([]string, error)` to Client
- Calls `namecheap.domains.getList` API command
- Parse the XML response to extract domain names
- API docs: https://www.namecheap.com/support/api/methods/domains/get-list/
- Return just the domain name strings (e.g. `["example.com", "other.io"]`)

### 4. Simplify NamecheapSettings struct
- Remove `APIUser`, `ClientIP`, and `Domains` fields from `NamecheapSettings`
- Add `CachedDomains []string \`json:"cached_domains"\`` for the fetched domain list
- Keep only `Username` and `APIKey` as user-entered fields
- Update `Configured()` to check `Username` and `APIKey`
- Update `NamecheapConfig()` to return the simplified `namecheap.Config`

### 5. Simplify settings TUI form
- Reduce `settings_namecheap.go` from 5 fields to 2: username, api key
- On save (ctrl+s or enter on last field):
  1. Build a namecheap client from the entered credentials
  2. Call `ListDomains` to validate credentials and fetch domains
  3. If successful: store settings with `CachedDomains` populated, show flash "saved — N domains found"
  4. If failed: show error flash, don't save
- Remove all `ncClientIP` and `ncDomains` field references

### 6. Update tests
- Add test for DNS IP detection (mock the DNS server with a local UDP listener if practical, or just test the function signature and error handling)
- Add test for `ListDomains` XML parsing using `httptest.Server` with canned XML
- Update `settings_test.go` for the simplified struct
- Update any integration tests that reference the old field names

## Target Repo
zarlcorp/zburn

## Agent Role
backend

## Files to Modify
- internal/namecheap/client.go (simplify Config, add ListDomains, cache IP)
- internal/namecheap/ip.go (new — DNS IP detection)
- internal/namecheap/client_test.go (add ListDomains tests, IP tests)
- internal/tui/config.go (simplify NamecheapSettings)
- internal/tui/settings_namecheap.go (2-field form, validation on save)
- internal/tui/settings_test.go (update for new struct)

## Notes
The `ListDomains` API call uses `namecheap.domains.getList` and returns XML with `<Domain>` elements. Each has a `Name` attribute. Pagination may be needed for large accounts but can be deferred — most users have fewer than 20 domains.
