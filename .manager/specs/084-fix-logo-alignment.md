# 084: Fix logo alignment in zburn TUI

## Objective
The logo appears misaligned in the TUI because `fmt.Sprintf("  %s", logo)` only indents the first line of the multiline logo string. Lines 2-3 start at column 0.

## Context
Both `password.go` and `menu.go` render the logo with:
```go
s := fmt.Sprintf("\n  %s\n  %s\n\n", logo, toolName)
```

The `  %s` adds a 2-space indent before the first line of the logo, but the logo is a 3-line string with embedded newlines. Lines 2 and 3 get no indent, causing visible misalignment.

## Requirements

### Fix rendering in `password.go` and `menu.go`

Use lipgloss to indent the logo block instead of `fmt.Sprintf` string interpolation. The logo is already rendered through `zstyle.StyledLogo()` which returns a styled multiline string — apply a `MarginLeft(2)` or `PaddingLeft(2)` to the style passed to StyledLogo, or wrap the result in a lipgloss style with left padding.

Example fix in password.go:
```go
logoStyle := lipgloss.NewStyle().Foreground(zstyle.ZburnAccent).PaddingLeft(2)
logo := zstyle.StyledLogo(logoStyle)
```

Or wrap after:
```go
logo := lipgloss.NewStyle().MarginLeft(2).Render(
    zstyle.StyledLogo(lipgloss.NewStyle().Foreground(zstyle.ZburnAccent)),
)
```

Either approach works — pick whichever is cleaner. The key requirement is that ALL lines of the logo get the same left indent.

Apply the same fix in menu.go.

### Update zstyle dependency

Update go.mod to use `pkg/zstyle v0.4.0` (which has the lowercase logo from 083). Run `go mod tidy`.

### Test

Run `go test ./...` to verify everything compiles and passes.

## Target Repo
zarlcorp/zburn

## Agent Role
backend

## Files to Modify
- internal/tui/password.go (logo rendering)
- internal/tui/menu.go (logo rendering)
- go.mod (zstyle version bump)

## Notes
Depends on 083 (lowercase logo) being tagged as pkg/zstyle/v0.4.0 on core first. Read every file before editing. Don't change functionality — only fix the visual alignment and pick up the new logo.
