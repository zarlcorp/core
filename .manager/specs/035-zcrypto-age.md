# 035: zcrypto — age-compatible encryption

## Objective
Add age-format encryption and decryption to `zcrypto`, enabling zarlcorp tools to produce and consume files compatible with the `age` CLI (`filippo.io/age`).

## Context
The zcrypto package currently provides AES-256-GCM encryption with custom framing (nonce prepended to ciphertext). This works for internal tool storage but cannot interoperate with external tools. Age is the de facto standard for simple file encryption in the Go ecosystem — users expect to `age -d` a file exported from zvault, or `age -e` a file for import.

The manifesto lists age-compatible encryption as a planned zcrypto capability. Issue #21.

## Requirements

### Password-based age encryption
- `EncryptAge(password string, src io.Reader, dst io.Writer) error` — encrypt using age's scrypt recipient format.
- `DecryptAge(password string, src io.Reader, dst io.Writer) error` — decrypt an age-encrypted stream.
- Output must be compatible with `age -d -` using the same password.
- Input from `age -e -p` must be decryptable.

### Key-based age encryption (stretch goal)
- `EncryptAgeKey(recipients []string, src io.Reader, dst io.Writer) error` — encrypt to one or more age public keys (X25519).
- `DecryptAgeKey(identity string, src io.Reader, dst io.Writer) error` — decrypt with an age identity (private key).
- Compatible with `age -r <pubkey>` and `age -i <identity>`.

### Implementation approach
- Use `filippo.io/age` as a dependency — do not reimplement the format. Age is MIT-licensed, well-audited, maintained by Filippo Valsorda.
- Wrap the age API with zcrypto's error handling conventions.
- Streaming — age supports streaming encryption, so these functions should not buffer the entire file in memory (unlike the existing `EncryptFile`/`DecryptFile`).

### Tests
- Round-trip: encrypt with zcrypto, decrypt with zcrypto.
- Interop: encrypt with zcrypto, decrypt with `age` CLI (if available in test environment, otherwise skip with `t.Skip`).
- Interop: encrypt with `age` CLI, decrypt with zcrypto (same skip logic).
- Wrong password: decryption fails cleanly.
- Empty input: defined behavior.

### No breaking changes
- Existing `Encrypt`/`Decrypt`/`EncryptFile`/`DecryptFile` remain unchanged.
- Age functions are additive.

## Target Repo
zarlcorp/core

## Agent Role
backend

## Files to Modify
- `pkg/zcrypto/age.go` — new file for age encryption functions
- `pkg/zcrypto/zcrypto_test.go` — add age test cases (or new `age_test.go`)
- `pkg/zcrypto/go.mod` — add `filippo.io/age` dependency
- `pkg/zcrypto/doc.go` — update package docs with age examples

## Notes
- Password-based is the priority. Key-based is a stretch goal — implement only if the API is clean and simple.
- The `filippo.io/age` library handles all the hard parts (scrypt params, header format, STREAM chunking). Our wrapper is thin.
- Error messages: `"age encrypt: %w"`, `"age decrypt: %w"` — direct context, no stuttering.
