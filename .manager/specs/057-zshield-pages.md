# 057: zshield GitHub Pages

## Objective
Create a GitHub Pages site for zshield with a single-viewport landing page and minimal docs.

## Context
zshield is scaffolded but not yet built. The page serves as a placeholder that establishes the tool's identity and links to the repo.

## Requirements

### Landing page (`docs/index.html`)
Must fit in a single viewport — no scrolling.
- Tool name "zshield" in monospace cyan
- One-liner: "DNS-level tracker and ad blocking. See what's tracking you, then kill it."
- Status badge or note: "Coming soon" in muted text
- Planned features as compact bullets:
  - Local DNS resolver blocking trackers and ads
  - Blocklist management (community lists + custom rules)
  - Real-time query log with TUI dashboard
  - Per-domain allow/deny overrides
  - Statistics: blocked vs allowed, top blocked domains
  - Runs as daemon with TUI attach
- Link to GitHub repo
- Link back to zarlcorp.github.io
- Footer: "Open source. MIT licensed."

### Docs page (`docs/docs.html`)
- Planned CLI commands (not yet implemented):
  - `zshield start` — start DNS resolver daemon
  - `zshield` — attach TUI dashboard
  - `zshield status` — show blocking stats
  - `zshield allow <domain>` — whitelist a domain
  - `zshield block <domain>` — blacklist a domain
- Note that these are planned, not live yet
- Link back to landing page and org site

### Styling
- Import shared CSS from `https://zarlcorp.github.io/shared.css`
- Same design as zburn/zvault pages — consistent tool page template

## Target Repo
zarlcorp/zshield

## Agent Role
frontend

## Files to Create
- docs/index.html — landing page
- docs/docs.html — documentation page (minimal)
- docs/style.css — page-specific overrides (if needed)

## Dependencies
- 053 (shared CSS design system)

## Notes
- Same template structure as zvault page — different content
- No install command — tool isn't released yet
- "Coming soon" aesthetic
