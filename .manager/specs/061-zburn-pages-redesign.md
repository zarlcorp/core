# 061: zburn pages redesign

## Objective
Redesign the zburn GitHub Pages site with the new riced aesthetic and tool-specific accent color.

## Context
The current zburn pages are functional but visually plain. The redesign applies the new shared.css design system with zburn's accent color.

## Requirements

### Landing page (`docs/index.html`)
Must fit in a single viewport.
- Status bar navigation at top
- "zburn" in tool accent color, monospaced
- One-liner: "Disposable identities — never give a service your real information again."
- Key features in a floating window container
- OS-aware install command (brew for macOS, curl for Linux)
- Links to docs page, org site, GitHub
- Footer

### Docs page (`docs/docs.html`)
- Status bar navigation
- Content organized in floating window sections:
  - Interactive TUI — description of the TUI experience
  - CLI Commands — all commands with examples
  - Data Storage — encryption details
- Uses `.doc-content` within `.window` containers
- Scrollable

### Tool accent color
Override the shared CSS accent variable for zburn:
```css
:root { --accent: var(--peach); } /* or whatever the designer chose in 058/059 */
```
The exact color will be established in 058 (zstyle) and 059 (shared.css).

### Styling
- Import shared CSS from `https://zarlcorp.github.io/shared.css`
- Import install.js from `https://zarlcorp.github.io/install.js`
- Minimal page-specific CSS in `docs/style.css`

## Target Repo
zarlcorp/zburn

## Agent Role
frontend

## Files to Modify
- docs/index.html — redesigned landing page
- docs/docs.html — redesigned docs page
- docs/style.css — accent color override + page-specific styles

## Dependencies
- 059 (shared CSS riced redesign)

## Notes
- Read the current zburn README.md for accurate command documentation
- Read the redesigned shared.css to understand new class names and patterns
- The floating window aesthetic is key — sections should feel like terminal windows
