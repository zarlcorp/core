# 038: Scaffold zvault repository

## Objective
Create the `zarlcorp/zvault` repository with the standard tool structure, a working `main.go` entrypoint, and the minimal Bubble Tea TUI shell — ready for feature development.

## Context
zvault is the second product: encrypted local secret storage. It builds on the same core packages as zburn, especially `zcrypto` for encryption and `zfilesystem` for storage.

This spec covers scaffolding only. Feature development (secret storage, search, auto-lock, etc.) is separate work.

Depends on: spec 036 (initial tags). Can run in parallel with spec 037 (zburn scaffold).
Issue #24.

## Requirements

### Repository structure
```
zarlcorp/zvault/
├── cmd/zvault/
│   └── main.go              # entrypoint — wires zapp, launches TUI or CLI
├── internal/
│   └── tui/
│       └── tui.go           # root Bubble Tea model (placeholder)
├── go.mod                   # module: github.com/zarlcorp/zvault
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
- Checks for subcommands (placeholder: `version`, `get`, `set`, `search`).
- Falls back to TUI mode if no subcommand.
- Imports: `zapp`, `zstyle`, `zcrypto`, `zfilesystem`.

### TUI shell
- Minimal Bubble Tea program: displays zvault name/version, exits on `q`.
- Uses `zstyle` for consistent theming.
- Placeholder only — vault browsing, secret display, etc. come later.

### Makefile targets
Same as zburn: `build`, `test`, `lint`, `run`, `clean`.

### go.mod dependencies
- `github.com/zarlcorp/core/pkg/zapp`
- `github.com/zarlcorp/core/pkg/zstyle`
- `github.com/zarlcorp/core/pkg/zcrypto`
- `github.com/zarlcorp/core/pkg/zfilesystem`
- `github.com/zarlcorp/core/pkg/zcache`
- `github.com/charmbracelet/bubbletea`
- `github.com/charmbracelet/lipgloss`

## Target Repo
zarlcorp/zvault

## Agent Role
backend

## Files to Modify
All files are new — this is a new repository.

## Notes
- Identical scaffold pattern to zburn (spec 037). The differentiation comes during feature development.
- zvault has more core dependencies than zburn — zcrypto, zfilesystem, and zcache are all needed for its feature set.
