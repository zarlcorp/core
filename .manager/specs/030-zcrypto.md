# 030: Build zcrypto package

## Objective
Implement the `zcrypto` package — encryption primitives for zarlcorp privacy tools. Provides AES-256-GCM symmetric encryption, key derivation (Argon2id, HKDF), secure random generation, secure memory erasure, and file encryption/decryption helpers. No custom cryptography — composition of proven Go stdlib and `x/crypto` primitives only.

## Context
zcrypto sits in the security layer of the package hierarchy, alongside zstyle. Every privacy tool (zburn, zvault, zshield) will import it for encryption operations. The package must be bulletproof — library code that prevents panics at all costs.

The `go.mod` already exists at `pkg/zcrypto/go.mod` with module path `github.com/zarlcorp/core/pkg/zcrypto`. The module is already listed in `go.work`.

Age-compatible encryption is deferred until a tool needs it.

## Requirements

### Symmetric encryption — `aes.go`

AES-256-GCM authenticated encryption.

```go
// Encrypt encrypts plaintext using AES-256-GCM with the given key.
// Key must be exactly 32 bytes. Returns ciphertext with nonce prepended.
func Encrypt(key, plaintext []byte) ([]byte, error)

// Decrypt decrypts ciphertext produced by Encrypt.
// Key must be exactly 32 bytes. Expects nonce prepended to ciphertext.
func Decrypt(key, ciphertext []byte) ([]byte, error)
```

Implementation details:
- Use `crypto/aes` + `crypto/cipher` (GCM mode)
- Generate random 12-byte nonce via `crypto/rand` for each Encrypt call
- Prepend nonce to ciphertext: `nonce || ciphertext`
- Decrypt extracts nonce from first 12 bytes
- Validate key length (must be 32 bytes), return error if wrong
- Validate ciphertext minimum length (must be > nonce size)

### Key derivation — `kdf.go`

Two key derivation functions for different use cases.

```go
// DeriveKey derives a 32-byte key from a password and salt using Argon2id.
// If salt is nil, a random 16-byte salt is generated.
// Returns the derived key and the salt used.
func DeriveKey(password, salt []byte) (key, usedSalt []byte, err error)

// ExpandKey derives a new key from existing key material using HKDF-SHA256.
// info provides context for domain separation (e.g. "file-encryption", "auth-token").
func ExpandKey(secret, salt, info []byte) ([]byte, error)
```

Implementation details:
- Argon2id parameters: time=1, memory=64*1024 (64MB), threads=4, keyLen=32
- Define these as package-level constants so they're documentated
- Salt for Argon2id: 16 bytes, generated via `crypto/rand` if nil
- HKDF uses `crypto/sha256` as the hash, outputs 32 bytes
- Dependencies: `golang.org/x/crypto/argon2`, `golang.org/x/crypto/hkdf`

### Secure random — `rand.go`

```go
// RandBytes returns n cryptographically random bytes.
func RandBytes(n int) ([]byte, error)

// RandHex returns a hex-encoded string of n random bytes (2n chars).
func RandHex(n int) (string, error)
```

