# 124: Credential form redesign with dual-mode fields

## Objective
Redesign the credential form with an identity header, dual-mode input fields (cycle/edit), username generation, and auto-generated passwords.

## Context
The credential form currently shows "add credential" with no context about which identity it belongs to. Username and password fields start empty, requiring manual entry. The user wants:
- A header showing the identity name/email
- Username pre-filled with the identity's email, with options to cycle through name-based and random handles
- Password auto-generated on form open, with space to regenerate
- A dual-mode field pattern: cycle mode (space cycles options) vs edit mode (free text)

## Requirements

### 1. Identity header in form
- Display the identity's name and email above the form fields
- Format: `Jane Doe  jane@domain.com` using `zstyle.MutedText` for context
- Requires passing identity info to the form (see req 5)

### 2. Dual-mode field behavior
Create a reusable pattern (not necessarily a separate type — can be logic within `credentialFormModel`) for fields that support two modes:

**Cycle mode** (default for generated fields):
- Space key → cycles to the next generated option
- Any printable character → switches to edit mode, the character is typed into the input
- Visual indicator: append `[generated]` in muted text after the field value, or style the value differently

**Edit mode** (user is typing custom text):
- Normal text input behavior — space types a literal space
- Esc key → returns to cycle mode, restores the current generated value
- Note: Esc from a non-dual-mode field (or when already in cycle mode) still navigates back as usual

**Important**: Only the username and password fields are dual-mode. Label, URL, TOTP secret, and notes remain plain text inputs. The `editing` flag (edit existing credential) should skip dual-mode — when editing, all fields are plain text inputs with the existing values.

### 3. Username field — dual-mode with generation
When adding a new credential (not editing):
- Default value: identity's email address
- Space cycles through options:
  1. Email: `jane@domain.com` (the identity's email)
  2. Name-based: `jane.doe`
  3. Name-based: `jdoe`
  4. Name-based: `janedoe`
  5. Random handle: `adjective + noun + 4 digits` (e.g. `swiftfox4821`) — regenerated each cycle
- The name-based options use the identity's first/last name, lowercased
- Cycling wraps around (after random handle, back to email)
- Typing any character enters edit mode

### 4. Password field — dual-mode with generation
When adding a new credential (not editing):
- Default value: a generated 20-character password via `zcrypto.GeneratePassword(20)`
- Space regenerates a new password (each press = fresh password)
- The password field should be visible (not masked) while in cycle mode so the user can see what was generated
- When the user switches to edit mode (types a character), the field masks as usual (`EchoPassword`)
- Esc returns to cycle mode and shows the current generated password (unmasked)
- Remove the old `ctrl+g` generate password shortcut — space replaces it

### 5. Pass identity info to form
- Change `addCredentialMsg` to carry `identity identity.Identity` instead of just `identityID string`
- Update `credentialFormModel` to store `identity identity.Identity` (replaces `identityID string`)
- Update `newCredentialFormModel` signature: `newCredentialFormModel(id identity.Identity, existing *credential.Credential)`
- In the root model's `addCredentialMsg` handler, pass `m.detail.identity` instead of `msg.identityID`
- For `editCredentialMsg`, look up the identity from `m.detail.identity` as well
- Update `credential_list.go` to emit the identity in `addCredentialMsg`
- Update `submit()` to use `m.identity.ID` for `c.IdentityID`

### 6. Update tests
- Update all `newCredentialFormModel("abc12345", ...)` calls to pass `testIdentity()` instead
- Add `TestCredentialFormIdentityHeader` — verify the form view contains the identity name and email
- Add `TestCredentialFormUsernameCycle` — verify space cycles through email → name variations → random handle
- Add `TestCredentialFormPasswordGenerated` — verify password is pre-filled and non-empty on new form
- Add `TestCredentialFormPasswordCycle` — verify space regenerates password
- Add `TestCredentialFormEditModeNoGeneration` — verify editing an existing credential shows plain text inputs with no cycling
- Add `TestCredentialFormUsernameEditMode` — verify typing a character switches to edit mode
- Add `TestCredentialFormUsernameEscCycle` — verify Esc returns to cycle mode
- Remove or update `TestCredentialFormGeneratePassword` (ctrl+g test) since ctrl+g is removed

## Target Repo
zarlcorp/zburn

## Agent Role
backend

## Files to Modify
- internal/tui/credential_form.go (dual-mode fields, identity header, username/password generation)
- internal/tui/credential_list.go (update addCredentialMsg to carry identity)
- internal/tui/tui.go (pass identity to form constructor)
- internal/tui/credential_test.go (update form constructor calls, add new tests)

## Notes
- The `adjectives` and `nouns` slices are in `internal/identity/data.go` — the random handle generation in the form can either import and reuse them, or use a simpler inline approach (e.g. `hex(random 4 bytes)` as a fallback). Reusing the identity package's `pick()` and word lists is preferred for consistency with email generation.
- The identity generator already has the `pick(adjectives) + pick(nouns) + digits` pattern in `Email()` case 7 — the form can call a similar function
- `zcrypto.GeneratePassword(20)` is already imported and used in the current form — just call it on init instead of on ctrl+g
- The dual-mode pattern is contained within `credentialFormModel` — no need for a separate reusable component yet. If we need it elsewhere later, we can extract it then.
- When in cycle mode, the `textinput.Model`'s value is set programmatically via `SetValue()`. The input is focused but typing a character triggers the switch to edit mode.
