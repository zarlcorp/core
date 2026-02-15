# 065: zarlcorp voice guide

## Objective
Create `.claude/VOICE.md` in core — a style document that defines how zarlcorp communicates across all written output.

## Context
zarlcorp speaks as a faceless entity, but the personality bleeds through from the founder's writing on zarl.dev. The voice is direct, opinionated, technically precise, and never corporate. This guide codifies that voice so agents and contributors write in a consistent tone.

## Requirements

### Scope
The voice guide covers:
- Website copy (landing pages, manifesto, docs)
- README files
- PR titles and descriptions
- Commit messages
- GitHub issue comments and descriptions
- Code review feedback
- Release notes

The voice guide does NOT cover:
- CLI output text
- TUI labels and prompts
- Code comments (covered by CLAUDE.md)

### Voice DNA — extracted from zarl.dev

The founder's writing has these patterns:

**Tone**: Direct, opinionated, self-aware, technically honest. Never hedges. Never corporate. Owns mistakes openly ("in hindsight I'm an idiot for not doing this from the start").

**Structure**: Short punchy statements mixed with technical depth. Fragments for emphasis. Rhetorical questions as transitions. Parenthetical asides for dry humor.

**What it NEVER does**: Corporate speak ("leverage", "synergy", "stakeholders"). Buzzwords without mockery. Academic passive voice. Unnecessary hedging. Neutral conclusions.

**Humor**: Self-deprecating, dry, sardonic. Acknowledges over-engineering. Emoji sparingly and only to punctuate a joke, never decorative.

### Translation to org voice

zarlcorp is not Bruno — it's the entity. The DNA carries over but the first person singular doesn't. Key rules:

1. **Say what was built and why it works** — never "we believe" or "our mission"
2. **Own opinions** — "this pattern sucks" not "some developers find this pattern challenging"
3. **Technical precision, casual delivery** — `go:embed` in the same sentence as "just works"
4. **Code first, explanation second** — show it working, then explain
5. **No corporate hedging** — "this is opinionated" not "this may not suit all use cases"
6. **Pragmatic over pure** — real-world trade-offs beat theoretical perfection
7. **Self-aware about being niche** — not everyone cares about this, and that's fine
8. **Humor in asides, not in explanations** — technical content is precise, commentary is dry
9. **All lowercase in headers, nav, UI labels** — the riced aesthetic extends to text

### Commit messages
- Lowercase, imperative, concise
- No period at end
- Format: `<id>: <what changed>` for spec work, plain imperative for ad-hoc
- Examples: `065: add voice guide`, `fix menu label casing in TUI`, `update zvault docs nav bar`

### PR titles
- Same format as commit messages
- Short — under 70 characters

### PR descriptions
- Lead with what changed, not why (the "why" is in the spec/issue)
- Bullet points over paragraphs
- Link the issue with "Closes #N"

### Issue comments
- Terse, factual, no filler
- "Sub-agent launched. Role: backend. Branch: work/065-voice-guide" — not "Hi! I've launched a sub-agent to work on this."

### Code review feedback
- Direct, specific, no softening
- "this allocates on every call — move it outside the loop" not "perhaps we could consider moving this outside the loop?"
- Praise when earned, briefly: "clean" or "nice separation"

### Website copy
- Punchy, declarative, no marketing fluff
- "Single binary. No accounts. No cloud." not "Our streamlined solution eliminates the need for complex infrastructure."
- Let the tools speak for themselves

### README tone
- Conversational but efficient
- Show install + first command in the first 10 lines
- No badges wall, no shields.io spam
- One-liner description, install, usage, done

## Target Repo
zarlcorp/core

## Agent Role
backend

## Files to Create
- `.claude/VOICE.md` — the voice guide

## Dependencies
None

## Notes
- Read the full voice analysis from zarl.dev posts (provided below) when writing
- Keep it concise — under 150 lines. This is a reference, not an essay.
- The guide should feel like it was written in the voice it describes
