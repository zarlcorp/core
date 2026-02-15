# 076: Install command copy button

## Objective
Add a small clipboard icon to every `.install-cmd` element that copies the command text to the clipboard when clicked.

## Context
Tool landing pages show install commands (brew/curl). Users currently have to manually select and copy. A copy button makes this frictionless. The install infrastructure lives in shared.css and install.js on the org site — all tool pages inherit it.

Only zburn has an install section today. zvault and zshield are "coming soon" with no install commands yet. Building this into the shared infrastructure means they'll get the button automatically when they add install sections.

## Requirements

### JavaScript (install.js)
Add a second self-invoking function (or extend the existing one) that:
1. Finds all `.install-cmd` elements
2. Appends a `<button class="copy-btn">` with a small clipboard SVG icon to each one
3. On click, copies the text content of the `.install-cmd` (excluding the `$ ` prompt) to the clipboard via `navigator.clipboard.writeText()`
4. The button should be minimal — no text label, just the icon

The SVG should be a simple clipboard outline, ~16x16, using `currentColor` so it inherits the text color from CSS.

### CSS (shared.css)
Add styles for `.copy-btn` inside `.install-cmd`:
- Position: absolute right inside the install-cmd block (the install-cmd needs `position: relative`)
- Color: `var(--overlay1)` (same as surrounding text)
- Hover: `var(--accent)` or `var(--text)` — something visible but not loud
- Background: transparent, no border
- Cursor: pointer
- Size: small, doesn't crowd the command text
- Vertically centered within the install-cmd

### No HTML changes needed
The button is injected by JavaScript. No tool page HTML needs to change.

## Target Repo
zarlcorp/zarlcorp.github.io

## Agent Role
frontend

## Files to Modify
- shared.css (add .copy-btn styles)
- install.js (add copy button injection + click handler)

## Notes
Icon only — no "copied!" feedback text or checkmark animation. Keep it simple. Test that the prompt `$ ` is excluded from the copied text.
