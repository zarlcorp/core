# 102: 2FA code extractor

## Objective
Create a shared internal utility in zburn that extracts 2FA verification codes from email and SMS message bodies.

## Context
Both the Gmail client (spec 100) and Twilio SMS client (spec 101) will fetch messages that contain 2FA codes. Rather than duplicate extraction logic, create a shared package that both can use. The extractor needs to handle the wide variety of formats services use for verification codes.

## Requirements

### Function
- Create `internal/codes/extract.go`:
  ```go
  // Extract returns all potential verification codes found in the text,
  // ordered by confidence (most likely first).
  func Extract(text string) []Code

  type Code struct {
      Value string // the code itself, e.g. "123456"
      Type  string // "numeric", "alphanumeric"
  }
  ```

### Patterns to match
- 4-digit numeric codes (e.g. "1234")
- 6-digit numeric codes (e.g. "123456") — most common
- 8-digit numeric codes (e.g. "12345678")
- Alphanumeric codes (e.g. "A1B2C3") — less common but used by some services

### Context clues (boost confidence)
- Text near the code contains keywords: "verification", "code", "OTP", "one-time", "confirm", "PIN", "security code", "2FA", "authenticate"
- Code appears after a colon, dash, or "is" (e.g. "Your code is: 123456")
- Code is on its own line or surrounded by whitespace

### Filtering
- Ignore numbers that are clearly not codes: years (2024, 2025, 2026), phone numbers (10+ digits), prices ($xx.xx), timestamps
- Ignore numbers embedded in URLs or email addresses
- 6-digit codes with nearby keywords should rank highest

### Testing
- Test common formats:
  - "Your verification code is 123456"
  - "123456 is your code"
  - "Code: 123456"
  - "Your OTP: 1234"
  - "Use 12345678 to verify"
  - "Enter A1B2C3 to confirm"
- Test filtering:
  - "Copyright 2025" should not extract "2025"
  - "Call us at 1234567890" should not extract as a code
  - "$123.45" should not extract
- Test multi-code: text with multiple potential codes returns all, highest confidence first
- Test empty/no-match returns empty slice

## Target Repo
zarlcorp/zburn

## Agent Role
backend

## Files to Modify
- internal/codes/extract.go
- internal/codes/extract_test.go

## Notes
Independent — no dependency on other specs. Pure string processing with regex. This is consumed by the TUI when displaying messages from Gmail or Twilio — the TUI can highlight the top code and offer a "copy code" action.
