# 066: manifesto rewrite

## Objective
Rewrite manifesto.html as a concise, punchy call to action. Strip the bloat.

## Context
The current manifesto has 6 sections (the-problem, the-belief, what-we-build, the-standard, decision-log, founding). It reads like an RFC, not a manifesto. "the-standard" is coding practices (belongs in architecture/dev docs). "decision-log" is moving to its own page (spec 067). "what-we-build" repeats the landing page tool cards and architecture page content.

The new manifesto should make someone feel something and want to act. Three sections max.

## Requirements

### Keep and tighten
- **the-problem** — why privacy tools need to exist. Keep the anger. Trim the explanation.
- **the-belief** — what zarlcorp stands for. The four bullet points are strong. Keep them.

### Condense
- **what-we-build** — trim to: we build terminal-first privacy tools, here's the suite (table), everything is open source MIT. Cut the "success" and "proof of concept" paragraphs — that's internal thinking, not a call to action.

### Remove entirely
- **the-standard** — coding practices. This is dev documentation, not manifesto material. Lives in CLAUDE.md and architecture page already.
- **decision-log** — moving to its own page (spec 067)
- **founding** — "the tools don't exist yet" is stale (zburn shipped). Cut it.

### Add a closing punch
End with something that makes people want to use the tools or contribute. Not a "founding" statement. A call to action.

### Apply voice guide
If spec 065 is done, reference `.claude/VOICE.md` for tone. If not, apply these rules:
- Declarative, not explanatory
- Short sentences
- No corporate padding
- Let the tools speak for themselves

### Nav bar
Keep the current org site nav: `zarlcorp | manifesto | decisions | architecture | github`
(Updated to include the new decisions page from spec 067)

## Target Repo
zarlcorp/zarlcorp.github.io

## Agent Role
frontend

## Files to Modify
- `manifesto.html` — rewritten content

## Dependencies
- 065 (voice guide) — soft dependency, nice to have but not blocking
- 067 (decisions page) — the nav bar update should include "decisions" link

## Notes
- The manifesto should fit comfortably on a screen. If you're scrolling past 2 viewports, it's too long.
- Read the current manifesto.html to understand what to cut
- The three remaining sections should each be one floating window
