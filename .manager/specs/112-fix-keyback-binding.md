# 112: Remove backspace from KeyBack binding

## Objective
Fix the `KeyBack` key binding in zstyle so that backspace no longer triggers "navigate back." Only Escape should navigate back.

## Context
`zstyle.KeyBack` is defined as `key.NewBinding(key.WithKeys("esc", "backspace"))`. Every TUI view that checks `key.Matches(msg, zstyle.KeyBack)` treats both Escape and Backspace as "go back." This means pressing backspace in any text input field navigates away instead of deleting the previous character — making data entry unusable.

## Requirements

### Fix KeyBack binding
In `pkg/zstyle/keys.go`, change `KeyBack` to only bind "esc":

```go
// before
KeyBack = key.NewBinding(key.WithKeys("esc", "backspace"), key.WithHelp("esc", "back"))

// after
KeyBack = key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back"))
```

That's the entire change. Backspace will no longer match `KeyBack`, so text inputs in downstream TUIs will properly consume it for character deletion.

### Run tests
`go test ./...` must pass.

## Target Repo
zarlcorp/core

## Agent Role
backend

## Files to Modify
- pkg/zstyle/keys.go

## Notes
This is a one-line fix but it affects every downstream TUI (zburn, zvault, zshield). The behavioral change is correct: Escape = back, Backspace = delete character. No TUI should use backspace for navigation.
