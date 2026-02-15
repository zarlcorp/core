# 069: TUI lowercase menus

## Objective
Make all zburn TUI text lowercase to match the riced aesthetic.

## Context
The zarlcorp visual identity is all-lowercase — nav bars, window titles, page headers. The TUI should match. Currently menu items and field labels use title case.

## Requirements

### Menu items (`internal/tui/menu.go`)
```
Current:                          Target:
"Generate identity"        →     "generate identity"
"Generate email (quick)"   →     "generate email (quick)"
"Browse saved identities"  →     "browse saved identities"
"Quit"                     →     "quit"
```

### Field labels (`internal/tui/generate.go`)
```
Current:     Target:
"Name"   →  "name"
"Email"  →  "email"
"Phone"  →  "phone"
"Street" →  "street"
"City"   →  "city"
"State"  →  "state"
"Zip"    →  "zip"
"Password" → "password"
```

### List header columns (`internal/tui/list.go`)
```
Current:                          Target:
"ID", "Name", "Email", "Created" → "id", "name", "email", "created"
```

### Tests
Update any test assertions that check for the old capitalized strings. All existing tests must pass.

## Target Repo
zarlcorp/zburn

## Agent Role
backend

## Files to Modify
- `internal/tui/menu.go` — menu item strings
- `internal/tui/generate.go` — field label strings
- `internal/tui/list.go` — header column strings
- `internal/tui/tui_test.go` — update test assertions for new casing

## Dependencies
None

## Notes
- This is a string-only change — no logic changes
- Don't change the detail.go field labels separately, they reference the same fields struct from generate.go
- Run `go test ./...` to verify
- Run `go build ./cmd/zburn/` to verify build
