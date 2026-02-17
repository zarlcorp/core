# 121: Redesign identity list

## Objective
Redesign the identity list view to show only essential information: name, email, and credential count. Remove the ID column, created date, and table header.

## Context
The current list view shows `id | name | email | created` with a table header row. The ID is an internal UUID that means nothing to users. The created date provides no actionable information. Users also have no way to tell which identities have credentials without drilling into each one.

Spec 120 established the email-first, minimal design direction. The list should follow the same philosophy: show only what matters.

## Requirements

### 1. Remove table header
- Delete the `header` line that renders column headers (`id`, `name`, `email`, `created`)
- The list should be clean rows with no header

### 2. Remove ID and created date columns
- Remove `id.ID` from the row format
- Remove `id.CreatedAt.Format(...)` from the row format
- These fields still exist on the struct for internal use

### 3. Add credential count badge
- Show credential count as `(N)` after the email, only when N > 0
- When count is 0, show nothing — no `(0)` badge
- Format: `Jane Doe  jane@domain.com  (3)`
- Use `zstyle.MutedText` for the badge

### 4. Load credential counts in bulk
- In `loadList()` (tui.go), after loading identities, load all credentials once via `m.credentials.List()`
- Build a `map[string]int` of identity ID → credential count
- Pass the counts to `newListModel()` — add a `credCounts map[string]int` field to `listModel`
- This avoids N+1 queries (one `List()` call instead of one per identity)

### 5. Clean row format
Each row should render as:
```
  > Jane Doe  jane@domain.com  (3)
    Bob Smith  bob@other.com
```
- Name and email separated by spacing (use `fmt.Sprintf` with widths or lipgloss)
- Selected row gets `>` prefix with `zstyle.Highlight`
- Non-selected rows get matching indent
- `truncate()` still applies to prevent overflow

### 6. Update tests
- `TestListViewShowsIdentities` — remove check for `abc12345` (ID), add check for credential count badge
- `TestListViewEmpty` — should still work unchanged
- `TestListNavigation` — should still work unchanged
- Add `TestListViewCredentialCount` — verify `(3)` appears when count > 0
- Add `TestListViewCredentialCountZero` — verify no `(0)` appears when count is 0

## Target Repo
zarlcorp/zburn

## Agent Role
backend

## Files to Modify
- internal/tui/list.go (row format, remove header, add credential counts)
- internal/tui/tui.go (loadList — bulk load credential counts)
- internal/tui/tui_test.go (update list view tests)

## Notes
- The `truncate()` helper in list.go should be kept — it's still useful for long names/emails
- The credential count loading reuses the existing `m.credentials.List()` + filter pattern already used in `countCredentials()` and `loadCredentialList()` — just applied in bulk
- Keep the sort by CreatedAt descending in `loadList()` even though we no longer display the date
