# 054: Org site redesign with full manifesto

## Objective
Redesign the zarlcorp.github.io landing page to include the full manifesto and OS-aware install commands.

## Context
The current org site has a manifesto excerpt and links to the full markdown. The redesign renders the full manifesto on the page and adds OS-aware install commands. Tool cards should link to the tool's GitHub Pages site (e.g. `zarlcorp.github.io/zburn`), not just the GitHub repo.

## Requirements

### Landing section (top of page)
- `zarlcorp` logo in monospace cyan
- Tagline: "Tools that fight back." in orange
- Tool cards linking to `https://zarlcorp.github.io/<tool>` for each tool:
  - zburn — "Disposable identities — never give a service your real information again."
  - zvault — "Encrypted local storage for secrets, keys, notes. Your data, your machine, your keys."
  - zshield — "DNS-level tracker and ad blocking. See what's tracking you, then kill it."
- OS-aware install command for zburn (the shipped tool)
- Footer with "Open source. MIT licensed." and GitHub link

### Manifesto section
- Full manifesto rendered as HTML below the landing section
- This section scrolls (the landing section above fits in one viewport)
- Style headings, paragraphs, code blocks, tables, horizontal rules
- Use the shared CSS doc-content styles
- Content source: render the full MANIFESTO.md content as HTML
  - Headings: h1-h4 with appropriate styling
  - Code blocks: monospace with bg-card background
  - Tables: clean bordered tables
  - Horizontal rules: subtle separators
  - Block quotes: left-bordered muted text

### Use shared CSS
- Import `shared.css` for all base styles
- Page-specific overrides in a minimal `style.css`

## Target Repo
zarlcorp/zarlcorp.github.io

## Agent Role
frontend

## Files to Modify
- index.html — full redesign
- style.css — page-specific overrides only

## Dependencies
- 053 (shared CSS design system)

## Notes
- The manifesto is ~460 lines of markdown. Render it as static HTML — no markdown parser at runtime, just hand-convert the key sections.
- The manifesto content should be current as of MANIFESTO.md in zarlcorp/core. Read it and convert.
- Keep the landing section above the fold — manifesto starts below
