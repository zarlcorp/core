# 031: Build zcrypto package

## Objective
Implement the `zcrypto` package — encryption primitives for zarlcorp privacy tools. Provides symmetric encryption (AES-256-GCM), key derivation (Argon2id), secure random generation, and secure memory erasure. Composition of proven Go stdlib and x/crypto primitives — no custom cryptography.

## Context
zcrypto sits in the security layer of the package hierarchy, alongside zstyle. It will be consumed by zburn (encrypted identity storage) and zvault (encrypted secret storage). Both products need the same crypto primitives, so they live in core.

The `go.mod` already exists at `pkg/zcrypto/go.mod` with module path `github.com/zarlcorp/core/pkg/zcrypto`. The module is already listed in `go.work`.

## Requirements

### Symmetric encryption — `aes.go`

AES-256-GCM authenticated encryption:

```go
// Encrypt encrypts plaintext with the given key using AES-256-GCM.
// Key must be 32 bytes. Returns nonce prepended to ciphertext.
func Encrypt(key, plaintext []byte) ([]byte, error)

// Decrypt decrypts ciphertext produced by Encrypt.
// Key must be 32 bytes. Expects nonce prepended to ciphertext.
func Decrypt(key, ciphertext []byte) ([]byte, error)
```

- Key must be exactly 32 bytes (AES-256), return error otherwise
- Nonce generated via `crypto/rand`, prepended to output
- Decrypt validates nonce length before attempting decryption
- Use `crypto/aes` + `crypto/cipher` from stdlib

### Key derivation — `kdf.go`

Argon2id for password-based key derivation:

```go
// DeriveKey derives a 32-byte key from a password and salt using Argon2id.
func DeriveKey(password, salt []byte) []byte

// DeriveKeyWithParams derives a key with custom Argon2id parameters.
func DeriveKeyWithParams(password, salt []byte, time, memory uint32, threads uint8) []byte
```

- Default parameters: time=1, memory=64*1024 (64MB), threads=4, keyLen=32
- Uses `golang.org/x/crypto/argon2`
- Salt should be at least 16 bytes — document this, but don't enforce (caller's responsibility)

### Secure random — `rand.go`

```go
// Bytes returns n cryptographically secure random bytes.
func Bytes(n int) ([]byte, error)

// Salt returns a 16-byte random salt suitable for key derivation.
func Salt() ([]byte, error)
```

- Uses `crypto/rand.Read`
- Wraps the error with context

### Secure memory erasure — `mem.go`

```go
// Zero overwrites a byte slice with zeros.
// Call this to erase sensitive data (keys, passwords) from memory.
func Zero(b []byte)
```

- Simple `for range` zero fill
- Not guaranteed against compiler optimizations, but best-effort. Document this caveat.

### What zcrypto does NOT do (yet)
- No file encryption helpers — add when a product needs them
- No age-compatible encryption — add when interop is needed
- No HKDF key expansion — add when needed
- No asymmetric encryption — not needed for local-only tools

### Package doc
Brief package comment on the main file (or doc.go) explaining zcrypto provides encryption primitives. Show a quick usage example:

```go
import "github.com/zarlcorp/core/pkg/zcrypto"

salt, _ := zcrypto.Salt()
key := zcrypto.DeriveKey([]byte("passphrase"), salt)
ciphertext, _ := zcrypto.Encrypt(key, []byte("secret"))
plaintext, _ := zcrypto.Decrypt(key, ciphertext)
defer zcrypto.Zero(key)
```

### Tests — `zcrypto_test.go`

Table-driven tests covering:
- Encrypt/Decrypt round-trip with valid key
- Encrypt with wrong key size returns error
- Decrypt with wrong key returns error (authentication failure)
- Decrypt with truncated ciphertext returns error
- DeriveKey produces deterministic output for same password+salt
- DeriveKey produces different output for different passwords
- DeriveKey produces different output for different salts
- Bytes returns requested length
- Bytes returns different values on successive calls
- Salt returns 16 bytes
- Zero clears the slice

## Acceptance Criteria
1. Encrypt/Decrypt round-trip works with 32-byte keys
2. Invalid key sizes produce clear errors
3. Wrong key on decrypt produces authentication error
4. DeriveKey is deterministic for same inputs
5. Random generation uses crypto/rand
6. Zero clears memory
7. All tests pass (`go test ./...` from `pkg/zcrypto/`)
8. `go.mod` has `golang.org/x/crypto` dependency for argon2
9. No dependency on any other zarlcorp package
10. Package compiles cleanly

## Target Repo
zarlcorp/core

## Agent Role
backend

## Files to Create/Modify
- `pkg/zcrypto/doc.go` — package doc with usage example (replace stub)
- `pkg/zcrypto/aes.go` — Encrypt, Decrypt
- `pkg/zcrypto/kdf.go` — DeriveKey, DeriveKeyWithParams
- `pkg/zcrypto/rand.go` — Bytes, Salt
- `pkg/zcrypto/mem.go` — Zero
- `pkg/zcrypto/zcrypto_test.go` — table-driven tests
- `pkg/zcrypto/go.mod` — add x/crypto dependency

## Dependencies
Depends on 029 (Go upgrade) — needs go.mod at correct version.

## Notes
- Only external dependency: `golang.org/x/crypto` for argon2
- Keep error messages direct: "decrypt: ciphertext too short", not "failed to decrypt: invalid ciphertext length"
- The Argon2id defaults (64MB memory) are reasonable for desktop tools. Server-side tools might want different params — that's why DeriveKeyWithParams exists.
- AES-256-GCM nonce is 12 bytes (standard for GCM)
- This is library code — must be bulletproof, prevent panics, validate all inputs
