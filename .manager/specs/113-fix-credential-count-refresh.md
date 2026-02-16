# 113: Fix stale credential count on identity detail view

## Objective
Fix two issues in the zburn TUI: (1) update the zstyle dependency to get the KeyBack fix from spec 112, and (2) refresh the credential count on the identity detail view after credential operations.

## Context
After adding or deleting a credential, the identity detail view still shows the old credential count. The count is set once in `handleViewIdentity()` and never refreshed when navigating back from credential views. The integration tests pass because they explicitly send `viewIdentityMsg` which re-triggers `handleViewIdentity()`, but real UI navigation uses `navigateMsg{view: viewDetail}` which just switches the active view without refreshing.

## Requirements

### 1. Bump core dependency
Run `go get github.com/zarlcorp/core@latest` to pull in the KeyBack fix from spec 112.

### 2. Refresh credential count on navigation back to detail
In `internal/tui/tui.go`, when handling `navigateMsg` with `view: viewDetail`, refresh the credential count on the existing detail model:

```go
case viewDetail:
    // refresh credential count
    if m.credentials != nil {
        count, err := m.countCredentials(m.detail.identity.ID)
        if err == nil {
            m.detail.credentialCount = count
        }
    }
    m.active = viewDetail
```

The key change is adding the count refresh before setting `m.active = viewDetail`. The detail model already exists (it was created when the user first navigated to the identity), so we just update its credential count.

### 3. Run tests
`go test ./...` must pass.

## Target Repo
zarlcorp/zburn

## Agent Role
backend

## Files to Modify
- internal/tui/tui.go
- go.mod / go.sum (dependency bump)

## Depends On
- 112 (KeyBack fix must be merged and tagged in core first)

## Notes
The `navigateMsg` handler is around line 288-307 in tui.go. Look for the `case navigateMsg:` switch and find the `viewDetail` case within it.
