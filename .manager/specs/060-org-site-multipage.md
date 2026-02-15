# 060: Org site multi-page redesign

## Objective
Redesign the zarlcorp.github.io org site as a multi-page site with the new riced aesthetic. Break the current single-page monolith into landing, manifesto, and architecture pages.

## Context
The current site dumps the entire 468-line manifesto on one page. The redesign breaks this into separate pages, each using the new shared.css riced design system.

## Requirements

### Landing page (`index.html`)
Must fit in a single viewport — this is the storefront.
- Status bar navigation at top (`.bar` style)
- `zarlcorp` logo in accent color, monospaced
- Tagline: "Tools that fight back." in subtext color
- Tool cards as mini floating windows, each with:
  - Tool name in its accent color
  - One-liner description
  - Link to tool's GitHub Pages site
- OS-aware install command for zburn
- Navigation links to manifesto and architecture pages
- Footer with "Open source. MIT licensed." and GitHub link

### Manifesto page (`manifesto.html`)
- Status bar navigation
- Full manifesto content rendered as HTML
- Uses `.doc-content` within `.window` containers
- Break into logical sections with window containers:
  - The Problem
  - The Belief
  - What We Build
  - The Standard
  - Decision Log
  - Founding
- Each major section in its own floating window
- Scrollable — this is documentation

### Architecture page (`architecture.html`)
- Status bar navigation
- Technical content from the manifesto:
  - Platform overview (core repo structure)
  - Package roster and layering
  - Product details (zburn, zvault, zshield)
  - Agent model
  - Release pipeline
- Each section in its own floating window
- Code blocks and ASCII diagrams styled for the terminal aesthetic

### Use shared CSS
- Import the redesigned `shared.css` for all base styles
- Page-specific overrides in `style.css` (minimal)

## Target Repo
zarlcorp/zarlcorp.github.io

## Agent Role
frontend

## Files to Create
- manifesto.html — manifesto page
- architecture.html — architecture page

## Files to Modify
- index.html — redesigned landing page
- style.css — page-specific overrides

## Dependencies
- 059 (shared CSS riced redesign)

## Notes
- Read MANIFESTO.md from /Users/bruno/src/zarlcorp/core/MANIFESTO.md for content
- The manifesto is the soul of the org — it deserves a good presentation
- Architecture content is the technical reference — package diagrams, ASCII art, etc.
- Split the manifesto content intelligently — "The Problem" through "The Standard" on manifesto page, technical details on architecture page
- Tool card links point to `https://zarlcorp.github.io/<tool>`
