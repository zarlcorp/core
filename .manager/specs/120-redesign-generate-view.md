# 120: Redesign generate identity view

## Objective
Redesign the generate identity view with a cleaner layout: email first, grouped fields, no internal ID, and a placeholder avatar image rendered via the Kitty graphics protocol.

## Context
The current generate view shows fields in a flat list starting with the internal ID. The user wants email as the primary field, grouped layout, and a visual identity card feel with a placeholder avatar.

The `github.com/charmbracelet/x/ansi/kitty` package (already a transitive dependency via charmbracelet ecosystem) provides `EncodeGraphics()` for rendering images in terminals that support the Kitty graphics protocol. Ghostty supports this.

## Requirements

### 1. Remove ID from field list
- Remove `{"id", id.ID}` from `identityFields()` in `generate.go`
- The ID is an internal UUID — users don't need to see it
- The ID still exists on the `identity.Identity` struct for storage/operations

### 2. Reorder fields — email first
Update `identityFields()` to return fields in this order:
```
email
name
phone
street
city, state zip   (combined into one line)
dob
```
- Combine city/state/zip into a single "address" field: `"Portland, OR 97201"`
- This reduces 4 address fields to 2 (street + combined city/state/zip)
- Total fields: 6 (was 9)

### 3. Grouped visual layout
In `View()`, add visual grouping with blank lines:
```
  [avatar]  email    jane.doe@domain.com  [domain.com]  space to cycle
            name     Jane Doe
            phone    (555) 123-4567

            street   123 Oak Ave
            address  Portland, OR 97201

            dob      1990-06-15
```
- Avatar image displayed to the left of the first few fields (see req 4)
- Blank line between contact info and address block
- Blank line between address block and DOB

### 4. Placeholder avatar image
- Embed a generic avatar PNG in the binary using `//go:embed`
- Create `internal/tui/avatar.go` with the embedded image
- Create `internal/tui/avatar.png` — a simple generic silhouette avatar (64x64 or 48x48 pixels)
- In `View()`, render the avatar to the left of the identity fields using `kitty.EncodeGraphics()`
- The avatar is displayed once when the view renders, positioned at the top-left of the identity card
- Fall back gracefully if the terminal doesn't support Kitty graphics — just skip the image, show fields only

### 5. Kitty graphics integration
- Import `github.com/charmbracelet/x/ansi/kitty`
- Use `kitty.EncodeGraphics(buf, img, &kitty.Options{Action: kitty.TransmitAndPut, Format: kitty.PNG, Chunk: true})` to render
- Write the escape sequence into the view string
- The image is a static placeholder — no per-identity generation yet

### 6. Update tests
- Update `TestIdentityFields` — now 6 fields, no "id" field, email is first
- Update `TestGenerateViewShowsFields` — check for email, name, phone, combined address
- Add test for avatar embed (image bytes are non-empty)
- Update any tests that reference field indices

## Target Repo
zarlcorp/zburn

## Agent Role
backend

## Files to Modify
- internal/tui/generate.go (field order, layout, avatar rendering)
- internal/tui/avatar.go (new — embed directive and helper)
- internal/tui/avatar.png (new — placeholder image)
- internal/tui/tui_test.go (update field tests)

## Notes
- The avatar is a PLACEHOLDER — we'll add generated avatars (thispersondoesnotexist.com or similar) in a future spec
- The `charmbracelet/x/ansi/kitty` package is already available as a transitive dep but may need to be imported directly — run `go mod tidy` after adding the import
- Terminal detection for Kitty support is not required for v1 — assume the terminal supports it (Ghostty does). If it doesn't, the escape sequences will be ignored or show as garbage, which is acceptable for now
- The `allFieldsText()` method used by "copy all" should also reflect the new field order
