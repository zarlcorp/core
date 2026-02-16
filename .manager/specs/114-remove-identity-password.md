# 114: Remove password field from Identity

## Objective
Remove the top-level `Password` field from the `Identity` struct. Credentials now own per-site passwords, making the identity-level password redundant.

## Context
Before the credential vault (spec 103), each identity had a generated password as part of the persona. Now that credentials store per-site passwords with their own generation (Ctrl+G in credential form), the identity password serves no purpose and confuses the domain model.

The `Password()` method on the generator MUST be kept — it's used by the credential TUI form for generating passwords. Only the identity struct field and its usage in `Generate()` are removed.

Existing stored identities may have a `password` key in their JSON. Go's JSON decoder silently ignores unknown fields, so no migration is needed.

## Requirements

### 1. Remove field from struct
In `internal/identity/identity.go`, remove `Password string` from the `Identity` struct.

### 2. Remove from generator
In `internal/identity/generator.go`, remove the `Password: g.Password(...)` line from `Generate()`. Do NOT remove the `Password()` method itself — it's used by credential form.

### 3. Remove from TUI views
- `internal/tui/generate.go` — remove password from the fields displayed (should go from 10 to 9 fields)
- `internal/tui/detail.go` — remove password from the identity detail display

### 4. Remove from CLI output
In `internal/cli/cli.go`:
- `printIdentity()` — remove the password line
- The JSON output via `printJSON(id)` will automatically stop including password since the struct field is gone

### 5. Update smoke test
In `test/smoke.sh`, test 3 checks `identity --json` for expected fields. Remove `password` from the field list.

### 6. Update all tests
Find and update any tests that reference `Identity.Password` or expect a password field in identity JSON output. This includes:
- `internal/tui/tui_test.go`
- `internal/tui/integration_test.go`
- `internal/identity/` tests (if any reference Password)
- Any other test files

### Acceptance criteria
- `Identity` struct has no `Password` field
- `Generator.Password()` method still exists and works
- Generate view shows 9 fields
- `identity --json` CLI output has no `password` key
- All tests pass

## Target Repo
zarlcorp/zburn

## Agent Role
backend

## Files to Modify
- internal/identity/identity.go
- internal/identity/generator.go
- internal/tui/generate.go
- internal/tui/detail.go
- internal/cli/cli.go
- test/smoke.sh
- internal/tui/tui_test.go
- internal/tui/integration_test.go

## Notes
This is a straightforward field removal. The key risk is missing a test or display location that still references the password. The agent should grep for `Password` across the repo to catch all references.
