# 044: zshield README

## Objective
Create a README.md for the zshield repo that communicates what the tool does, how to install it, and how to use it.

## Context
zshield is DNS-level tracker and ad blocking. The repo has been scaffolded with a TUI shell but has no README. The core repo README serves as a style reference.

The Manifesto describes zshield as: "DNS-level tracker and ad blocking. Single binary Pi-hole. See what's tracking you, then kill it."

## Requirements
- Title: `zshield` with a one-line tagline from the Manifesto
- Brief description (2-3 sentences max) explaining what it does and why
- Install section: `go install` and `brew install zarlcorp/tap/zshield` (note: brew formula not yet created)
- Usage section showing both TUI mode (`zshield`) and CLI subcommands (`zshield version`, `zshield start`, `zshield status`, `zshield allow <domain>`, `zshield block <domain>`)
- Development section: `make build`, `make test`, `make lint`
- Link to core repo and Manifesto
- MIT license footer
- No badges, no screenshots, no emoji, no verbose prose
- Match the tone of the core README: terse, direct, lowercase comment style

## Target Repo
zarlcorp/zshield

## Agent Role
backend

## Files to Create
- README.md

## Notes
zshield is pre-release — most CLI subcommands are not yet implemented. The README should describe the intended interface without claiming features work. The `version` subcommand and TUI mode do work.
