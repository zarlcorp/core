# 095: zstore — encrypted collection storage

## Objective
Create `pkg/zstore`, a generic encrypted key-value store built on `zfilesystem` and `zcrypto`. This is the shared storage foundation for zburn, zvault, and any future tool needing encrypted local persistence.

## Context
zburn currently has a bespoke `internal/store` that encrypts identity files via AES-256-GCM with Argon2id key derivation. zvault needs the same pattern for secrets, tasks, notes. Rather than duplicate, we extract the generic parts into a shared package in core.

The store uses HKDF (`zcrypto.ExpandKey`) to derive per-collection sub-keys from the master key, so compromising one collection's key doesn't expose others. The collection name is used as the HKDF info parameter.

## Requirements

### Store
- `Open(fs zfilesystem.ReadWriteFileFS, password []byte, opts ...Option) (*Store, error)` — creates or opens a store
  - On first run: generates salt, creates verification token, stores both via the filesystem
  - On subsequent runs: reads salt, derives key, verifies password by decrypting the verification token
  - Salt: 16 bytes via `zcrypto.RandBytes`, stored as `salt` file
  - Verification: encrypt a known token (`"zstore-verify-ok"`) with the master key, store as `verify` file
  - Master key derived via `zcrypto.DeriveKey(password, salt)`
- `Close() error` — erases master key and all sub-keys from memory via `zcrypto.Erase`
- Variadic options following `zoptions` pattern (unexported fields, exported option functions)

### Collection
- `Collection[V any](store *Store, name string) *Collection[V]` — returns a typed collection
  - Creates a subdirectory named `name` via `MkdirAll`
  - Derives a sub-key via `zcrypto.ExpandKey(masterKey, salt, []byte(name))`
  - V must be JSON-serializable
- `Put(id string, value V) error` — marshal to JSON, encrypt with collection sub-key, write to `{name}/{id}.enc`
- `Get(id string) (V, error)` — read file, decrypt, unmarshal. Return `zstore.ErrNotFound` if file doesn't exist
- `Delete(id string) error` — remove the file. Return `zstore.ErrNotFound` if file doesn't exist
- `List() ([]V, error)` — walk the collection directory, decrypt and unmarshal all `.enc` files, return unsorted (consumers sort as needed)
- `Len() (int, error)` — count `.enc` files without decrypting

### Error sentinels
- `var ErrNotFound = errors.New("not found")`
- `var ErrWrongPassword = errors.New("wrong password")`

### File layout on disk
```
<root>/
├── salt           (16 bytes, unencrypted)
├── verify         (AES-256-GCM encrypted verification token)
├── identities/    (collection)
│   ├── abc123.enc
│   └── def456.enc
├── credentials/   (collection)
│   └── ...
└── config/        (collection)
    └── ...
```

### Testing
- All tests use `zfilesystem.NewMemFS()` — no disk, no cleanup
- Test Open (first run + subsequent run + wrong password)
- Test Collection CRUD (Put, Get, Delete, List, Len)
- Test ErrNotFound on Get and Delete for missing keys
- Test multiple collections with different types share a store but are isolated
- Test Close erases keys (verify key slice is zeroed)
- Contract-style: if we want multiple store backends later, the test suite should be reusable

### Dependencies
- `github.com/zarlcorp/core/pkg/zcrypto` (DeriveKey, ExpandKey, Encrypt, Decrypt, Erase, RandBytes)
- `github.com/zarlcorp/core/pkg/zfilesystem` (ReadWriteFileFS)

## Target Repo
zarlcorp/core

## Agent Role
backend

## Files to Modify
- pkg/zstore/zstore.go (Store type, Open, Close)
- pkg/zstore/collection.go (Collection type, CRUD)
- pkg/zstore/options.go (variadic options)
- pkg/zstore/zstore_test.go (tests)
- pkg/zstore/go.mod, go.sum

## Notes
This is the foundation for specs 098-105. No external dependencies beyond what zcrypto already has. The API mirrors zburn's existing store closely but generalized with generics and per-collection sub-keys.
