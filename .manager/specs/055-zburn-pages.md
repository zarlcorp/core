# 055: zburn GitHub Pages

## Objective
Create a GitHub Pages site for zburn with a single-viewport landing page and a docs page.

## Context
zburn v0.1.0 is released. It has a full CLI and interactive TUI. The page needs to showcase the tool and provide usage documentation.

## Requirements

### Landing page (`docs/index.html`)
Must fit in a single viewport — no scrolling.
- Tool name "zburn" in monospace cyan
- One-liner: "Disposable identities — never give a service your real information again."
- OS-aware install command (brew for macOS, curl for Linux)
- Key features as compact bullets or icons:
  - Generate burner emails, names, addresses, passwords
  - Encrypted local storage
  - Interactive TUI + scriptable CLI
  - Copy to clipboard
- Link to docs page
- Link back to zarlcorp.github.io
- Footer: "Open source. MIT licensed."

### Docs page (`docs/docs.html`)
- Usage section:
  - TUI mode: description of the interactive experience (password prompt, menu, generate, browse, detail, clipboard)
  - CLI commands with examples:
    - `zburn email` — generate a burner email
    - `zburn identity` — generate a complete identity
    - `zburn identity --json` — JSON output
    - `zburn identity --save` — save to encrypted store
    - `zburn list` — list saved identities
    - `zburn list --json` — JSON output
    - `zburn forget <id>` — delete a saved identity
    - `zburn version` — print version
- Data storage: explain that identities are encrypted at rest with AES-256-GCM, master password with Argon2id key derivation, stored in `~/.local/share/zburn/`
- Link back to landing page and org site

### Styling
- Import shared CSS from `https://zarlcorp.github.io/shared.css`
- Minimal page-specific CSS in `docs/style.css` if needed
- Same dark theme, same typography, same color palette

### GitHub Pages setup
- Pages served from `docs/` folder on main branch
- The agent should note in a blocker if Pages needs to be enabled manually (it does — via repo settings)

## Target Repo
zarlcorp/zburn

## Agent Role
frontend

## Files to Create
- docs/index.html — landing page
- docs/docs.html — documentation page
- docs/style.css — page-specific overrides (if needed)

## Dependencies
- 053 (shared CSS design system)

## Notes
- Read the current zburn README.md, cli.go, and main.go for accurate command documentation
- The landing page MUST fit in one viewport — be ruthless about what goes above the fold
- The docs page can scroll — it's documentation
