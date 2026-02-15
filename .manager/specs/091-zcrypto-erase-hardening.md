# 091: zcrypto Erase hardening — compiler-resistant zeroing

## Objective
Replace the current zeroing loop + runtime.KeepAlive approach with a more robust method that resists compiler dead-store elimination.

## Context
Review of pkg/ shared libraries noted that the current Erase implementation uses a for-loop to zero bytes followed by `runtime.KeepAlive`. The compiler can optimize away the zeroing loop as dead stores (values are never read afterward). KeepAlive prevents GC of the slice but doesn't prevent the write elimination.

## Requirements
- Replace the manual loop with a compiler-resistant approach. Options (in order of preference):
  1. Use `clear(b)` (Go 1.21+ built-in) — the compiler intrinsic is less likely to be optimized away, and pair with `runtime.KeepAlive`
  2. Use a package-level `var eraser func([]byte)` that is set to the zeroing function at init time — the indirection prevents the compiler from proving the writes are dead
  3. Use `unsafe.Pointer` + `memclr` pattern
- Keep the honest doc comment about best-effort guarantees (GC may copy)
- All existing tests must pass
- Add a test that verifies the slice is actually zeroed after Erase (it already likely exists, but confirm)

## Target Repo
zarlcorp/core

## Agent Role
backend

## Files to Modify
- pkg/zcrypto/erase.go

## Notes
Independent — no dependency on other specs. Very small scope.
