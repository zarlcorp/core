# 063: zshield pages redesign

## Objective
Redesign the zshield GitHub Pages site with the new riced aesthetic and tool-specific accent color.

## Context
zshield is not yet released. The page is a "Coming soon" placeholder with planned features and CLI commands.

## Requirements

### Landing page (`docs/index.html`)
Must fit in a single viewport.
- Status bar navigation at top
- "zshield" in tool accent color, monospaced
- One-liner: "DNS-level tracker and ad blocking. See what's tracking you, then kill it."
- "Coming soon" status indicator
- Planned features in a floating window
- Links to docs page, org site, GitHub
- No install command (not released)
- Footer

### Docs page (`docs/docs.html`)
- Status bar navigation
- Planned CLI commands in floating window sections
- Note that commands are planned, not live
- Uses `.doc-content` within `.window` containers

### Tool accent color
Override the shared CSS accent variable for zshield:
```css
:root { --accent: var(--teal); } /* or whatever the designer chose */
```

### Styling
- Import shared CSS from `https://zarlcorp.github.io/shared.css`
- Minimal page-specific CSS in `docs/style.css`

## Target Repo
zarlcorp/zshield

## Agent Role
frontend

## Files to Modify
- docs/index.html — redesigned landing page
- docs/docs.html — redesigned docs page
- docs/style.css — accent color override + page-specific styles

## Dependencies
- 059 (shared CSS riced redesign)

## Notes
- Same template structure as zvault but different content and accent color
- "Coming soon" aesthetic
- No install commands
