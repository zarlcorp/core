# 028: Build zstyle package

## Objective
Implement the `zstyle` package — zarlcorp's visual identity for TUIs. Provides a color palette, lipgloss style presets, and standard keybinding constants. Not a Bubble Tea wrapper — just constants, styles, and helpers that tools import for visual consistency.

## Context
zstyle sits in the presentation layer of the package hierarchy, alongside zcrypto. It depends only on the Charmbracelet ecosystem (lipgloss for styles, bubbles/key for keybinding types). Every zarlcorp TUI tool will import zstyle so they all look and feel the same.

The `go.mod` already exists at `pkg/zstyle/go.mod` with module path `github.com/zarlcorp/core/pkg/zstyle`. The module is already listed in `go.work`.

## Requirements

### Color palette — `colors.go`

Theme: **cyan and orange on a dark terminal background.** Assume the user's terminal has a dark background — never set background colors on the main content area.

Define colors as `lipgloss.Color` constants:

```go
// core brand colors
Cyan      = lipgloss.Color("#00E5FF")  // primary — headings, active elements, highlights
Orange    = lipgloss.Color("#FF6E40")  // accent — calls to action, selected items, emphasis

// semantic colors
Success   = lipgloss.Color("#69F0AE")  // green — completed, passed, ok
Error     = lipgloss.Color("#FF5252")  // red — errors, failures, destructive actions
Warning   = lipgloss.Color("#FFD740")  // amber — warnings, caution
Info      = lipgloss.Color("#40C4FF")  // light blue — informational

// neutral tones
Muted     = lipgloss.Color("#78909C")  // grey-blue — secondary text, timestamps, metadata
Subtle    = lipgloss.Color("#37474F")  // dark grey — borders, separators, inactive elements
Bright    = lipgloss.Color("#ECEFF1")  // near-white — primary text when emphasis needed
```

The agent should use these exact hex values as starting points. Minor tweaks for readability on dark terminals are acceptable but document any changes.

### Style presets — `styles.go`

Pre-built lipgloss styles that tools import directly. All styles assume dark background.

```go
// text styles
Title     lipgloss.Style  // bold, Cyan foreground
Subtitle  lipgloss.Style  // bold, Muted foreground
Highlight lipgloss.Style  // Orange foreground
Muted     lipgloss.Style  // Muted foreground (secondary text)

// status indicators
StatusOK    lipgloss.Style  // Success foreground
StatusErr   lipgloss.Style  // Error foreground
StatusWarn  lipgloss.Style  // Warning foreground

// structural
Border      lipgloss.Style  // rounded border in Subtle color
ActiveBorder lipgloss.Style // rounded border in Cyan
```

Styles are package-level variables (not constants — lipgloss.Style isn't const-able). Use `lipgloss.NewStyle()` to build them.

### Keybinding constants — `keys.go`

Standard keybindings using `github.com/charmbracelet/bubbles/key`. Define `key.Binding` values:

```go
KeyQuit     key.Binding  // q, ctrl+c — quit the application
KeyHelp     key.Binding  // ? — toggle help
KeyUp       key.Binding  // k, up — navigate up
KeyDown     key.Binding  // j, down — navigate down
KeyEnter    key.Binding  // enter — confirm/select
KeyBack     key.Binding  // esc, backspace — go back
KeyTab      key.Binding  // tab — next field/section
KeyFilter   key.Binding  // / — search/filter
```

Each binding should have a `Help()` that returns a short description (e.g. `key.WithHelp("q", "quit")`).

### What zstyle does NOT do
- No Bubble Tea models or update functions — this is not a component library
- No runtime configuration — colors are compile-time constants
- No terminal detection or adaptive themes — assume dark background
- No dependency on other zarlcorp packages — lipgloss and bubbles only

### Package doc
Brief package comment on the main file explaining zstyle provides the zarlcorp visual identity. Show a quick usage example:

```go
import "github.com/zarlcorp/core/pkg/zstyle"

fmt.Println(zstyle.Title.Render("My Tool"))
fmt.Println(zstyle.StatusOK.Render("✓ done"))
```

### Tests — `zstyle_test.go`

Since this is mostly constants and preset styles, tests should verify:
- All color constants are non-empty
- All style presets render without panicking
- All keybindings have non-empty key sets and help text
- No nil values in any exported variable

Table-driven tests grouping by category (colors, styles, keys).

## Acceptance Criteria
1. Color palette defines all 9 colors as `lipgloss.Color` values
2. Style presets cover text (4), status (3), and structural (2) categories
3. Keybinding constants define all 8 standard bindings with help text
4. All tests pass (`go test ./...` from `pkg/zstyle/`)
5. `go.mod` has lipgloss and bubbles dependencies
6. No dependency on any other zarlcorp package
7. Package compiles cleanly with no unused imports

## Target Repo
zarlcorp/core

## Agent Role
backend

## Files to Create/Modify
- `pkg/zstyle/colors.go` — color palette constants
- `pkg/zstyle/styles.go` — lipgloss style presets
- `pkg/zstyle/keys.go` — standard keybinding constants
- `pkg/zstyle/zstyle_test.go` — tests for all exported values
- `pkg/zstyle/go.mod` — add lipgloss and bubbles dependencies

## Dependencies
None — zstyle is independent of other zarlcorp packages.

## Notes
- Use `lipgloss` v1.1.0 and `bubbles` v1.0.0 (both at stable v1)
- Colors are `lipgloss.Color` (a string type), not `lipgloss.AdaptiveColor` — we assume dark terminals
- Style variables use a `Style` suffix only where there's ambiguity with color names. For `Muted`, the style should be named `MutedStyle` to avoid collision with the `Muted` color constant, or use a different approach (e.g. colors in a `Color` struct/namespace, styles as bare names). Agent should pick the cleanest approach that avoids name collisions.
- Keybindings use `key.NewBinding` with `key.WithKeys(...)` and `key.WithHelp(...)`
- Run `go mod tidy` after adding dependencies to pick up transitive deps
