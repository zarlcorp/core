# 117: Realistic email generation

## Objective
Generate realistic email addresses for burner identities using the identity's first/last name, with variety across multiple patterns. Use Namecheap domains when configured, with rotation and a cycle button.

## Context
Currently `Email()` generates random strings like `sparklyoctopus4821@zburn.id` — disconnected from the generated name and using a hardcoded domain. After spec 116 simplifies Namecheap settings and fetches domains, we can use real domains. The email local part should look like a real person's email.

## Requirements

### 1. Name-based email patterns
Modify the identity generator so `Email` uses the identity's first/last name. Support these patterns, chosen randomly per identity:

- `firstname.lastname` — `john.doe`
- `firstinitiallastname` — `jdoe`
- `firstnamelastname` — `johndoe`
- `firstname.lastname` + 2 digits — `john.doe42`
- `firstinitiallastname` + 2 digits — `jdoe42`
- `firstinitial.lastname` — `j.doe`
- `lastname.firstname` — `doe.john`
- `adjective` + `noun` + 4 digits — `sparklyoctopus4821` (current random style, kept for variety)

All names lowercased. The random pattern selection uses `crypto/rand` like everything else.

### 2. Change Email signature
`Email()` currently takes no arguments. Change to:
```go
func (g *Generator) Email(firstName, lastName, domain string) string
```
Update `Generate()` to pass the already-generated name and domain through.

### 3. Domain support
- `Generate()` needs a domain parameter or the generator needs domain config
- Option: change `Generate()` signature to accept a domain: `Generate(domain string) Identity`
- When domain is empty, fall back to `zburn.id`
- The TUI will pass the current rotated domain from Namecheap settings

### 4. Domain rotation in TUI
- When creating a new identity, auto-select a domain from `CachedDomains` (from spec 116)
- Rotate through domains round-robin across identity generations
- Add a key binding to cycle to the next domain before generating (show current domain in the generate view)
- If no Namecheap domains configured, use `zburn.id` as default

### 5. Update tests
- Test each email pattern produces valid format
- Test that the name is correctly incorporated
- Test domain parameter works (custom domain and fallback)
- Update any existing email tests for the new signature

## Target Repo
zarlcorp/zburn

## Agent Role
backend

## Files to Modify
- internal/identity/generator.go (new Email signature, name-based patterns)
- internal/identity/generator_test.go (test all patterns)
- internal/identity/data.go (may need to keep adjectives/nouns for the random pattern)
- internal/tui/tui.go (domain rotation state, cycle key, pass domain to Generate)

## Notes
Depends on spec 116 for `CachedDomains` in settings. However, the generator changes can be built independently — the generator just accepts a domain string. The TUI integration that reads from CachedDomains is the part that depends on 116.
