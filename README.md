# zarlcorp/core

Shared Go packages for zarlcorp privacy tools.

## Packages

| Package | Description | Status |
|---------|-------------|--------|
| zapp | Application lifecycle toolkit | ready |
| zcache | Generic caching with multiple backends | ready |
| zcrypto | Encryption primitives | stub |
| zfilesystem | Filesystem abstraction | ready |
| zoptions | Generic functional options | ready |
| zstyle | TUI visual identity — colors, styles, keybindings | ready |
| zsync | Thread-safe data structures | ready |

## Install

Import individual packages as needed:

```bash
go get github.com/zarlcorp/core/pkg/zapp
go get github.com/zarlcorp/core/pkg/zstyle
go get github.com/zarlcorp/core/pkg/zcache
go get github.com/zarlcorp/core/pkg/zsync
```

## Quick Example

```go
package main

import (
    "context"
    "fmt"

    "github.com/zarlcorp/core/pkg/zapp"
    "github.com/zarlcorp/core/pkg/zstyle"
)

func main() {
    app := zapp.New()
    ctx, cancel := zapp.SignalContext(context.Background())
    defer cancel()

    fmt.Println(zstyle.Title.Render("hello from zarlcorp"))

    <-ctx.Done()
    app.Close()
}
```

## Learn More

- [MANIFESTO.md](./MANIFESTO.md) — Philosophy and architecture
- [LICENSE](./LICENSE) — MIT License

---

MIT License
