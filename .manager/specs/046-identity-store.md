# 046: Encrypted Identity Store

## Objective
Build the persistence layer that stores and retrieves identities encrypted at rest.

## Context
zburn needs to save generated identities so users can recall them later. All storage is local, encrypted with a master password. Uses zcrypto (AES-256-GCM + Argon2id) and zfilesystem for testability.

## Requirements

### Store interface and implementation
Create `internal/store/store.go`:
```go
type Store struct {
    fs   zfilesystem.ReadWriteFileFS
    key  []byte // derived from master password
    salt []byte
}

func Open(fs zfilesystem.ReadWriteFileFS, password string) (*Store, error)
func (s *Store) Save(id identity.Identity) error
func (s *Store) Get(id string) (identity.Identity, error)
func (s *Store) List() ([]identity.Identity, error)
func (s *Store) Delete(id string) error
func (s *Store) Close() error
```

### Storage format
- Base directory: the store receives a `ReadWriteFileFS` rooted at the data dir. The caller decides the path (typically `~/.local/share/zburn/` or `~/.zburn/`).
- Salt file: `salt` — 16 bytes, created on first `Open`, reused on subsequent opens
- Identity files: `identities/<id>.enc` — each identity is JSON-marshaled, then encrypted with AES-256-GCM
- On `Open`: read salt (or create if first run), derive key from password+salt via `zcrypto.DeriveKey`
- On `Close`: erase key from memory via `zcrypto.Erase`

### Operations
- `Save`: marshal identity to JSON, encrypt with `zcrypto.Encrypt`, write to `identities/<id>.enc`
- `Get`: read `identities/<id>.enc`, decrypt with `zcrypto.Decrypt`, unmarshal JSON
- `List`: walk `identities/` dir, decrypt each file, return slice sorted by `CreatedAt` descending
- `Delete`: remove `identities/<id>.enc` file. Use `zcrypto.Erase` on the decrypted content before returning.

### Error handling
- Wrong password: `Open` should still succeed (we can't verify the password until we try to decrypt). `Get`/`List` will return a decrypt error if the password is wrong.
- OR: store a known verification token encrypted with the key. On `Open`, try to decrypt it — if it fails, return an error immediately.
- Prefer the verification token approach for better UX.

### Testing
- Use `zfilesystem.NewMemFS()` for all tests — no real filesystem
- Test round-trip: save then get
- Test list returns all saved identities
- Test delete removes identity
- Test wrong password fails verification
- Test first-run creates salt and verification token
- Test re-open with correct password succeeds

## Target Repo
zarlcorp/zburn

## Agent Role
backend

## Files to Create
- internal/store/store.go
- internal/store/store_test.go

## Dependencies
- `github.com/zarlcorp/core/pkg/zcrypto` — encryption, key derivation, secure erase
- `github.com/zarlcorp/core/pkg/zfilesystem` — filesystem abstraction
- `internal/identity` (from spec 045) — Identity type

## Notes
- The store receives a `ReadWriteFileFS` — it does NOT decide where files go. The caller (main.go) creates an `OSFileSystem` rooted at the user's data dir and passes it in. This keeps the store testable.
- Import the Identity type from `internal/identity`. If 045 isn't complete yet, define a minimal local type and swap later.
- Use `zcrypto.Erase` on sensitive data (keys, plaintext) when done with it.
- The `encoding/json` package is fine for serialization — identities are small.
