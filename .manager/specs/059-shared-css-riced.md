# 059: Shared CSS riced redesign

## Objective
Complete overhaul of shared.css with a Catppuccin Mocha palette and unixporn/riced aesthetic — floating window containers, monospaced typography, tiling WM feel.

## Context
The current shared.css is functional but generic. The redesign creates a design system that feels like a beautifully riced desktop — think hyprland with floating bars, gaps between containers, rounded corners, monospaced everything.

This is the CSS foundation all zarlcorp pages import. Tool pages get their own accent color via CSS variable override.

## Requirements

### Color system (Catppuccin Mocha)
Replace all CSS custom properties with Catppuccin Mocha values:
```css
:root {
  --base: #1e1e2e;
  --mantle: #181825;
  --crust: #11111b;
  --surface0: #313244;
  --surface1: #45475a;
  --surface2: #585b70;
  --overlay0: #6c7086;
  --overlay1: #7f849c;
  --overlay2: #9399b2;
  --subtext0: #a6adc8;
  --subtext1: #bac2de;
  --text: #cdd6f4;
  --rosewater: #f5e0dc;
  --flamingo: #f2cdcd;
  --pink: #f5c2e7;
  --mauve: #cba6f7;
  --red: #f38ba8;
  --maroon: #eba0ac;
  --peach: #fab387;
  --yellow: #f9e2af;
  --green: #a6e3a1;
  --teal: #94e2d5;
  --sky: #89dceb;
  --sapphire: #74c7ec;
  --blue: #89b4fa;
  --lavender: #b4befe;
  /* per-tool accent — overridden by tool pages */
  --accent: var(--lavender);
  --accent-dim: var(--overlay2);
}
```

### Typography
- Primary font: monospaced stack everywhere — `"JetBrains Mono", "Fira Code", "Cascadia Code", "SF Mono", Menlo, Consolas, monospace`
- Consider importing JetBrains Mono from Google Fonts for consistent rendering
- Body text in monospace, slightly smaller base size (16px) for the terminal feel
- No sans-serif anywhere — commit fully to the monospaced aesthetic

### Riced container aesthetic

The core visual element: **floating containers** that feel like windows in a tiling WM.

- `.window` — a floating container with:
  - Background: `--surface0`
  - Border: 1px solid `--surface1`, rounded corners (8-12px)
  - Subtle gap/margin between windows (like tiling WM gaps)
  - Optional: subtle box-shadow for depth

- `.window-title` — optional title bar at top of a window:
  - Background: slightly different shade or transparent
  - Could include decorative dots (red/yellow/green circles) like macOS — designer's call on whether this is tasteful or gimmicky
  - Title text in `--subtext1`

- `.window-content` — padding inside windows

### Status bar / info bar
- `.bar` — a horizontal bar element (like polybar/waybar):
  - Background: `--mantle`
  - Monospaced text
  - Can contain status items, navigation, breadcrumbs
  - Sits at top of page as navigation

### Layout
- Page background: `--base`
- Content max-width stays ~800px but with visible gaps between sections
- Sections wrapped in `.window` containers with gaps between them
- Grid layouts for tool cards use gaps to create the tiling effect

### Existing classes to redesign
Keep the same class names but restyle completely:
- `.container` — centered layout with padding
- `.hero` — could be a `.window` with prominent styling
- `.logo` — monospaced, accent color
- `.tagline` — `--subtext1`, slightly muted
- `.tool-card` — redesign as mini windows with tool accent colors
- `.install-cmd` — terminal-style with prompt, fits the aesthetic naturally
- `.footer` — muted bar at bottom
- `.section-heading` — `--overlay1` text
- `.doc-content` — prose styling within windows
- `.feature-list` — styled list with accent-colored bullets
- `.nav` — redesign as a status bar

### Per-tool accent override
Tool pages override the accent with one CSS variable:
```css
:root { --accent: var(--peach); } /* zburn */
:root { --accent: var(--mauve); } /* zvault */
:root { --accent: var(--teal); } /* zshield */
```

All accent-dependent styles (headings, active borders, links, highlights) use `var(--accent)` so they automatically adapt.

### Responsive
- Mobile: single column, windows stack vertically
- Desktop: grid where appropriate
- The riced aesthetic should scale down gracefully

### install.js
Keep the existing OS detection script as-is — it works fine.

## Target Repo
zarlcorp/zarlcorp.github.io

## Agent Role
frontend

## Files to Modify
- shared.css — complete rewrite
- install.js — no changes needed

## Dependencies
- 058 (zstyle Catppuccin palette) — for exact color values and per-tool accent choices

## Notes
- Reference r/unixporn for aesthetic inspiration
- The key visual differentiator is the floating window/container pattern — everything lives in bordered, rounded, gapped containers
- Don't overdo the fake window chrome — it should feel inspired by tiling WMs, not be a pixel-perfect imitation
- The designer has creative latitude on how far to push the rice — subtle vs full
- Test that the palette has sufficient contrast for accessibility (Catppuccin is designed for this)
