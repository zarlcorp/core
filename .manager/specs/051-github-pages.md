# 051: GitHub Pages for all repos

## Objective
Enable GitHub Pages on all zarlcorp repos with a consistent landing page style.

## Context
zarlcorp repos need public-facing pages. GitHub Pages with a simple static site gives each tool a web presence without adding dependencies or build complexity.

## Requirements

### Organization page (zarlcorp/.github)
- Create `profile/README.md` if not already present (org-level profile)
- Create a GitHub Pages site in a new repo `zarlcorp/zarlcorp.github.io` with:
  - Simple landing page listing all tools (zburn, zvault, zshield)
  - Link to each tool's repo and GitHub Pages
  - Manifesto excerpt or link
  - Clean, minimal design — dark theme matching the terminal aesthetic
  - Pure HTML/CSS, no JavaScript, no build tools

### Per-tool pages (zburn, zvault, zshield)
- Enable GitHub Pages via a `docs/` folder or `gh-pages` branch on each repo
- Each tool gets a simple `index.html` with:
  - Tool name and one-line description
  - Install instructions (brew + go install)
  - Link to README for full docs
  - Link back to zarlcorp org page
- Same dark theme as org page

### Implementation approach
- Use a single HTML template with CSS variables for per-tool customization
- No static site generators — plain HTML/CSS files
- GitHub Actions not needed — Pages serves static files from the repo

## Target Repo
zarlcorp/zarlcorp.github.io (new repo for org page)
zarlcorp/zburn, zarlcorp/zvault, zarlcorp/zshield (docs/ folders)

## Agent Role
frontend

## Files to Create
- Org site: index.html, style.css
- Per-tool: docs/index.html, docs/style.css (or shared)

## Notes
- This is the one place where we write HTML — it's a static landing page, not an app
- Keep it minimal — the tools speak for themselves
- Dark background, cyan/orange accents matching zstyle colors
