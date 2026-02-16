# 108: zburn v0.5.0 release

## Objective
Cut a new zburn release (v0.5.0) that includes all the feature expansion work from specs 098-105.

## Context
v0.4.0 was released before the feature expansion landed. Main now has 10 additional commits including: zstore migration, credential vault, Namecheap/Gmail/Twilio clients, 2FA code extractor, settings TUI, and burn cascade. This is a major feature release.

## Requirements

### Release notes
Create a GitHub release with comprehensive notes covering:
- **Credential vault** — store passwords, TOTP secrets, URLs per persona
- **Settings** — configure Namecheap, Gmail, Twilio integrations from the TUI
- **Burn cascade** — confirmation dialog + best-effort cleanup of all associated resources
- **Email forwarding** — Namecheap API integration for persona email
- **Gmail OAuth2** — read emails from a shared inbox
- **Twilio SMS** — provision phone numbers, send/receive SMS
- **2FA codes** — extract verification codes from email/SMS
- **Encrypted storage** — migrated to zstore with HKDF-derived per-collection keys

### Release process
1. Ensure all tests pass on main
2. Tag `v0.5.0` on current main HEAD
3. Create GitHub release with notes
4. Verify homebrew tap auto-updates (spec 078 workflow)

## Target Repo
zarlcorp/zburn

## Agent Role
PM task — not delegated to agents

## Depends On
- 107 (integration tests should pass before release)

## Notes
The version is injected via ldflags at build time, so no code changes are needed. This is a PM-executed task: tag, release, verify homebrew. The homebrew tap update should trigger automatically via the GitHub Actions workflow from spec 078.
