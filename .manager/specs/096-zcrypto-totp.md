# 096: zcrypto TOTP code generation

## Objective
Add RFC 6238 TOTP code generation to `pkg/zcrypto`. Given a base32-encoded secret, generate the current 6-digit TOTP code.

## Context
zburn's credential vault will store TOTP secrets and generate codes on demand. zvault may also need this for stored secrets. TOTP is a natural fit for zcrypto since it's a cryptographic operation (HMAC-SHA1 + time).

## Requirements

### Functions
- `TOTPCode(secret string) (string, error)` — generates the current 6-digit TOTP code
  - `secret` is a base32-encoded string (standard encoding, case-insensitive, padding optional)
  - Uses current time with 30-second period (standard)
  - Uses HMAC-SHA1 (the RFC default, what 99% of services use)
  - Returns zero-padded 6-digit string (e.g. "007832")
  - Returns error if secret is not valid base32
- `TOTPCodeAt(secret string, t time.Time) (string, error)` — same but at a specific time (for testing)
- `TOTPCode` calls `TOTPCodeAt` with `time.Now()`

### Algorithm (RFC 6238 / RFC 4226)
1. Decode base32 secret to bytes
2. Calculate time step: `counter = floor(unix_time / 30)`
3. Encode counter as big-endian uint64 (8 bytes)
4. HMAC-SHA1(secret_bytes, counter_bytes) → 20-byte hash
5. Dynamic truncation: take 4 bytes at offset `hash[19] & 0x0f`
6. Convert to uint32, mask with `0x7fffffff`, mod `1000000`
7. Zero-pad to 6 digits

### Testing
- Test against known test vectors from RFC 6238 appendix B
- Test with secrets from real services (Google Authenticator format)
- Test invalid base32 returns error
- Test that codes change every 30 seconds (using TOTPCodeAt with times in different windows)
- Test zero-padding (codes with leading zeros)

### No external dependencies
- Standard library only: `crypto/hmac`, `crypto/sha1`, `encoding/base32`, `encoding/binary`, `time`

## Target Repo
zarlcorp/core

## Agent Role
backend

## Files to Modify
- pkg/zcrypto/totp.go
- pkg/zcrypto/totp_test.go

## Notes
Independent — no dependency on other specs. Small scope, ~40 lines of implementation.
