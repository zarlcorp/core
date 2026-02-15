# 047: CLI Subcommands

## Objective
Wire up zburn's CLI subcommands so users can generate identities, list saved ones, and forget them from the command line.

## Context
zburn currently only has a `version` subcommand. The Manifesto specifies these CLI commands:
- `zburn email` — generate and print a burner email
- `zburn identity` — generate a complete identity
- `zburn list` — show stored identities
- `zburn forget <id>` — securely erase an identity

This spec wires up the identity generator (045) and store (046) into the CLI.

## Requirements

### Subcommands
Extend `cmd/zburn/main.go` to handle:

**`zburn email`**
- Generate a random email using `identity.Generator`
- Print it to stdout
- No store interaction needed

**`zburn identity`**
- Generate a complete identity using `identity.Generator`
- Print all fields in a readable format
- Prompt: "Save this identity? [y/N]" — if yes, prompt for master password and save to store
- Use `--save` flag to skip the prompt and always save
- Use `--json` flag to output as JSON instead of formatted text

**`zburn list`**
- Prompt for master password
- List all saved identities: ID, name, email, created date
- Use `--json` flag for JSON output
- If no identities saved, print "no saved identities"

**`zburn forget <id>`**
- Prompt for master password
- Delete the identity with the given ID
- Print confirmation or error

### Password prompting
- Use `term.ReadPassword` from `golang.org/x/term` for secure password input (no echo)
- Prompt: "master password: "
- On first use (no salt file exists), prompt twice to confirm: "master password: " then "confirm password: "

### Data directory
- Default: `$XDG_DATA_HOME/zburn` if set, otherwise `~/.local/share/zburn`
- Create the directory if it doesn't exist

### Output formatting
For `zburn identity` (non-JSON), print:
```
  id:       a1b2c3d4
  name:     Jane Smith
  email:    brightfox4521@zburn.id
  phone:    (555) 234-5678
  address:  742 Oak Ave, Portland, OR 97201
  dob:      1985-03-14
  password: Kx#m9pQ2wR!vN8sL4jYt
```

For `zburn list` (non-JSON), print:
```
  a1b2c3d4  Jane Smith      brightfox4521@zburn.id  2026-02-15
  e5f6g7h8  John Doe        quietmoon8832@zburn.id  2026-02-14
```

### Testing
- Test that each subcommand is recognized (no "unknown command" error)
- Test output formatting functions
- Integration-style tests are optional — the core logic is tested in 045 and 046

## Target Repo
zarlcorp/zburn

## Agent Role
backend

## Files to Modify
- cmd/zburn/main.go — add subcommand dispatch

## Files to Create
- internal/cli/cli.go — subcommand implementations
- internal/cli/cli_test.go — tests

## Dependencies
- `internal/identity` (spec 045) — Generator, Identity type
- `internal/store` (spec 046) — Store for save/list/forget
- `golang.org/x/term` — secure password input

## Notes
- Keep the CLI thin — it should be a wrapper around the identity generator and store. Business logic lives in those packages, not here.
- Use `os.UserHomeDir()` and `os.Getenv("XDG_DATA_HOME")` for the data directory.
- The `--json` flag is useful for scripting: `zburn identity --json | jq .email`
- Don't add a flag parsing library. Use `os.Args` directly — the command set is small enough.
