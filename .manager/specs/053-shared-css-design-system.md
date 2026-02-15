# 053: Shared CSS design system

## Objective
Create a reusable CSS design system for all zarlcorp web pages — org site and per-tool pages.

## Context
The current org site has inline styles in a single `style.css`. As we add tool-specific pages on each repo, we need a shared CSS file that all pages import for consistent look and feel. The design system lives in the org repo and is referenced by tool pages.

## Requirements

### Shared CSS file (`shared.css`)
Extract and extend the current `style.css` into a design system:

**CSS Variables (from zstyle palette)**:
- `--bg: #0d1117` — page background
- `--bg-card: #161b22` — card/block background
- `--cyan: #00E5FF` — primary accent (headings, links, active)
- `--orange: #FF6E40` — secondary accent (taglines, CTAs)
- `--bright: #ECEFF1` — primary text
- `--muted: #78909C` — secondary text
- `--subtle: #37474F` — borders, separators
- `--success: #69F0AE` — install commands, OK states
- `--error: #FF5252` — error states
- `--warning: #FFD740` — warning states

**Typography**:
- Monospace stack: `"SF Mono", "Fira Code", "Cascadia Code", Menlo, Consolas, monospace`
- Sans-serif stack: `-apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif`
- Base font size: 18px desktop, 16px mobile

**Reusable classes**:
- `.container` — max-width centered layout
- `.hero` — page hero section
- `.logo` — monospace cyan heading
- `.tagline` — orange subtitle
- `.tool-card` — card with subtle border, hover to cyan
- `.install-cmd` — terminal-style command block with non-selectable prompt
- `.footer` — muted centered footer
- `.btn` — subtle bordered button/link
- `.section-heading` — muted section label
- `.doc-content` — documentation prose styling (headings, code blocks, lists)

**OS-aware install component**:
- Small inline JS snippet (~5 lines) that detects macOS vs Linux via `navigator.platform`
- Shows `brew install zarlcorp/tap/<tool>` for macOS
- Shows `curl -sL https://github.com/zarlcorp/<tool>/releases/latest/download/<tool>_linux_amd64 -o <tool> && chmod +x <tool>` for Linux
- Copyable command (click-to-copy or just selectable text with non-selectable prompt)

**Responsive**:
- Single column on mobile, grid on desktop
- Tool landing pages must fit in one viewport (no scroll)

### Hosting
The shared CSS lives at `zarlcorp.github.io/shared.css` and is referenced by tool pages via:
```html
<link rel="stylesheet" href="https://zarlcorp.github.io/shared.css">
```

## Target Repo
zarlcorp/zarlcorp.github.io

## Agent Role
frontend

## Files to Create
- shared.css — the design system
- install.js — small inline OS detection snippet (can be inlined in pages instead)

## Files to Modify
- style.css — refactor to import/use shared.css patterns
- index.html — update to use shared class names if they change

## Notes
- The current style.css is a good starting point — extend it, don't start from scratch
- The JS snippet should be tiny — just `navigator.platform` check and DOM swap
- Tool pages will reference this via absolute URL from their `docs/` folders
