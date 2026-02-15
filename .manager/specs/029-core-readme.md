# 029: Core repo README

## Objective
Add a README.md to zarlcorp/core that introduces the project, lists the packages, shows usage, and links to the MANIFESTO.

## Context
The core repo has 7 packages (zsync, zcache, zoptions, zfilesystem, zapp, zstyle, zcrypto stub) and a MANIFESTO.md, but no README. GitHub visitors see nothing. The README should be concise — the MANIFESTO has the full philosophy and roadmap.

## Requirements

### Header
```
# zarlcorp/core

Shared Go packages for zarlcorp privacy tools.
```

### Package table
List all packages with a one-line description and status (ready / stub):

| Package | Description | Status |
|---------|-------------|--------|
| zapp | Application lifecycle toolkit | ready |
| zcache | Generic caching with multiple backends | ready |
| zcrypto | Encryption primitives | stub |
| zfilesystem | Filesystem abstraction | ready |
| zoptions | Generic functional options | ready |
| zstyle | TUI visual identity — colors, styles, keybindings | ready |
| zsync | Thread-safe data structures | ready |

### Install
Show how to import individual packages:
```go
go get github.com/zarlcorp/core/pkg/zapp
go get github.com/zarlcorp/core/pkg/zstyle
```

### Quick example
Short code block showing zapp + zstyle together — something like:
```go
package main

import (
    "context"
    "fmt"

    "github.com/zarlcorp/core/pkg/zapp"
    "github.com/zarlcorp/core/pkg/zstyle"
)

func main() {
    app := zapp.New(zapp.WithName("example"))
    ctx, cancel := zapp.SignalContext(context.Background())
    defer cancel()

    fmt.Println(zstyle.Title.Render("hello from zarlcorp"))

    <-ctx.Done()
    app.Close()
}
```

### Links
- Link to MANIFESTO.md for philosophy and roadmap
- Link to LICENSE (MIT)

### Footer
```
MIT License
```

### What the README does NOT include
- No badges (no CI configured on core yet)
- No contributing guide (premature)
- No changelog (use git tags)
- No lengthy philosophy — that's in the MANIFESTO

## Acceptance Criteria
1. README.md exists at repo root
2. All 7 packages listed with accurate status
3. Import examples use correct module paths
4. Code example compiles conceptually (uses real API)
5. Links to MANIFESTO.md and LICENSE

## Target Repo
zarlcorp/core

## Agent Role
backend

## Files to Create
- `README.md`

## Notes
- Keep it under 80 lines. The MANIFESTO is the detailed document.
- Use the actual API from zapp (New, WithName, SignalContext, Close) and zstyle (Title.Render) — these are real and merged.
