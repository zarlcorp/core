# 068: tool page nav consistency

## Objective
Make all tool page nav bars consistent with the org site bar pattern. Fix capitalization and HTML structure.

## Context
The org site bar (`zarlcorp | manifesto | decisions | architecture | github`) sets the standard. Tool pages should follow the same pattern with tool-appropriate links:

```
zarlcorp | documentation | github
```

Current state:
- **zburn**: `zarlcorp | docs | github` — close but "docs" should be "documentation"
- **zvault**: `zarlcorp | Docs | Org | GitHub` — wrong capitalization, extra "Org" link, uses `<nav class="nav">` wrapper inside the bar
- **zshield**: `zarlcorp | docs | github` — same as zburn

All need to use the same HTML structure as the org site.

## Requirements

### Nav bar template for all tool pages
Both `index.html` and `docs.html` in each tool:
```html
<nav class="bar">
  <div class="container">
    <a href="https://zarlcorp.github.io" class="nav-brand">zarlcorp</a>
    <a href="docs.html">documentation</a>
    <a href="https://github.com/zarlcorp/<tool>">github</a>
  </div>
</nav>
```

On the docs page itself, the "documentation" link changes to the tool name linking back to index:
```html
<nav class="bar">
  <div class="container">
    <a href="https://zarlcorp.github.io" class="nav-brand">zarlcorp</a>
    <a href="index.html"><tool></a>
    <a href="https://github.com/zarlcorp/<tool>">github</a>
  </div>
</nav>
```

### Fix zvault HTML structure
zvault currently uses `<div class="bar"><div class="container"><nav class="nav">` — fix to match the standard `<nav class="bar"><div class="container">` pattern.

### All lowercase
Every link label: `zarlcorp`, `documentation`, `github`, `<tool>` — all lowercase.

## Target Repos
- zarlcorp/zburn — `docs/index.html`, `docs/docs.html`
- zarlcorp/zvault — `docs/index.html`, `docs/docs.html`
- zarlcorp/zshield — `docs/index.html`, `docs/docs.html`

## Agent Role
frontend

## Files to Modify
- zburn: `docs/index.html`, `docs/docs.html`
- zvault: `docs/index.html`, `docs/docs.html`
- zshield: `docs/index.html`, `docs/docs.html`

## Dependencies
None

## Notes
- This is a surgical fix — only change the nav bars, don't touch any other content
- zvault needs the most work (HTML structure + capitalization + extra link)
- zburn and zshield just need "docs" → "documentation"
