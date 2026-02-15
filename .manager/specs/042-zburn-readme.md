# 042: zburn README

## Objective
Create a README.md for the zburn repo that communicates what the tool does, how to install it, and how to use it.

## Context
zburn is the flagship zarlcorp tool — a disposable identity generator. The repo has been scaffolded with a TUI shell but has no README. The core repo README serves as a style reference: concise, no fluff, shows install + quick example + links.

The Manifesto describes zburn as: "Disposable identities — burner emails, names, addresses, phone numbers, passwords. Never give a service your real information again."

## Requirements
- Title: `zburn` with a one-line tagline from the Manifesto
- Brief description (2-3 sentences max) explaining what it does and why
- Install section: `go install` and `brew install zarlcorp/tap/zburn` (note: brew formula is placeholder until first release)
- Usage section showing both TUI mode (`zburn`) and CLI subcommands (`zburn version`, `zburn email`, `zburn identity`, `zburn list`, `zburn forget <id>`)
- Development section: `make build`, `make test`, `make lint`
- Link to core repo and Manifesto
- MIT license footer
- No badges, no screenshots, no emoji, no verbose prose
- Match the tone of the core README: terse, direct, lowercase comment style

## Target Repo
zarlcorp/zburn

## Agent Role
backend

## Files to Create
- README.md

## Notes
zburn is pre-release — most CLI subcommands are not yet implemented. The README should describe the intended interface (as documented in the Manifesto) but should not claim features that don't exist. Use future-tense or document the commands without claiming they work yet. The `version` subcommand and TUI mode do work.
