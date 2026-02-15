# 072: Manifesto — two sections only

## Objective
The manifesto is still overwritten. Strip it down to exactly two windows: the-problem and the-solution. No belief section, no tool table, no closing window.

## Context
The 066 rewrite trimmed the manifesto from 6 sections to 4, but it's still too much. A manifesto should hit hard and fast — here's what's broken, here's what we do about it.

## Requirements

### Structure
Two windows only, plus the page header and footer:

```
page-header: "manifesto" / "software should serve its user and no one else."

window: the-problem
window: the-solution

footer
```

### the-problem
Keep the core message from the current version but make it tighter:
- Services capture your data, your attention, your dependency
- They call it "improving your experience" — they mean "building your profile"
- The technology to build private software has existed for decades. The incentives point the other way.
- 3-4 paragraphs max. Short sentences. Declarative.

### the-solution
Collapse the-belief, what-we-build, and the closing into one section. This is the answer to the problem:
- What zarlcorp builds: terminal-first privacy tools, single binaries, no accounts, no cloud
- The principles: open source, local-first, single binary, MIT licensed
- The tools (brief — not a full table, just mention them)
- End with: code is open source, read it, verify it, fork it
- 4-5 paragraphs max. Same voice — direct, no hedging.

### Voice
Follow `.claude/VOICE.md`:
- Declarative, not aspirational
- No corporate hedging
- Short sentences
- Lowercase page subtitle (already is)
- All banned words apply

### What to remove
- the-belief window (fold into the-solution)
- what-we-build window (fold into the-solution)
- ~ closing window (fold into the-solution)
- The full tool table — too much detail for a manifesto. Mention the tools by name, don't describe each one.

## Target Repo
zarlcorp/zarlcorp.github.io

## Agent Role
frontend

## Files to Modify
- manifesto.html

## Notes
Independent of 070/071 (CSS changes). Can run in parallel with 070.
