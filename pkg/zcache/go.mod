module github.com/zarlcorp/core/pkg/zcache

go 1.24.4

require (
	github.com/redis/go-redis/v9 v9.12.1
	github.com/zarlcorp/core/pkg/zfilesystem v0.0.0
	github.com/zarlcorp/core/pkg/zoptions v0.0.0
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/zarlcorp/core/pkg/zsync v0.0.0 // indirect
)

replace (
	github.com/zarlcorp/core/pkg/zfilesystem => ../zfilesystem
	github.com/zarlcorp/core/pkg/zoptions => ../zoptions
	github.com/zarlcorp/core/pkg/zsync => ../zsync
)
