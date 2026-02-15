# 080: Fix first-run password prompt

## Objective
When zburn is run for the first time, the password prompt says "master password:" — identical to subsequent runs. The user has no way to know they're creating a new password vs entering an existing one.

## Context
The `passwordModel` has a `firstRun` flag that controls whether confirmation is required, but the initial prompt text doesn't change. The View function always shows "master password:" for the first entry.

## Requirements

### Fix in `internal/tui/password.go`

Update the View function to distinguish between first run and subsequent runs:

- **First run, initial entry:** `"create master password:"`
- **First run, confirmation:** `"confirm password:"` (already correct)
- **Subsequent runs:** `"master password:"` (already correct)

The logic should check both `m.firstRun` and `m.confirming`:

```go
var prompt string
if m.firstRun {
    if m.confirming {
        prompt = "confirm password:"
    } else {
        prompt = "create master password:"
    }
} else {
    prompt = "master password:"
}
```

### Tests

Add or update tests in `internal/tui/password_test.go` that verify:
1. First-run View contains "create master password"
2. First-run confirming View contains "confirm password"
3. Non-first-run View contains "master password" (not "create")

Run `go test ./...` to verify.

## Target Repo
zarlcorp/zburn

## Agent Role
backend

## Files to Modify
- internal/tui/password.go
- internal/tui/password_test.go

## Notes
Small change. Depends on 079 being merged first only if we want to do a combined release, but can be developed independently.
