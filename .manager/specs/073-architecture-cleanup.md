# 073: Architecture page cleanup

## Objective
The architecture page has 8 windows and reads like a brain dump. Trim it to 3-4 focused sections that answer "how is zarlcorp software structured?" without repeating the manifesto or embedding godoc.

## Context
Current windows: how-we-build, zarlcorp/core, package-roster, package-layering, products, product-details, release-pipeline, agent-model. Half of this content either belongs elsewhere (tool docs, godoc) or is internal process (agent model).

## Requirements

### New structure — 3 windows

**Window 1: "principles"**
Brief technical principles — the HOW, not the WHY (manifesto covers why):
- TUI-first with Bubble Tea
- Single Go binaries, no runtime dependencies
- Local storage, encrypted at rest
- Each tool is its own repo, imports shared packages from core
- 2-3 short paragraphs. No overlap with manifesto.

**Window 2: "the-platform"**
How core and tools fit together:
- Core repo provides shared packages (list them by name with ONE LINE each — no API details)
- Package dependency diagram (keep the ASCII layering art — it's good, but simplify the labels)
- Standard tool structure (the directory tree from the current "products" section)
- How tools consume core (the import example)
- Drop: `.claude/` and `.manager/` from the repo tree (internal agent config)
- Drop: detailed package APIs (struct names, interface methods) — that's godoc

**Window 3: "release-pipeline"**
Keep the current release pipeline section mostly as-is. It's already concise:
- Tag → GoReleaser → GitHub Releases → Homebrew tap
- Reusable CI from zarlcorp/.github

### What to remove entirely
- **how-we-build**: "agent-driven development" paragraphs — internal process. The TUI-first and single binary points fold into "principles"
- **package-roster** detail: the individual package breakdown with struct/interface listings. Replace with a brief table (name | one-line purpose)
- **product-details**: full CLI usage per tool (belongs on each tool's docs page)
- **agent-model**: entirely internal workflow. Nobody visiting the website needs to know about /discuss, /plan, /delegate, or agent personas

### Page header
Change subtitle from "How the platform is built." to something tighter:
"how it fits together." (lowercase per voice guide)

### Voice
Follow `.claude/VOICE.md`:
- Declarative, short sentences
- No corporate hedging
- Lowercase subtitle
- Technical precision, casual delivery

## Target Repo
zarlcorp/zarlcorp.github.io

## Agent Role
frontend

## Files to Modify
- architecture.html

## Notes
Independent of 070/071/072. Can run in parallel with everything.
