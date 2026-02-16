# 094: GitHub Pages for core shared library

## Objective
Build a GitHub Pages site for `zarlcorp/core` matching the shared Catppuccin Mocha design system used by zburn, zvault, and zshield.

## Context
Every tool repo has a GitHub Pages site deployed from `docs/`. The core shared library has no site yet. It should follow the same pattern — landing page + docs page — but adapted for a Go package library instead of a CLI tool.

The shared design system lives at `zarlcorp.github.io/shared.css`. Each site imports it and adds a small `style.css` for overrides (accent color, spacing).

Reference sites:
- `zarlcorp/zburn/docs/` — peach accent
- `zarlcorp/zvault/docs/` — mauve accent
- `zarlcorp/zshield/docs/` — teal accent
- `zarlcorp/zarlcorp.github.io/` — org landing, lavender accent

## Requirements

### Accent color
`--sky` (#89dceb)

### Files to create

**`docs/style.css`**
- Set `--accent: var(--sky)`
- Minimal overrides matching the tool site patterns (landing-links, docs-page spacing)

**`docs/index.html`** — landing page
- Import `shared.css` from `https://zarlcorp.github.io/shared.css`
- Import local `style.css`
- Nav bar: `zarlcorp` (brand link to zarlcorp.github.io) | `documentation` (link to docs.html) | `github` (link to github.com/zarlcorp/core)
- Hero: logo text "core", tagline "shared go packages for zarlcorp privacy tools."
- Features window listing all 7 packages:
  - zapp — application lifecycle, resource cleanup, signal handling
  - zcache — generic caching with pluggable backends
  - zcrypto — encryption primitives, key derivation, secure erase
  - zfilesystem — filesystem abstraction with OS and in-memory implementations
  - zoptions — generic functional options pattern
  - zstyle — catppuccin mocha palette, lipgloss presets, standard keybindings
  - zsync — thread-safe map, set, and blocking queue
- Install section with `go get` commands (show 3-4 key packages):
  ```
  go get github.com/zarlcorp/core/pkg/zapp
  go get github.com/zarlcorp/core/pkg/zstyle
  go get github.com/zarlcorp/core/pkg/zcache
  go get github.com/zarlcorp/core/pkg/zsync
  ```
- Landing links row: docs / zarlcorp / github
- Footer: open source. MIT licensed.

**`docs/docs.html`** — package reference
- Same nav bar as index.html but with `core` linking back to index.html instead of `documentation`
- One window per package (7 windows total), each with:
  - window-title: package name (e.g., "zapp")
  - doc-content with: brief description, key exported types/functions (as a list or table), one short usage example in a `<pre><code>` block
- Package details (high level — key exports and one example each):
  - **zapp**: `New()`, `SignalContext()`, `app.Track()`, `app.Close()` — lifecycle management
  - **zcache**: `Cache[K,V]` interface, `MemoryCache`, `FileCache`, `RedisCache` — pluggable caching
  - **zcrypto**: `Encrypt()`, `Decrypt()`, `DeriveKey()`, `Erase()` — AES-256-GCM encryption
  - **zfilesystem**: `ReadWriteFileFS` interface, `NewMemFS()`, `NewOSFileSystem()` — filesystem abstraction
  - **zoptions**: `Option[T]` type, `Apply()` — generic functional options
  - **zstyle**: `Logo`, `Colors` (Catppuccin palette), `DefaultKeyMap` — TUI visual identity
  - **zsync**: `ZMap[K,V]`, `ZSet[T]`, `Queue[T]` — thread-safe collections
- Footer: back to core link, MIT licensed

### Design constraints
- All copy lowercase (matching zarlcorp voice guide)
- Use the existing shared.css classes: `.bar`, `.landing`, `.hero`, `.logo`, `.tagline`, `.window`, `.window-title`, `.window-content`, `.doc-content`, `.feature-list`, `.install-cmd`, `.footer`, `.container`
- No JavaScript required (pure static HTML/CSS)
- Title tag format: "core — shared go packages."

### GitHub Pages deployment
GitHub Pages should be configured to deploy from the `docs/` folder on main branch. This can be set via repo settings or the `gh` CLI.

## Target Repo
zarlcorp/core

## Agent Role
frontend

## Files to Modify
- docs/index.html (create)
- docs/style.css (create)
- docs/docs.html (create)

## Notes
- Do NOT modify the org landing page (zarlcorp.github.io) — keep it tools-only
- Do NOT use the frontend-design skill — write clean HTML/CSS by hand matching the existing pattern
- Reference existing tool sites for exact structure and class usage
- The install section should NOT have OS-aware detection (no install.js) — just show `go get` commands
