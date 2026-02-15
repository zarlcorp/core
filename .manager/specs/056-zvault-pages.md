# 056: zvault GitHub Pages

## Objective
Create a GitHub Pages site for zvault with a single-viewport landing page and minimal docs.

## Context
zvault is scaffolded but not yet built. The page serves as a placeholder that establishes the tool's identity and links to the repo.

## Requirements

### Landing page (`docs/index.html`)
Must fit in a single viewport — no scrolling.
- Tool name "zvault" in monospace cyan
- One-liner: "Encrypted local storage for secrets, keys, notes. Your data, your machine, your keys."
- Status badge or note: "Coming soon" in muted text
- Planned features as compact bullets:
  - Store secrets encrypted at rest
  - Master password with Argon2id key derivation
  - Organize with tags and folders
  - Search across entries
  - Auto-lock after inactivity
  - Export/import (encrypted format only)
- Link to GitHub repo
- Link back to zarlcorp.github.io
- Footer: "Open source. MIT licensed."

### Docs page (`docs/docs.html`)
- Planned CLI commands (not yet implemented):
  - `zvault` — interactive TUI
  - `zvault get <path>` — retrieve a secret
  - `zvault set <path>` — store a secret
  - `zvault search <query>` — find entries
- Note that these are planned, not live yet
- Link back to landing page and org site

### Styling
- Import shared CSS from `https://zarlcorp.github.io/shared.css`
- Same design as zburn pages — consistent tool page template

## Target Repo
zarlcorp/zvault

## Agent Role
frontend

## Files to Create
- docs/index.html — landing page
- docs/docs.html — documentation page (minimal)
- docs/style.css — page-specific overrides (if needed)

## Dependencies
- 053 (shared CSS design system)

## Notes
- This is largely a template copy of zburn's page structure with different content
- No install command — tool isn't released yet
- "Coming soon" aesthetic — clean, not empty
