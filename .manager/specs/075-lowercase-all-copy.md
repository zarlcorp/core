# 075: Lowercase all website copy

## Objective
All text on zarlcorp pages should be lowercase — headings, body copy, taglines, feature lists, footers. The riced aesthetic extends to text. Only technical acronyms and proper nouns stay uppercase.

## Context
Voice guide rule 9: "all lowercase in headers, nav, labels." The user wants this extended to ALL copy on the website, not just headers. Currently page titles, headings, and sentence starts all use standard capitalization.

## Requirements

### What to lowercase
- Page titles (`<h1>`): "Manifesto" → "manifesto", "Architecture" → "architecture", "Decision Log" → "decision log"
- Section headings (`<h2>`, `<h3>`, `<h4>`): "The Problem" → "the problem", "The Solution" → "the solution", "Principles" → "principles", "The Platform" → "the platform", "Release Pipeline" → "release pipeline", "Founding Decisions" → "founding decisions", etc.
- `<title>` tags: "zarlcorp — Tools that fight back." → "zarlcorp — tools that fight back."
- Taglines: "Tools that fight back." → "tools that fight back.", "Disposable identities" → "disposable identities"
- Body copy sentence starts: "Every tool you use..." → "every tool you use...", "They call it..." → "they call it..."
- Feature list items: "Generate burner emails" → "generate burner emails", "Encrypted local storage" → "encrypted local storage"
- Footer text: "Open source. MIT licensed." → "open source. MIT licensed."
- Bold labels on decisions page: "Org name: zarlcorp" → "org name: zarlcorp", "Entity without a face" → "entity without a face"
- Window titles: already lowercase, no changes needed
- Nav links: already lowercase, no changes needed

### What stays uppercase
- Technical acronyms: MIT, AES-256-GCM, Go, TUI, CLI, DNS, LIFO, HTTP, JSON, HTTPS, FIFO, CI, SIGINT, SIGTERM, RWMutex, HKDF
- Brand names that are conventionally capitalized: Bubble Tea, Charmbracelet, GoReleaser, Homebrew, GitHub, Argon2id, SeaweedFS, Pi-hole, Bitwarden, Uber, Catppuccin Mocha, Lipgloss, JetBrains Mono
- Package/code references in `<code>` tags: `zapp`, `zcache`, etc. — these are already lowercase
- HTML entity references stay as-is

### Pages to update

**zarlcorp.github.io (org site):**
- index.html — title tag, tagline, footer
- manifesto.html — title tag, page title, h2 headings, all body copy
- architecture.html — title tag, page title, h2/h3 headings, all body copy
- decisions.html — title tag, page title, h2/h3 headings, bold labels, all body copy

**zburn (docs/):**
- index.html — title tag, tagline, feature list items, footer
- docs.html — title tag, h3 headings, all body copy

**zvault (docs/):**
- index.html — title tag, tagline, feature list items, status badge, footer

**zshield (docs/):**
- index.html — title tag, tagline, feature list items, status badge, footer

## Target Repo
zarlcorp/zarlcorp.github.io, zarlcorp/zburn, zarlcorp/zvault, zarlcorp/zshield (multi-repo)

## Agent Role
frontend

## Files to Modify
- zarlcorp.github.io: index.html, manifesto.html, architecture.html, decisions.html
- zburn: docs/index.html, docs/docs.html
- zvault: docs/index.html
- zshield: docs/index.html

## Notes
Read every file before editing — don't guess at the current content. Be thorough but careful with acronyms and brand names.
