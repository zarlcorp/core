# 067: decisions page

## Objective
Create a decisions.html page on the org site. Move the decision log from the manifesto to its own page. Add it to the nav bar.

## Context
The decision log is a historical record of choices made during zarlcorp's founding and ongoing development. It belongs on the public site as a living document — but not in the manifesto, which should be a call to action.

## Requirements

### Create `decisions.html`
- Page header: "Decisions" with subtitle "What we chose and why."
- Same page structure as architecture.html (`.page` wrapper, `.page-header`, floating windows)
- One floating window containing the full decision log content currently in manifesto.html (the "decision-log" section)
- Organized chronologically: "2026-02-15 — Founding Decisions" as the first batch
- Room to grow — future decisions get added as new date-headed sections within the same window, or as new windows per date

### Update nav bar on ALL org site pages
Update the nav in `index.html`, `manifesto.html`, `architecture.html`, and `decisions.html`:
```
zarlcorp | manifesto | decisions | architecture | github
```

### Styling
No new CSS needed — uses existing `.page`, `.page-header`, `.window`, `.doc-content` classes from shared.css.

## Target Repo
zarlcorp/zarlcorp.github.io

## Agent Role
frontend

## Files to Modify
- `decisions.html` — new page
- `index.html` — update nav bar, update landing-nav links
- `manifesto.html` — update nav bar (decision-log content removed in spec 066)
- `architecture.html` — update nav bar

## Dependencies
None (can run in parallel with 066)

## Notes
- Copy the decision log content verbatim from the current manifesto.html "decision-log" window
- The landing page's `.landing-nav` links should be updated to include decisions: "read the manifesto / decisions / explore the architecture"