Implementation: use `crypto/rand.Read`. Return error if the system RNG fails (don't panic).

### Secure erasure — `erase.go`

```go
// Erase zeroes out a byte slice to prevent sensitive data from lingering in memory.
func Erase(b []byte)
```

Implementation: iterate and zero each byte. Use a loop that the compiler can't optimize away (write through a volatile-style pattern or use `runtime.KeepAlive`-adjacent technique). This is best-effort — Go's GC can copy memory, but zeroing the known location is still worthwhile.

### File encryption — `file.go`

```go
// EncryptFile encrypts the contents of src and writes the result to dst.
// Uses AES-256-GCM with the given 32-byte key.
// File format: salt (16 bytes) || nonce (12 bytes) || ciphertext
func EncryptFile(key []byte, src io.Reader, dst io.Writer) error

// DecryptFile decrypts a file encrypted by EncryptFile.
func DecryptFile(key []byte, src io.Reader, dst io.Writer) error
```

Implementation details:
- Read entire src into memory (privacy tools deal with small files — secrets, identities, not gigabyte blobs)
- Use `Encrypt`/`Decrypt` internally
- File format is simply the output of `Encrypt` (nonce || ciphertext)
- Errors use direct context: `"read source: %w"`, `"encrypt: %w"`, etc.

### What zcrypto does NOT do
- No age-compatible encryption (deferred)
- No asymmetric encryption (no RSA, no ECDH) — not needed for local-only tools
- No streaming encryption — files are read into memory
- No key storage or management — that's the consumer's job
- No dependency on other zarlcorp packages

### Package doc
Brief package comment on `aes.go` (or whichever file is first alphabetically) explaining zcrypto provides encryption primitives for zarlcorp tools. Short usage example showing Encrypt/Decrypt round-trip.

### Tests — `zcrypto_test.go`

Comprehensive table-driven tests. Security-critical code needs thorough coverage:

**Encrypt/Decrypt:**
- Round-trip: encrypt then decrypt recovers plaintext
- Different plaintexts produce different ciphertexts (nonce uniqueness)
- Wrong key fails to decrypt
- Tampered ciphertext fails to decrypt
- Empty plaintext works
- Invalid key length returns error
- Ciphertext too short returns error

**DeriveKey:**
- Same password + salt produces same key (deterministic)
- Different passwords produce different keys
- nil salt auto-generates and returns salt
- Provided salt is used as-is

**ExpandKey:**
- Same inputs produce same output (deterministic)
- Different info values produce different keys (domain separation)

**RandBytes/RandHex:**
- Returns correct length
- Two calls produce different output (statistical)

**Erase:**
- Slice is zeroed after call

**File encryption:**
- Round-trip through EncryptFile/DecryptFile via bytes.Buffer
- Wrong key fails

## Acceptance Criteria
1. `Encrypt`/`Decrypt` round-trip works with AES-256-GCM
2. `DeriveKey` uses Argon2id with documented parameters
3. `ExpandKey` uses HKDF-SHA256
4. `RandBytes`/`RandHex` use `crypto/rand`
5. `Erase` zeroes the slice
6. `EncryptFile`/`DecryptFile` work with `io.Reader`/`io.Writer`
7. All error messages use direct context (no "failed to" prefixes)
8. All tests pass (`go test -race ./...` from `pkg/zcrypto/`)
9. `go.mod` has `golang.org/x/crypto` dependency
10. No dependency on any other zarlcorp package
11. No panics — all error paths return errors

## Target Repo
zarlcorp/core

## Agent Role
backend

## Files to Create/Modify
- `pkg/zcrypto/aes.go` — Encrypt, Decrypt
- `pkg/zcrypto/kdf.go` — DeriveKey (Argon2id), ExpandKey (HKDF)
- `pkg/zcrypto/rand.go` — RandBytes, RandHex
- `pkg/zcrypto/erase.go` — Erase
- `pkg/zcrypto/file.go` — EncryptFile, DecryptFile
- `pkg/zcrypto/zcrypto_test.go` — comprehensive tests
- `pkg/zcrypto/go.mod` — add x/crypto dependency

## Dependencies
None — zcrypto is independent of other zarlcorp packages.

## Notes
- Use `golang.org/x/crypto` for argon2 and hkdf — these are maintained by the Go team
- Argon2id parameters (time=1, memory=64MB, threads=4) are OWASP recommended minimums
- The nonce-prepended ciphertext format is standard and self-contained — no separate nonce storage needed
- `Erase` is best-effort due to Go's GC. Document this clearly in the function comment.
- Error wrapping: `"decrypt: %w"`, `"derive key: %w"`, `"read source: %w"` — direct context, never "failed to"
- Run `go mod tidy` after adding x/crypto
