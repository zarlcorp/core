# Package Review — 031

Systematic review of all 7 shared packages against zarlcorp coding standards.

## Per-Package Findings

### 1. zoptions (foundation, no deps)

**Status:** Clean. No issues found.

- Package doc with usage example: present
- Single exported type `Option[T]` with doc comment: present
- External test package (`zoptions_test`): yes
- Table-driven tests: yes
- No error handling needed (no errors in this package)
- No dependencies beyond stdlib: correct

### 2. zsync (foundation, no deps)

**Findings:**

- `set.go:68` — Method doc on `ZSet.Ordered` was identical wording to the package-level `Ordered` function, creating confusion about which does what.

**Fixed:**

- `set.go:68` — Clarified method doc to distinguish from the package-level `Ordered` function.

**Notes (not fixed, by design):**

- Package doc is verbose but appropriate for library code per the "Library/SDK: Comprehensive with examples" strategy.
- `ZMap.Get` returns `ErrNotFound` (error) rather than `(V, bool)`. This is a deliberate API choice — noted but not changed per spec constraint.

### 3. zcache (data layer)

**Findings:**

- `file.go:271` — Error message `"filesystem write failed: %w"` uses banned "failed" prefix.
- `file.go:277` — Error message `"filesystem read failed: %w"` uses banned "failed" prefix.
- `file.go:282` — Error message `"filesystem data integrity check failed"` uses banned "failed" prefix.
- `redis.go:193` — Error message `"redis ping failed: %w"` uses banned "failed" prefix.
- `redis.go:31` — `WithClient` missing doc comment.
- `redis.go:37` — `WithPrefix` missing doc comment.
- `redis.go:43` — `WithTTL` missing doc comment.
- `redis.go:29` — `RedisOption` type alias missing doc comment.
- `cache.go:40-42` — Unnecessary `var ( )` block for single `ErrNotFound` declaration; missing doc comment.

**Fixed:**

- `file.go` (Healthy) — Rewrote error messages: `"write health check: %w"`, `"read health check: %w"`, `"health check data mismatch"`.
- `redis.go` (Healthy) — Rewrote error message: `"ping redis: %w"`.
- `redis.go` — Added doc comments to `WithClient`, `WithPrefix`, `WithTTL`, `RedisOption`.
- `cache.go` — Simplified to plain `var` declaration, added doc comment on `ErrNotFound`.

**Notes (not fixed, by design):**

- `file.go` `Clear()` and `Len()` type-assert to `*zfilesystem.MemFS` for optimization. This breaks the `FileSystem` interface abstraction but is a performance choice, not a bug. Changing it would alter behavior.
- `file.go` `NewFileCache` creates a temp directory unconditionally even when options might override the filesystem. Minor waste but fixing would change constructor semantics.

### 4. zfilesystem (data layer)

**Findings:**

- `filesystem.go:56` — Doc comment on `File` interface missing trailing period.
- `memfs.go:66` — Doc comment on `MkdirAll` missing trailing period.
- `memfs.go:71` — Doc comment on `OpenFile` missing trailing period.

**Fixed:**

- `filesystem.go:56` — Added period: `"File is a file that can be read from and written to."`
- `memfs.go:66` — Added period to MkdirAll doc.
- `memfs.go:71` — Added period to OpenFile doc.

**Notes (not fixed, by design):**

- `MemFS.ClearCacheFiles()` and `MemFS.CountCacheFiles()` are specialized methods for cache implementations that leak cache concerns into a filesystem package. However, they exist as performance optimizations and removing them would change the public API.

### 5. zstyle (presentation)

**Findings:**

- `zstyle_test.go:1` — Uses internal test package (`package zstyle`) instead of external (`package zstyle_test`). Standards require testing the public API via external test package.

**Fixed:**

- `zstyle_test.go` — Converted to external test package (`package zstyle_test`), added `zstyle` import, prefixed all references with `zstyle.`.

### 6. zcrypto (security)

**Status:** Clean. No issues found.

- Package doc with usage example: present
- All exported functions/types have doc comments: yes
- External test package (`zcrypto_test`): yes
- Error messages follow direct context pattern: yes (e.g. `"create cipher: %w"`, `"generate nonce: %w"`)
- No panics: correct
- Erase function has appropriate best-effort caveat in doc
- Table-driven tests where appropriate: yes

### 7. zapp (top layer)

**Status:** Clean. No issues found.

- Package doc with usage example: present
- All exported types/functions have doc comments: yes
- External test package (`zapp_test`): yes
- Error handling via `errors.Join`: correct
- LIFO close order: tested
- Idempotent Close via `sync.Once`: tested
- Concurrent Track safety: tested
- `CloserFunc` adapter with doc: present
- `SignalContext` wrapper: clean

**Notes (not fixed, by design):**

- `App.name` field is set by `WithName` but never readable externally. No accessor exists. Not adding one per spec constraint (no API changes).

## Cross-Package Issues

### Consistency issues found and fixed

| Issue | Packages | Status |
|-------|----------|--------|
| Error messages with "failed" prefix | zcache | Fixed |
| Missing doc comments on exported symbols | zcache/redis.go | Fixed |
| Internal test package instead of external | zstyle | Fixed |
| Doc comments missing trailing periods | zfilesystem | Fixed |

### Consistency issues not found (already consistent)

- **Options pattern**: zoptions used consistently by zapp (via type alias) and zcache (direct import). Same pattern everywhere.
- **Constructor naming**: `NewX()` consistently used across all packages.
- **Error handling style**: Direct context wrapping (e.g. `"create cipher: %w"`) used everywhere except the 4 "failed" instances in zcache, now fixed.
- **Mutex usage**: All packages use unexported `mu` field, never embedded. Correct per standards.
- **External test packages**: All packages use `_test` suffix except zstyle, now fixed.
- **Package docs with examples**: Present on all packages.
- **Interface satisfaction checks**: Used in zcache (all 3 implementations), zfilesystem (both implementations), zsync (ZQueue io.Closer). Consistent pattern.
- **Table-driven tests**: Used throughout where appropriate.
- **File organization**: Consistent pattern of main file + implementation files + single test file (or test per implementation).
- **Comment style**: Lowercase terse style used consistently for inline comments.

## Verification

All packages pass:
- `go vet ./...` — clean
- `go test -race ./...` — all pass
- No remaining "failed to" / "unable to" / "could not" error prefixes
- All exported symbols have doc comments
- All test files use external test packages

## Summary of Changes

| File | Change |
|------|--------|
| `zcache/file.go` | Rewrote 3 error messages in Healthy() to remove "failed" prefix |
| `zcache/redis.go` | Rewrote 1 error message in Healthy(), added 4 doc comments |
| `zcache/cache.go` | Simplified ErrNotFound declaration, added doc comment |
| `zfilesystem/filesystem.go` | Added period to File interface doc comment |
| `zfilesystem/memfs.go` | Added periods to MkdirAll and OpenFile doc comments |
| `zstyle/zstyle_test.go` | Converted to external test package |
| `zsync/set.go` | Clarified ZSet.Ordered method doc comment |
