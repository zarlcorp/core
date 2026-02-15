# 039: Scaffold zshield repository

## Objective
Create the `zarlcorp/zshield` repository with the standard tool structure, a working `main.go` entrypoint, and the minimal Bubble Tea TUI shell — ready for feature development.

## Context
zshield is the third product: DNS-level tracker and ad blocking. It's the most complex tool (daemon mode, DNS resolver, real-time dashboard) and depends on `znet` which doesn't exist yet. Scaffolding establishes the repo structure; DNS and networking features come later when `znet` is built.

Depends on: spec 036 (initial tags). Can run in parallel with specs 037 and 038.
Issue #25.

## Requirements

### Repository structure
```
zarlcorp/zshield/
├── cmd/zshield/
│   └── main.go              # entrypoint — wires zapp, launches TUI or daemon
├── internal/
│   └── tui/
│       └── tui.go           # root Bubble Tea model (placeholder)
├── go.mod                   # module: github.com/zarlcorp/zshield
├── go.sum
├── Makefile                 # build, test, lint, run targets
├── LICENSE                  # MIT
├── README.md                # minimal — what it is, how to install, how to run
└── .github/
    └── workflows/
        └── ci.yml           # references reusable CI or inline
```

### main.go
- Uses `zapp.New()` for lifecycle management.
- Uses `zapp.SignalContext()` for graceful shutdown.
- Checks for subcommands (placeholder: `version`, `start`, `status`, `allow`, `block`).
- Falls back to TUI mode (attach to running daemon) if no subcommand.
- Imports: `zapp`, `zstyle`.
- Does NOT import `znet` — that package doesn't exist yet.

### TUI shell
- Minimal Bubble Tea program: displays zshield name/version, exits on `q`.
- Uses `zstyle` for consistent theming.
- Placeholder only — DNS dashboard, query log, stats come later.

### Makefile targets
Same as zburn/zvault: `build`, `test`, `lint`, `run`, `clean`.

### go.mod dependencies
- `github.com/zarlcorp/core/pkg/zapp`
- `github.com/zarlcorp/core/pkg/zstyle`
- `github.com/zarlcorp/core/pkg/zcache`
- `github.com/charmbracelet/bubbletea`
- `github.com/charmbracelet/lipgloss`

## Target Repo
zarlcorp/zshield

## Agent Role
backend

## Files to Modify
All files are new — this is a new repository.

## Notes
- Lightest scaffold of the three tools — no zcrypto or zfilesystem dependencies until features are built.
- zshield is lowest priority per the manifesto. Scaffold it to establish the repo but don't expect feature work until zburn and zvault are further along.
