# 045: Identity Generation Engine

## Objective
Build the core identity generation engine — pure functions that produce random but realistic-looking personal data.

## Context
zburn generates disposable identities so users never give real information to services. This is the foundation package that both the CLI and TUI will call. It has no I/O, no storage — just generation.

The Manifesto says: "Disposable identities — burner emails, names, addresses, phone numbers, passwords. Never give a service your real information again."

## Requirements

### Identity type
Create an `Identity` struct in `internal/identity/identity.go`:
```go
type Identity struct {
    ID        string    // short random hex ID (8 chars)
    FirstName string
    LastName  string
    Email     string    // <random>@zburn.id format
    Phone     string    // US format: (555) XXX-XXXX using 555 prefix (fictional)
    Street    string
    City      string
    State     string    // US state abbreviation
    Zip       string
    DOB       time.Time
    Password  string    // generated strong password
    CreatedAt time.Time
}
```

### Generator
Create a `Generator` in `internal/identity/generator.go`:
- `New() *Generator` — creates a generator
- `Generate() Identity` — generates a complete random identity
- `Email() string` — generates just an email
- `Password(length int) string` — generates a password of given length
- `Name() (first, last string)` — generates a name pair

### Data sources
Create `internal/identity/data.go` with embedded name/address pools:
- First names: ~100 common US first names (mix of gender)
- Last names: ~100 common US last names
- Cities: ~50 US cities
- States: all 50 US state abbreviations
- Street names: ~50 common street name patterns

### Generation rules
- Use `crypto/rand` for all randomness (never `math/rand`)
- Email format: `<adjective><noun><4digits>@zburn.id`
- Phone: always use 555 prefix (reserved for fictional use)
- Password: mix of upper, lower, digits, symbols. Default 20 chars.
- DOB: random date between 21 and 65 years ago
- Street: random number + random street name + random suffix (St, Ave, Blvd, Dr)

### Testing
- Table-driven tests for each generator method
- Test that generated values are non-empty and well-formed
- Test password length and character class requirements
- Test email format matches expected pattern
- Test DOB falls within expected age range
- Test that consecutive calls produce different results (randomness check)

## Target Repo
zarlcorp/zburn

## Agent Role
backend

## Files to Create
- internal/identity/identity.go (types)
- internal/identity/generator.go (generation logic)
- internal/identity/data.go (name/address pools)
- internal/identity/generator_test.go (tests)

## Notes
- No dependencies beyond stdlib. Use `crypto/rand` for randomness.
- The `@zburn.id` email domain is a placeholder — no actual mail forwarding yet. Keep it as a constant so it's easy to change later.
- Phone numbers MUST use 555 prefix — this is the North American reserved range for fictional use. Never generate real-looking phone numbers.
- This package is pure — no file I/O, no network, no side effects. Just data generation.
