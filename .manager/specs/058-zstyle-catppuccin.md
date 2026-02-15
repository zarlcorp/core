# 058: zstyle Catppuccin Mocha palette

## Objective
Replace the current ad-hoc color palette in zstyle with Catppuccin Mocha. Add per-tool accent colors and expand the style presets for the riced aesthetic.

## Context
zarlcorp is adopting Catppuccin Mocha as the unified color palette across web and TUI. This is the foundation — shared.css and all TUI apps will derive their colors from these definitions.

The unixporn/riced aesthetic means: rounded borders, floating containers with gaps, monospaced everything, tasteful warm pastels on dark backgrounds.

## Requirements

### Replace colors.go with Catppuccin Mocha palette

Full Catppuccin Mocha base:
- `Base = #1e1e2e` — primary background
- `Mantle = #181825` — darker background
- `Crust = #11111b` — darkest background
- `Surface0 = #313244` — elevated surface
- `Surface1 = #45475a` — borders, separators
- `Surface2 = #585b70` — inactive elements
- `Overlay0 = #6c7086` — muted text
- `Overlay1 = #7f849c` — secondary text
- `Overlay2 = #9399b2`
- `Subtext0 = #a6adc8`
- `Subtext1 = #bac2de`
- `Text = #cdd6f4` — primary text

Catppuccin Mocha accents:
- `Rosewater = #f5e0dc`
- `Flamingo = #f2cdcd`
- `Pink = #f5c2e7`
- `Mauve = #cba6f7`
- `Red = #f38ba8`
- `Maroon = #eba0ac`
- `Peach = #fab387`
- `Yellow = #f9e2af`
- `Green = #a6e3a1`
- `Teal = #94e2d5`
- `Sky = #89dceb`
- `Sapphire = #74c7ec`
- `Blue = #89b4fa`
- `Lavender = #b4befe`

### Per-tool accent colors

Define tool-specific accent colors. The designer should pick what looks best, but suggested starting point:
- `Zburn` accent — a warm color (Peach, Maroon, or Red)
- `Zvault` accent — a cool/regal color (Mauve, Lavender, or Blue)
- `Zshield` accent — a protective/calm color (Teal, Sapphire, or Green)

Export these as named variables (e.g. `ZburnAccent`, `ZvaultAccent`, `ZshieldAccent`).

### Semantic color mapping

Map semantic colors to Catppuccin accents:
- `Success` → Green
- `Error` → Red
- `Warning` → Yellow
- `Info` → Blue

Keep the old semantic names (`Success`, `Error`, etc.) so existing code doesn't break. Remove `Cyan`, `Orange`, `Muted`, `Subtle`, `Bright` and replace with Catppuccin equivalents. Provide backward-compatible aliases if needed during transition.

### Update styles.go

Update all lipgloss style presets to use the new palette:
- `Title` — tool accent color (default to Lavender or Blue), bold
- `Subtitle` — Subtext1, bold
- `Highlight` — Peach or tool accent
- `MutedText` — Overlay1
- `Border` — Surface1 foreground, rounded
- `ActiveBorder` — tool accent foreground, rounded
- `StatusOK/Err/Warn` — mapped to semantic colors

### Update tests

Update zstyle_test.go to verify the new color values.

### CSS export helper

Add a function or constant block that exports the palette as CSS custom properties string. This allows shared.css to reference the exact same hex values. Something like:

```go
const CSSVariables = `
:root {
  --base: #1e1e2e;
  --mantle: #181825;
  ...
}
`
```

This isn't strictly required but is useful for keeping web and TUI in sync.

## Target Repo
zarlcorp/core

## Agent Role
backend

## Files to Modify
- pkg/zstyle/colors.go — full palette replacement
- pkg/zstyle/styles.go — updated presets
- pkg/zstyle/zstyle_test.go — updated tests

## Dependencies
None — this is the foundation.

## Notes
- The Catppuccin project is at github.com/catppuccin/catppuccin — reference for exact hex values
- Mocha is the darkest variant, which fits dark terminal backgrounds
- The CSS export is a nice-to-have, not a blocker
- Keep the package API surface small — colors + styles + keys, nothing else
