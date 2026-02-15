# 037: Scaffold zburn repository

## Objective
Create the `zarlcorp/zburn` repository with the standard tool structure, a working `main.go` entrypoint, and the minimal Bubble Tea TUI shell — ready for feature development.

## Context
zburn is the flagship product: disposable identity generation. It's the first tool to validate the full stack (core packages → tool → release → Homebrew). The manifesto defines its capabilities and CLI interface.

This spec covers scaffolding only — the repo structure, build system, and a "hello world" TUI. Feature development (email generation, identity storage, etc.) is separate work.

Depends on: spec 036 (initial tags, so go.mod can reference tagged core packages).
Issue #23.

## Requirements

### Repository structure
```
zarlcorp/zburn/
├── cmd/zburn/
│   └── main.go              # entrypoint — wires zapp, launches TUI or CLI
├── internal/
│   └── tui/
│       └── tui.go           # root Bubble Tea model (placeholder)
├── go.mod                   # module: github.com/zarlcorp/zburn
├── go.sum
├── Makefile                 # build, test, lint, run targets
├── LICENSE                  # MIT
├── README.md                # minimal — what it is, how to install, how to run
└── .github/
    └── workflows/
        └── ci.yml           # references reusable CI from zarlcorp/.github (or inline if 040 not done yet)
```

### main.go
- Uses `zapp.New()` for lifecycle management.
- Uses `zapp.SignalContext()` for graceful shutdown.
- Checks for subcommands (placeholder: just `version`).
- Falls back to TUI mode if no subcommand.
- Imports: `zapp`, `zstyle` (for consistent theming).

### TUI shell
- Minimal Bubble Tea program that displays the zburn name/version and exits on `q`.
- Uses `zstyle` colors and styles for the header.
- This is the skeleton — feature screens (email gen, identity list, etc.) come later.

### Makefile targets
- `build` — `go build -o bin/zburn ./cmd/zburn`
- `test` — `go test -race ./...`
- `lint` — `golangci-lint run`
- `run` — `go run ./cmd/zburn`
- `clean` — remove `bin/`

### go.mod dependencies
- `github.com/zarlcorp/core/pkg/zapp` (tagged version from spec 036)
- `github.com/zarlcorp/core/pkg/zstyle` (tagged version)
- `github.com/charmbracelet/bubbletea`
- `github.com/charmbracelet/lipgloss`

### CI
- If spec 040 (reusable CI) is done, reference the shared workflow.
- Otherwise, inline a simple CI: build + test + lint on PR.

## Target Repo
zarlcorp/zburn

## Agent Role
backend

## Files to Modify
All files are new — this is a new repository.

## Notes
- The repo is created by the manager via `gh repo create zarlcorp/zburn --public --license MIT`.
- Keep the scaffold minimal. No feature code. The goal is a compiling, lintable, testable skeleton.
- Follow the coding standards in CLAUDE.md — the agent should create a `.claude/CLAUDE.md` in the tool repo that references core standards.
