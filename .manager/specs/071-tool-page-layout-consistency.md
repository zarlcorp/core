# 071: Tool page layout consistency

## Objective
Remove per-tool CSS overrides that cause layout drift. After 070 moves the base landing/bar styles to shared.css, tool pages should ONLY override the accent color and tool-specific elements. The logo, tagline, and content should sit in the same position across all pages.

## Context
Currently each tool page defines its own `.landing`, `.landing .hero`, `.landing .tagline` rules with slightly different values:
- zburn: `min-height: 100vh; max-height: 100vh; overflow: hidden` + hero `2rem 0 1.25rem`
- zvault: same as zburn + tagline `1.05rem`
- zshield: same + hero `2.5rem 0 1.25rem` + tagline `max-width: 540px`

After 070, shared.css provides the base `.landing` and `.landing .hero` styles. Tool pages must remove their conflicting overrides so the shared base takes effect.

## Requirements

### zburn/docs/style.css
Remove:
- `.landing` block (lines 9-16) — the `flex-direction: column; justify-content: center; min-height: 100vh; max-height: 100vh; overflow: hidden` conflicts with shared base
- `.landing .hero` block (lines 18-20) — different padding than shared base
- `.landing .tagline` block (lines 22-24) — sets font-size to 1rem which is already the shared default
- `.landing .install` block (lines 34-37) — now in shared.css

Keep:
- `:root { --accent: var(--peach); }` — tool accent
- `.landing .features` padding
- `.landing-links` styles (tool-specific bottom links)
- `.landing .footer` only if the padding differs intentionally (check if shared base is sufficient)
- `.docs-page` styles

### zvault/docs/style.css
Remove:
- `.landing` block (lines 9-16) — same conflict as zburn
- `.landing .hero` block (lines 18-20)
- `.landing .tagline` block (lines 22-24) — 1.05rem is inconsistent, should use shared default

Keep:
- `:root { --accent: var(--mauve); }` — tool accent
- `.status-badge` styles — tool-specific
- `.features` styles — tool-specific

### zshield/docs/style.css
Remove:
- `.landing` block (lines 9-16) — same conflict
- `.landing .hero` block (lines 18-20) — `2.5rem 0 1.25rem` is different from shared `2rem 0 1.5rem`
- `.landing .tagline` block (lines 22-24) — `max-width: 540px` creates different horizontal layout

Keep:
- `:root { --accent: var(--teal); }` — tool accent
- `.status-badge` styles — tool-specific
- `.landing .features` padding — tool-specific

### HTML consistency check
Verify all tool landing pages use the same wrapper structure:
```html
<section class="landing">  <!-- not <div> -->
  <div class="container">
    <header class="hero">
      <div class="logo">toolname</div>
      <div class="tagline">...</div>
    </header>
    ...
  </div>
</section>
```

zburn currently uses `<div class="landing">` — change to `<section class="landing">` for consistency.

### Verification
- All three tool pages render with the logo and content in the same vertical/horizontal position
- Nav bar has zarlcorp brand on left, links on right (inheriting from shared.css)
- Only the accent color and tool-specific content differs between pages

## Target Repo
zarlcorp/zburn, zarlcorp/zvault, zarlcorp/zshield (multi-repo)

## Agent Role
frontend

## Files to Modify
- zburn: docs/style.css, docs/index.html
- zvault: docs/style.css
- zshield: docs/style.css

## Notes
Depends on 070 being merged first — tool pages load shared.css from `https://zarlcorp.github.io/shared.css`. The shared styles must be live before removing the local overrides.
