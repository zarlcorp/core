# 062: zvault pages redesign

## Objective
Redesign the zvault GitHub Pages site with the new riced aesthetic and tool-specific accent color.

## Context
zvault is not yet released. The page is a "Coming soon" placeholder with planned features and CLI commands.

## Requirements

### Landing page (`docs/index.html`)
Must fit in a single viewport.
- Status bar navigation at top
- "zvault" in tool accent color, monospaced
- One-liner: "Encrypted local storage for secrets, keys, notes. Your data, your machine, your keys."
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
Override the shared CSS accent variable for zvault:
```css
:root { --accent: var(--mauve); } /* or whatever the designer chose */
```

### Styling
- Import shared CSS from `https://zarlcorp.github.io/shared.css`
- Minimal page-specific CSS in `docs/style.css`

## Target Repo
zarlcorp/zvault

## Agent Role
frontend

## Files to Modify
- docs/index.html — redesigned landing page
- docs/docs.html — redesigned docs page
- docs/style.css — accent color override + page-specific styles

## Dependencies
- 059 (shared CSS riced redesign)

## Notes
- Same template structure as zburn but different content and accent color
- "Coming soon" aesthetic — clean, minimal, not empty
- No install commands since tool isn't released
