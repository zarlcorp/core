# 107: zburn integration tests

## Objective
Add integration tests that exercise multi-step flows through the zburn TUI root model with a real zstore backend. Verify that the new features (credentials, settings, burn cascade) work correctly end-to-end.

## Context
Specs 098-105 added significant new functionality: credential vault, settings for external services, and burn cascade. Each feature has good unit-level tests, but there are no tests that exercise full user flows through the root model with a real encrypted store. These integration tests ensure the features work together correctly.

## Requirements

### Credential lifecycle tests
- Open store → save identity → add credential for that identity → verify credential count in detail view
- Add credential → edit credential → verify changes persisted
- Add multiple credentials → delete one → verify others remain
- Credential isolation: credentials for identity A are not visible when viewing identity B

### Settings persistence tests
- Open store → configure Namecheap settings → reload configs → verify settings persisted
- Configure Gmail settings with mock OAuth → verify GmailConfigured() returns true
- Configure Twilio settings → verify TwilioConfigured() returns true
- Disconnect Gmail → verify GmailConfigured() returns false
- Feature gating: verify Configured() methods match stored state

### Burn cascade integration tests
- Create identity with 3 credentials → burn identity → verify identity deleted → verify all 3 credentials deleted
- Burn identity with no credentials → verify clean burn
- Burn with mocked external services (email forwarder, phone releaser) → verify cascade calls all services
- Burn with failing external service → verify best-effort continues and identity still deleted

### Store lifecycle tests
- Fresh store: open with new password → verify empty collections
- Reopening: open store → save data → close → reopen with same password → verify data persisted
- Wrong password: open store → close → reopen with wrong password → verify error

### Test infrastructure
- Use `t.TempDir()` for each test's data directory
- Create real `zstore.Store` instances (no mocks for the store layer)
- Use the existing test helpers (`keyMsg`, `enterKey`, etc.) from `tui_test.go`
- Tests go in `internal/tui/integration_test.go` to keep them separate from unit tests

## Target Repo
zarlcorp/zburn

## Agent Role
testing

## Depends On
- 098-105 (all merged)

## Files to Modify
- internal/tui/integration_test.go (new — integration test suite)

## Notes
The version is set via ldflags at build time (`var version = "dev"`), so no version string changes needed. Focus on testing the actual user flows, not re-testing what unit tests already cover. Use real zstore instances — the whole point is verifying the integration between TUI, zstore, and the business logic packages.
