# 074: Fixed landing page position

## Objective
Replace `justify-content: center` on `.landing` with a fixed `padding-top` so the logo/title always appears at the same vertical position on every page, regardless of content volume below it.

## Context
Currently `.landing` uses `justify-content: center` which vertically centers content in the viewport. Pages with more content (org homepage with tool grid) push the logo higher, while pages with less content (zvault, zshield) show the logo lower. The user wants every page to look the same — logo at the same spot, content flowing below.

## Requirements

### Change in shared.css
Find the `.landing` rule block:
```css
.landing {
  min-height: calc(100vh - 40px);
  display: flex;
  flex-direction: column;
  justify-content: center;
}
```

Replace `justify-content: center` with a fixed padding-top. Use roughly `8vh` or `10vh` — enough to give breathing room below the nav bar without pushing content too far down. The exact value should match the visual feel of the org homepage (logo roughly in the upper third of the viewport).

The result:
```css
.landing {
  min-height: calc(100vh - 40px);
  display: flex;
  flex-direction: column;
  padding-top: 8vh;
}
```

Adjust the `padding-top` value if needed after checking visually — the goal is the logo sits at the same vertical position as it currently does on the org homepage.

### No other changes needed
All pages (org, zburn, zvault, zshield) inherit `.landing` from shared.css. This single change fixes all of them.

## Target Repo
zarlcorp/zarlcorp.github.io

## Agent Role
frontend

## Files to Modify
- shared.css

## Notes
One-line change. No dependencies.
