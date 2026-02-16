# 118: TUI domain rotation for identity generation

## Objective
Wire up Namecheap cached domains into the identity generation flow. Show the current domain next to the email field, allow cycling with space, and pass the selected domain to `Generate()`.

## Context
Spec 116 added `CachedDomains` to `NamecheapSettings`. Spec 117 changed `Generate(domain string)` to accept a domain parameter. The TUI currently passes `""` (defaulting to `zburn.id`). This spec connects the two.

## Requirements

### 1. Domain state on the root model
- Add `domains []string` and `domainIdx int` to the root `Model` in `tui.go`
- On startup (after loading settings), populate `domains` from `NamecheapSettings.CachedDomains`
- If `domains` is empty, leave it nil — `Generate("")` already falls back to `zburn.id`
- `domainIdx` starts at 0 each session, no persistence

### 2. Pass domain to Generate
- In the `viewGenerate` navigation handler (`tui.go` ~line 296), pass the current domain:
  ```go
  domain := ""
  if len(m.domains) > 0 {
      domain = m.domains[m.domainIdx]
  }
  id := m.gen.Generate(domain)
  ```

### 3. Show domain next to email in generate view
- In `generate.go`, the email field currently shows just `id.Email`
- Add a `domain` field to `generateModel` so the view knows what domain is active
- Update `newGenerateModel` to accept the current domain
- In `View()`, next to the email line, show the domain hint: `[domain.com]  space to cycle`
- Only show the hint if a domain was provided (not when using the `zburn.id` fallback)
- Use `zstyle.MutedText` for the hint

### 4. Space key cycles domain
- In `generate.go` `handleKey`, add a `" "` (space) handler:
  - Send a new message type `cycleDomainMsg{}` up to the root model
- In `tui.go`, handle `cycleDomainMsg`:
  - Advance `domainIdx = (domainIdx + 1) % len(domains)`
  - Regenerate the identity with the new domain: `m.gen.Generate(newDomain)`
  - Rebuild `m.generate = newGenerateModel(id, newDomain)`
- Only cycle if `len(domains) > 1` — no point cycling with one domain
- If no domains configured, space does nothing

### 5. Update help text
- In `generate.go` `View()`, add `space domain` to the help line when domains are available
- When no domains configured, omit it

### 6. Update tests
- Test that `cycleDomainMsg` advances `domainIdx` and regenerates with the new domain
- Test that space key in generate view produces `cycleDomainMsg`
- Test that no domains configured = space does nothing
- Test that single domain = space does nothing
- Test that domain hint appears in View() when domain is set
- Test that domain hint is absent when using default

## Target Repo
zarlcorp/zburn

## Agent Role
backend

## Files to Modify
- internal/tui/tui.go (domain state, pass to Generate, handle cycleDomainMsg)
- internal/tui/generate.go (domain field, space handler, view hint, help text)
- internal/tui/tui_test.go (cycle tests)
- internal/tui/integration_test.go (if any generate tests need updating)

## Notes
The `n` key already regenerates a new identity from within the generate view. When pressing `n`, it should also use the current domain (not reset it). The navigation handler for `viewGenerate` already covers this since `n` sends `navigateMsg{view: viewGenerate}`.
