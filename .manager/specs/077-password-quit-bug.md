# 077: Fix password view quitting on 'q' keypress

## Objective
Typing 'q' as part of a master password in zburn's TUI quits the application instead of entering the character into the password field.

## Context
The password model (`internal/tui/password.go`) checks for `zstyle.KeyQuit` (which matches `"q"` and `"ctrl+c"`) before passing the keystroke to the textinput. Line 53:

```go
if key.Matches(msg, zstyle.KeyQuit) && msg.Type != tea.KeyEnter {
    return m, tea.Quit
}
```

The guard `msg.Type != tea.KeyEnter` was meant to prevent enter from quitting, but it doesn't prevent regular character `q` from quitting. The textinput never sees the keystroke because the quit fires first.

## Requirements

### Fix in password.go
Remove the `zstyle.KeyQuit` check from the password model's Update method. The password view has a focused textinput at all times — `q` is valid password input and must reach the textinput.

Replace lines 52-55:
```go
case tea.KeyMsg:
    if key.Matches(msg, zstyle.KeyQuit) && msg.Type != tea.KeyEnter {
        return m, tea.Quit
    }
```

With only `ctrl+c` as quit:
```go
case tea.KeyMsg:
    if msg.Type == tea.KeyCtrlC {
        return m, tea.Quit
    }
```

This way:
- `q` → goes to textinput (password character)
- `ctrl+c` → quits (universal escape hatch)
- `enter` → handled by the existing `zstyle.KeyEnter` check below

### Test
Ensure existing tests still pass. If there are no tests for the password view, the agent should write a basic test that verifies a `tea.KeyMsg` with rune 'q' does NOT produce `tea.Quit`.

## Target Repo
zarlcorp/zburn

## Agent Role
backend

## Files to Modify
- internal/tui/password.go

## Notes
One-line fix. The other views (menu, generate, list, detail) correctly use `q` to quit because they don't have focused text inputs — no changes needed there.
