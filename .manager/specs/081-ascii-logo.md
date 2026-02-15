# 081: Add zarlcorp ASCII logo to zstyle

## Objective
Add a shared ASCII art "zarlcorp" logo to `core/pkg/zstyle` that all TUI tools can display on their password/splash screen. Consistent branding across the tool suite.

## Context
Currently each tool shows its own name as a plain styled title (e.g., `zstyle.Title.Render("zburn")`). The user wants a proper ASCII art zarlcorp logo that appears on the first screen (password entry) of every tool.

## Requirements

### Add logo to `pkg/zstyle`

Create a `Logo` variable in zstyle that renders an ASCII art "zarlcorp" text. Requirements:
- Clean, minimal ASCII art — not overly large
- Should fit comfortably in an 80-column terminal
- Use the Catppuccin Mocha accent color (the tool's `--accent` equivalent, or the default accent)
- Exported as a pre-rendered string or a function that takes a lipgloss style

Suggested approach — a simple blocky/pixel font, roughly 3-5 lines tall. Something like:
```
               _
 ______ _ _ __| | ___ ___  _ __ _ __
|_  / _` | '__| |/ __/ _ \| '__| '_ \
 / / (_| | |  | | (_| (_) | |  | |_) |
/___\__,_|_|  |_|\___\___/|_|  | .__/
                                |_|
```

Or a cleaner, more minimal style that fits the riced aesthetic. The agent should pick something that looks good in a terminal with JetBrains Mono or a standard monospace font.

### Export it

```go
// Logo is the ASCII art zarlcorp logo for TUI splash screens.
var Logo string
```

And a styled version:
```go
// StyledLogo returns the logo rendered with the given lipgloss style.
func StyledLogo(style lipgloss.Style) string
```

### Tests

Basic test that Logo is not empty and StyledLogo returns a non-empty string.

Run `go test ./pkg/zstyle/...` to verify.

## Target Repo
zarlcorp/core

## Agent Role
backend

## Files to Modify
- pkg/zstyle/logo.go (new file)
- pkg/zstyle/logo_test.go (new file)

## Notes
This just adds the logo to the shared package. The actual integration into each tool's TUI (replacing the plain title with the logo on the password screen) will be done per-tool. For zburn, that's part of 082.
