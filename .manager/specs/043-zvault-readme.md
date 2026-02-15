# 043: zvault README

## Objective
Create a README.md for the zvault repo that communicates what the tool does, how to install it, and how to use it.

## Context
zvault is encrypted local secret storage. The repo has been scaffolded with a TUI shell but has no README. The core repo README serves as a style reference.

The Manifesto describes zvault as: "Encrypted local storage for secrets, keys, notes. Your data, your machine, your keys."

## Requirements
- Title: `zvault` with a one-line tagline from the Manifesto
- Brief description (2-3 sentences max) explaining what it does and why
- Install section: `go install` and `brew install zarlcorp/tap/zvault` (note: brew formula not yet created)
- Usage section showing both TUI mode (`zvault`) and CLI subcommands (`zvault version`, `zvault get <path>`, `zvault set <path>`, `zvault search <query>`)
- Development section: `make build`, `make test`, `make lint`
- Link to core repo and Manifesto
- MIT license footer
- No badges, no screenshots, no emoji, no verbose prose
- Match the tone of the core README: terse, direct, lowercase comment style

## Target Repo
zarlcorp/zvault

## Agent Role
backend

## Files to Create
- README.md

## Notes
zvault is pre-release — most CLI subcommands are not yet implemented. The README should describe the intended interface without claiming features work. The `version` subcommand and TUI mode do work.
