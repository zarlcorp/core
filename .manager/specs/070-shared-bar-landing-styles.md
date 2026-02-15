# 070: Move bar and landing styles to shared.css

## Objective
The bar nav link styles and landing page base layout currently live only in the org site's `style.css`. Move them to `shared.css` so all pages (org + tool pages) get the same nav bar and landing layout.

## Context
After spec 068 made the nav bar HTML consistent, the pages still look different because the CSS styling for `.bar a`, `.bar .nav-brand`, and `.landing` layout are defined per-page instead of in the shared design system. Tool pages have unstyled nav links (no right-alignment, wrong colors) and different vertical positioning.

## Requirements

### Move to shared.css
Add these rule blocks to `shared.css` (they currently exist in org `style.css` lines 6-29):

```css
.bar a {
  color: var(--overlay1);
  font-weight: 400;
  letter-spacing: 0.02em;
}

.bar a:hover {
  color: var(--accent);
}

.bar .nav-brand {
  color: var(--accent);
  font-weight: 700;
  font-size: 0.85rem;
  margin-right: auto;
  letter-spacing: 0.04em;
}
```

The `margin-right: auto` on `.nav-brand` is critical — it pushes documentation/github links to the right.

Note: `.bar .container` is ALREADY in shared.css (lines 115-120). Do NOT duplicate it.

### Add base landing layout to shared.css
Add a standardized `.landing` base that all pages use:

```css
.landing {
  min-height: calc(100vh - 40px);
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.landing .hero {
  padding: 2rem 0 1.5rem;
}

.landing .install {
  padding: 1.5rem 0;
  text-align: center;
}

.landing .footer {
  padding: 1.5rem 0 2rem;
}
```

### Clean up org style.css
Remove the now-duplicated rules from the org site's `style.css`:
- Remove `.bar .container` (lines 6-11) — already in shared.css
- Remove `.bar a` (lines 13-17)
- Remove `.bar a:hover` (lines 19-21)
- Remove `.bar .nav-brand` (lines 23-29)
- Remove `.landing` (lines 33-37) — now in shared.css
- Remove `.landing .hero` (lines 39-41) — now in shared.css
- Remove `.landing .install` (lines 47-50) — now in shared.css
- Remove `.landing .footer` (lines 67-69) — now in shared.css

Keep org-specific styles that don't belong in shared:
- `.landing .tools` padding
- `.landing-nav` styles (org-specific bottom nav)
- `.page`, `.page-header`, `.page-title`, `.page-subtitle` (doc page styles)

### Verification
After changes, the org landing page and subpages (manifesto, architecture, decisions) must render identically to before. Open them in a browser or diff the computed styles.

## Target Repo
zarlcorp/zarlcorp.github.io

## Agent Role
frontend

## Files to Modify
- shared.css
- style.css

## Notes
This must merge before 071 (tool page alignment) since tool pages load shared.css from this repo's GitHub Pages URL.
