# zarlcorp Coding Standards

Cross-language principles and coding style for all zarlcorp projects. See MANIFESTO.md for the big picture — why we exist and what we're building.

## Language-Specific Guides

- [CLAUDE_GO.md](./CLAUDE_GO.md) - Go patterns, repository/service/transport layers, app lifecycle

## Architecture

### TUI-First
Every tool is a single Go binary with a Bubble Tea TUI. No browser, no JavaScript. The terminal is the primary interface. Web UIs are added only when the TUI isn't sufficient.

### Single Binary Deployment
One binary serves everything. Download, run, done. No Docker, no databases, no environment variables. Local storage is files on disk, encrypted at rest.

### Layer Separation
Each layer owns its types - no shared "domain" package:
- Repository layer: database/storage types
- Service layer: business logic types (no tags)
- Transport layer: request/response types (with wire format tags)

Map between layers at boundaries with explicit conversion functions.

### Abstraction Process
1. Build concrete implementation first
2. "Poke the problem with reality"
3. Delete prototype
4. Rebuild with understanding
5. Extract interfaces from actual usage patterns

### Project Context
- **Personal/Simple**: Keep concrete, move fast
- **Complex domains**: Structure after understanding emerges
- **Libraries**: Design for unknown use cases

## Core Principles

**Build Philosophy**: Start concrete → understand problem → delete prototype → rebuild with understanding → extract abstractions only when needed

**Error Philosophy**: Errors tell a story - build narrative without stuttering, wrap at every failure point, log once at boundaries

**Interface Philosophy**: Small interfaces, consumer-side definition, emergent not design-first

**Type Philosophy**: Semantic types for complex domains, scope-based naming, avoid primitive obsession

**Code Quality Philosophy**: No duplication in branches, extract common operations, prefer early returns over if/else chains

## Error Handling

### Never Use These Prefixes
- "failed to", "unable to", "could not", "error"
- Use direct context instead: `"open file: " + err`

### Logging Strategy
- Never log every error occurrence
- Log once at application boundaries with full context
- Let error chain build the story through the stack

## Naming Conventions

### Scope-Based Naming
Smaller scope = shorter names:
- Loop variables: `i`, `j`, `k`
- Short-lived: `u`, `r`, `w`
- Larger scope: `requestID`, `wordsPerMinute`

### Abbreviations
- Avoid unless universally understood (`URL`, `ID`, `HTTP`)
- When in doubt, spell it out

## Code Quality

### Early Returns Over If/Else
```go
// BAD
if condition {
    doA()
} else {
    doB()
}

// GOOD
if condition {
    doA()
    return
}
doB()
```

### Extract Common Operations
```go
// BAD - duplicated in both branches
if condition {
    doSpecific()
    commonOp()
    return
}
doOther()
commonOp()  // duplicated!

// GOOD
if condition {
    doSpecific()
} else {
    doOther()
}
commonOp()
```

**Rule**: ALWAYS review code for duplication before submitting.

## Testing

### Philosophy
- Test exposed API, not internals
- Prefer table-driven tests
- Use contract tests for implementation conformity
- Avoid mocking - prefer in-memory implementations

### Test Hierarchy
1. Real implementations (best)
2. In-memory implementations with seed data
3. Mocks (avoid if possible)

## Comments & Documentation

### Style
- Lowercase, terse, minimal punctuation
- Focus on "why" not "what"

### When to Comment
**Exported APIs:** Brief description for docs

**Inline:** Explain WHY decisions were made, not WHAT the code does

### Strategy
- **Library/SDK**: Comprehensive with examples
- **Internal code**: Minimal, focus on naming clarity

## Performance vs Readability

### Library Code
- Must be fast and bulletproof
- Prevent panics at all costs

### Application Code
- Prioritize readability and simplicity
- Optimize only when measured bottlenecks exist

## Dependencies

Every dependency has a cost — build time, supply chain risk, upgrade burden. Prefer the standard library. When a dependency is justified, pin it and understand what it does. The Charmbracelet ecosystem is the one major dependency family we embrace.

## Anti-Patterns

Quick reference - applies across languages:

- Primitive obsession (use semantic types)
- Shared domain package (each layer owns types)
- Large interfaces for normal operations
- Verbose error prefixes
- Logging every error occurrence
- Duplicated code in branches
- if/else chains instead of early returns
- Design-first interfaces
- Premature abstraction
- Fire-and-forget async operations
- Mocking when fakes work
