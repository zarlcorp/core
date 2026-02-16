# 101: Twilio SMS client

## Objective
Create an internal Twilio client in zburn for provisioning phone numbers and reading incoming SMS messages.

## Context
zburn personas need real phone numbers for 2FA verification. The user has a funded Twilio account with API access, supporting both UK (+44) and US (+1) numbers. Numbers are provisioned per persona and released on burn.

## Requirements

### Client
- Create `internal/twilio/client.go`:
  ```go
  type Config struct {
      AccountSID string
      AuthToken  string
  }

  type Client struct { ... }

  func NewClient(cfg Config) *Client
  ```
- All API calls go to `https://api.twilio.com/2010-04-01/Accounts/{AccountSID}/`
- Twilio uses HTTP Basic Auth (AccountSID:AuthToken)
- Responses are JSON (request with `.json` suffix)

### Methods

#### Number management
- `SearchNumbers(ctx context.Context, country string) ([]AvailableNumber, error)`
  - Calls `GET /AvailablePhoneNumbers/{country}/Local.json`
  - `country` is ISO country code: "GB" for UK, "US" for US
  - Returns a list of available numbers with capabilities
- `BuyNumber(ctx context.Context, phoneNumber string) (*PhoneNumber, error)`
  - Calls `POST /IncomingPhoneNumbers.json` with `PhoneNumber={number}`
  - Returns the provisioned number details including SID
- `ReleaseNumber(ctx context.Context, numberSID string) error`
  - Calls `DELETE /IncomingPhoneNumbers/{SID}.json`
  - Releases the number back to Twilio

#### SMS reading
- `ListMessages(ctx context.Context, to string, limit int) ([]SMSMessage, error)`
  - Calls `GET /Messages.json?To={to}&PageSize={limit}`
  - Returns recent messages sent TO the given number
  - Ordered by date descending (most recent first)

### Types
```go
type AvailableNumber struct {
    PhoneNumber  string // e.g. "+447123456789"
    FriendlyName string
    Capabilities struct {
        SMS   bool
        Voice bool
    }
}

type PhoneNumber struct {
    SID         string
    PhoneNumber string
    FriendlyName string
}

type SMSMessage struct {
    SID       string
    From      string
    To        string
    Body      string
    DateSent  time.Time
    Status    string
}
```

### Testing
- Test number search response parsing with canned JSON
- Test number purchase request/response
- Test number release (DELETE request)
- Test message listing and parsing
- All tests use httptest.Server — do NOT call real Twilio API
- Test error handling (insufficient funds, number unavailable, auth failure)
- Test Basic Auth header is correctly set

### Dependencies
- Standard library only: `net/http`, `encoding/json`, `net/url`

## Target Repo
zarlcorp/zburn

## Agent Role
backend

## Files to Modify
- internal/twilio/client.go
- internal/twilio/client_test.go

## Notes
Independent — no dependency on other specs. The client takes Config as a parameter; credential storage in zstore is handled by the settings TUI (spec 104). UK numbers are +44, US numbers are +1. The user wants both country options available.
