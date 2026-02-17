# 123: Fix avatar rendering

## Objective
Fix the kitty graphics avatar not rendering by adding the missing `Transmission: kitty.Direct` option.

## Context
The `renderAvatar()` function in generate.go uses `kitty.EncodeGraphics()` but omits the `Transmission` option. The default value is `0`, which matches none of the switch cases in `EncodeGraphics` (`Direct = 'd'`, `File = 'f'`, etc.), so the image data is never encoded. The resulting escape sequence is empty — just `_Gf=100,a=T\` with no payload.

## Requirements

### 1. Add Transmission option
In `renderAvatar()` in generate.go, add `Transmission: kitty.Direct` to the options:
```go
opts := &kitty.Options{
    Action:       kitty.TransmitAndPut,
    Transmission: kitty.Direct,
    Format:       kitty.PNG,
    Chunk:        true,
}
```

### 2. Update test
Add a test that verifies `renderAvatar()` returns a non-empty string containing base64 image data (not just the empty `_Gf=100,a=T\` control sequence). Check that the output length is > 100 bytes (the real payload is ~400 bytes).

## Target Repo
zarlcorp/zburn

## Agent Role
backend

## Files to Modify
- internal/tui/generate.go (add Transmission option)
- internal/tui/tui_test.go (add renderAvatar test)

## Notes
- This is a one-line fix plus a test
- The avatar is rendered in both generate.go and detail.go via `renderAvatar()`, so fixing it once fixes both views
- Verified with POC: adding `Transmission: kitty.Direct` produces 419 bytes of escape sequence with actual base64 PNG data
