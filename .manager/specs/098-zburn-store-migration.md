# 098: zburn store migration + credential model

## Objective
Replace zburn's `internal/store` with `pkg/zstore` from core, add a credential model, and update the identity generator to use `zcrypto.GeneratePassword`.

## Context
Spec 095 creates `pkg/zstore` as a generic encrypted collection store. Spec 097 extracts password generation to zcrypto. This spec migrates zburn to use both, and adds the credential data model that the TUI (spec 103) will present.

Clean break — no migration of existing encrypted data. Users re-create identities.

## Requirements

### Store migration
- Remove `internal/store/` package entirely
- Replace with `zstore.Open(fs, password)` in the app initialization
- Identity collection: `zstore.Collection[identity.Identity](store, "identities")`
- Update all call sites that reference the old store (TUI models, CLI commands)
- The `identity.Identity` struct stays in `internal/identity` — it's zburn-specific

### Credential model
- Create `internal/credential/credential.go`:
  ```go
  type Credential struct {
      ID         string    `json:"id"`
      IdentityID string    `json:"identity_id"`
      Label      string    `json:"label"`
      URL        string    `json:"url"`
      Username   string    `json:"username"`
      Password   string    `json:"password"`
      TOTPSecret string    `json:"totp_secret,omitempty"`
      Notes      string    `json:"notes,omitempty"`
      CreatedAt  time.Time `json:"created_at"`
      UpdatedAt  time.Time `json:"updated_at"`
  }
  ```
- Credential collection: `zstore.Collection[credential.Credential](store, "credentials")`
- ID generation: 8-character hex (same as identity)

### Password generator swap
- `internal/identity/generator.go`: replace the `Password` method body with a call to `zcrypto.GeneratePassword(length)`
- Remove the password character class constants and `pickByte` from generator.go (they move to zcrypto)
- Keep `pick`, `randIntn`, `mustRead` — those are used for identity generation (names, addresses)

### Testing
- All existing tests must pass after migration
- Add test that credentials round-trip through the store (put + get)
- Add test that deleting an identity doesn't delete its credentials (they're separate collections)
- Test credential ID generation

## Target Repo
zarlcorp/zburn

## Agent Role
backend

## Depends On
- 095 (zstore)
- 097 (zcrypto password generator)

## Files to Modify
- internal/store/ (delete entire package)
- internal/credential/credential.go (new)
- internal/identity/generator.go (swap password generation to zcrypto)
- internal/tui/*.go (update store references)
- internal/cli/cli.go (update store references if applicable)
- cmd/zburn/main.go (update store initialization)
- go.mod (update core dependency, remove old store deps if any)

## Notes
This is the bridge between the core infrastructure (095, 097) and the TUI work (103, 104, 105). The credential model is defined here but not yet surfaced in the TUI — that's spec 103.
