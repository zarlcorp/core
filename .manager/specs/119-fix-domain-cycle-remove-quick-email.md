# 119: Fix domain cycle and remove quick email menu option

## Objective
Two small TUI fixes: (1) cycling domains should only change the email, not regenerate the entire identity, and (2) remove the broken "generate email (quick)" menu option.

## Context
Spec 118 wired domain rotation into the generate view, but `handleCycleDomain` calls `m.gen.Generate(domain)` which creates an entirely new identity. It should only update the email field using the existing name + new domain.

The "generate email (quick)" menu option calls `handleQuickEmail` which generates a standalone email. With name-based email patterns (spec 117), this doesn't belong as a top-level menu item — users should generate a full identity instead.

## Requirements

### 1. Fix handleCycleDomain
In `tui.go` `handleCycleDomain()`, instead of:
```go
id := m.gen.Generate(domain)
m.generate = newGenerateModel(id, domain)
```

Do:
```go
id := m.generate.identity
id.Email = m.gen.Email(id.FirstName, id.LastName, domain)
m.generate = newGenerateModel(id, domain)
```

Keep the same identity (name, phone, address, DOB), only change the email.

### 2. Remove quick email from TUI menu
- Remove `menuEmail` from the `menuChoice` enum in `menu.go`
- Remove `"generate email (quick)"` from `menuItems`
- Remove the `case menuEmail:` in `selectItem()`
- Remove `quickEmailMsg` type from `menu.go`
- Remove `case quickEmailMsg:` handler from `tui.go` `Update()`
- Remove `handleQuickEmail()` method from `tui.go`
- Remove or update any tests that reference `quickEmailMsg` or `handleQuickEmail`

### 3. Update tests
- Fix `TestCycleDomainRegeneratesIdentity` — should now verify the name STAYS the same and only the email domain changes
- Remove/update the test that checks `quickEmailMsg` is emitted from menu selection
- Ensure all remaining tests pass

## Target Repo
zarlcorp/zburn

## Agent Role
backend

## Files to Modify
- internal/tui/tui.go (fix handleCycleDomain, remove quickEmailMsg handler and handleQuickEmail)
- internal/tui/menu.go (remove menuEmail, quickEmailMsg, menu item)
- internal/tui/tui_test.go (fix cycle test, remove quickEmail test)
- internal/tui/integration_test.go (if any reference quickEmailMsg)

## Notes
The CLI command `zburn email` (`CmdEmail` in `cli.go`) is NOT being removed — only the TUI menu option.
