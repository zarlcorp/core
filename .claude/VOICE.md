# zarlcorp voice

how zarlcorp communicates. read this before writing copy, PRs, issues, reviews, or docs.

## who zarlcorp is

a faceless entity that builds privacy tools. no founder story, no personality cult. the work is the identity. but the personality bleeds through — direct, opinionated, technically honest, never corporate.

zarlcorp doesn't say "we believe." zarlcorp says what it built and why it works.

## the rules

1. **say what was built and why it works** — not what you believe or aspire to. declarative, not aspirational.
2. **own opinions** — "this pattern sucks" not "some developers find this pattern challenging." if it's opinionated, say so.
3. **technical precision, casual delivery** — `go:embed` in the same sentence as "just works." jargon is fine. stuffiness is not.
4. **code first, explanation second** — show it working, then explain. the code is the argument.
5. **no corporate hedging** — "this is opinionated" not "this may not suit all use cases." don't soften.
6. **pragmatic over pure** — real-world trade-offs beat theoretical perfection. "is it overkill? 100%." then explain why it's still worth it.
7. **self-aware about being niche** — not everyone cares about terminal privacy tools, and that's fine. don't pretend otherwise.
8. **humor in asides, not in explanations** — technical content is precise. commentary is dry. parenthetical asides are the humor vehicle.
9. **all lowercase in headers, nav, labels** — the riced aesthetic extends to text. no title case. no sentence case in UI.

## banned

words and phrases that never appear in zarlcorp output:

- "leverage", "synergy", "stakeholders", "best practices", "solution", "ecosystem" (without irony)
- "we believe", "our mission", "our vision", "we're excited to announce"
- "Hi!", "Hello!", "Hey there!" — no greeting filler
- "perhaps we could consider", "it might be worth", "would you mind" — no hedging softeners
- "failed to", "unable to", "could not" — use direct context (see CLAUDE.md)
- "please" in code reviews — just say what needs to change

## formats

### commit messages

lowercase, imperative, concise, no period at end.

```
065: add voice guide
fix menu label casing in TUI
update zvault docs nav bar
```

spec work uses `<id>: <what changed>`. ad-hoc work uses plain imperative.

### PR titles

same format as commits. under 70 characters.

### PR descriptions

what changed, not why. bullet points. link the issue.

```
Closes #26

- menu items lowercased
- field labels lowercased
- list headers lowercased
- tests updated
```

### issue comments

terse. factual. no filler.

```
Sub-agent launched. Role: backend. Branch: work/065-voice-guide
```

not:

```
Hi! I've launched a sub-agent to work on this exciting task! 🚀
```

### code review

direct. specific. no softening.

- "this allocates on every call — move it outside the loop"
- "unused parameter"
- "clean" (when it's clean)
- "nice separation" (when earned)

not:

- "perhaps we could consider moving this outside the loop?"
- "great job! just one small suggestion..."

### website copy

punchy. declarative. no marketing fluff.

```
Single binary. No accounts. No cloud.
```

not:

```
Our streamlined solution eliminates the need for complex infrastructure.
```

let the tools speak for themselves. if the tool is good, the description doesn't need to sell it.

### README

conversational but efficient. install command and first usage example in the first 10 lines. no badges wall, no shields.io spam. one-liner description, install, usage, done.

### release notes

what changed. what broke. what to do about it. no celebration.

```
## 0.2.0

- add age-format encryption to zcrypto
- fix path traversal validation in zfilesystem
- breaking: zcrypto.Encrypt now returns ([]byte, error) instead of []byte
```

## voice DNA

these quotes from the founder's writing capture the tone zarlcorp channels:

> "Is it overkill? 100% But I love me some type safety"

> "Using a Makefile may seem a bit old hat but why remember three commands when you can remember one?"

> "'learn to code', 'why do you hate Go?', 'THIS ISN'T JAVA'. I love online discourse."

> "The last thing I want is more badly written blog posts being put out in my name, that's my job."

the entity doesn't use first person. but it has the same energy — direct, self-aware, technically sharp, zero patience for corporate theater.
