# 122: Redesign detail view

## Objective
Redesign the identity detail view to mirror the generate view layout: email first, grouped fields, avatar, and name as title instead of UUID.

## Context
Spec 120 redesigned the generate view with email-first field ordering, grouped sections (contact / address / personal), and a placeholder avatar via kitty graphics. The detail view still uses the old layout: UUID in the title, RFC3339 timestamp, flat field list.

The detail view should look and feel like the generate view — same `identityFields()`, same grouping, same avatar — but with "view" actions (copy, credentials, burn) instead of "generate" actions (save, new, cycle domain).

## Requirements

### 1. Replace UUID title with name
- Current: `zstyle.Title.Render("identity " + m.identity.ID)`
- New: `zstyle.Title.Render(m.identity.FirstName + " " + m.identity.LastName)`
- Remove the RFC3339 timestamp display entirely

### 2. Use shared identityFields()
- The detail view already calls `identityFields(id)` from generate.go
- After spec 120, this returns: email, name, phone, street, address (combined), dob
- No changes needed here — just verify the detail view uses the same function

### 3. Add grouped layout with section breaks
- Import and use the same `sectionBreaks` map from generate.go
- Add blank lines between contact (email/name/phone), address (street/address), and personal (dob) sections
- Match the exact spacing pattern from generate's `View()`

### 4. Add avatar
- Call `renderAvatar()` (from generate.go) in `newDetailModel()` or in `View()`
- Store the rendered string on `detailModel` (add `avatar string` field)
- Render it above the fields, same position as in generate view

### 5. Simplify credential section
- Current: `credentials (3)` subtitle + `w to view` hint
- New: `(3) credentials  w to view` — inline, no subtitle styling
- When count is 0: `no credentials  a to add` (or just omit the section)
- Use `zstyle.MutedText` for the whole line

### 6. Update tests
- `TestDetailViewShowsFields` — remove `abc12345` from checks (UUID no longer in title), verify name is in title
- `TestDetailViewShowsCredentialCount` — update expected format
- `TestDetailViewCredentialCountZero` — update expected format
- `TestDetailHelpShowsCredentials` — may need format update
- Add test verifying section breaks appear in view output

## Target Repo
zarlcorp/zburn

## Agent Role
backend

## Files to Modify
- internal/tui/detail.go (layout, title, avatar, grouped sections, credential display)
- internal/tui/tui_test.go (update detail tests)
- internal/tui/credential_test.go (update credential count display tests)

## Notes
- `identityFields()`, `sectionBreaks`, and `renderAvatar()` are all defined in generate.go and already package-scoped — detail.go can use them directly
- The `allFieldsText()` method on detailModel should continue to work unchanged since it iterates `m.fields`
- The detail view keeps its own action bar: `enter copy field  c copy all  w credentials  d burn  esc back  q quit`
- This spec depends on spec 121 being merged first (list redesign) to avoid merge conflicts in tui_test.go, but can be developed in parallel on a separate branch
