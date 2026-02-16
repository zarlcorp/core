# 097: zcrypto password generator extraction

## Objective
Extract the password generator from zburn's `internal/identity` package into `pkg/zcrypto` so it's available as a shared utility.

## Context
zburn currently generates passwords in `internal/identity/generator.go`. The credential vault will also need password generation, and zvault may want it too. The generator uses `crypto/rand` and is a natural fit for zcrypto alongside `RandBytes` and `RandHex`.

## Requirements

### Function
- `GeneratePassword(length int, opts ...PasswordOption) string`
  - Minimum length: 4 (auto-clamped like current implementation)
  - Default: guarantees at least one character from each class (lower, upper, digit, symbol)
  - Uses `crypto/rand` for all randomness
  - Fisher-Yates shuffle after placing guaranteed characters
  - Returns the password as a string
  - Panics on `crypto/rand` failure (same as current — unrecoverable)

### Options (variadic)
- `WithoutSymbols()` — exclude symbols from the character set (some services don't allow them)
- `WithCharset(chars string)` — override the full character set (advanced use)
- Keep options minimal for now — can add more later

### Character classes (same as current zburn)
- Lower: `abcdefghijklmnopqrstuvwxyz`
- Upper: `ABCDEFGHIJKLMNOPQRSTUVWXYZ`
- Digits: `0123456789`
- Symbols: `!@#$%^&*()-_=+[]{}|;:,.<>?`

### Testing
- Test default generation produces valid passwords of requested length
- Test minimum length clamping (length < 4 → 4)
- Test all character classes are represented in default mode
- Test WithoutSymbols excludes symbols
- Test randomness: generate 100 passwords, verify they're not all identical
- Test various lengths (4, 8, 20, 64, 128)

## Target Repo
zarlcorp/core

## Agent Role
backend

## Files to Modify
- pkg/zcrypto/password.go
- pkg/zcrypto/password_test.go

## Notes
Independent — no dependency on other specs. After this lands, spec 098 will update zburn to import `zcrypto.GeneratePassword` instead of using its internal implementation. The internal generator in zburn stays for identity-specific generation (names, emails, etc) but delegates password generation to zcrypto.
