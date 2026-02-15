// Package zapp provides application lifecycle management for zarlcorp tools.
//
// zapp is a toolkit, not a framework. The consumer owns main and wires
// things together explicitly. zapp handles resource tracking with ordered
// cleanup, signal-based context cancellation, and functional options.
//
// # Usage
//
//	func main() {
//		app := zapp.New(zapp.WithName("myservice"))
//
//		ctx, cancel := zapp.SignalContext(context.Background())
//		defer cancel()
//
//		db := openDB()
//		app.Track(db)
//
//		srv := startServer(ctx, db)
//		app.Track(zapp.CloserFunc(func() error {
//			return srv.Shutdown(context.Background())
//		}))
//
//		<-ctx.Done()
//
//		if err := app.Close(); err != nil {
//			slog.Error("shutdown", "err", err)
//			os.Exit(1)
//		}
//	}
package zapp

import (
	"context"
	"errors"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/zarlcorp/core/pkg/zoptions"
)

// Option configures an App.
type Option = zoptions.Option[App]

// App tracks resources and tears them down in LIFO order on close.
type App struct {
	mu      sync.Mutex
	once    sync.Once
	name    string
	closers []io.Closer
	err     error
}

// New creates an App with the binary basename as the default name.
// Use WithName to override.
func New(opts ...Option) *App {
	a := &App{
		name: filepath.Base(os.Args[0]),
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

// Track registers a closer for cleanup. Safe for concurrent use.
func (a *App) Track(c io.Closer) {
	a.mu.Lock()
	a.closers = append(a.closers, c)
	a.mu.Unlock()
}

// Close tears down all tracked resources in LIFO order. Returns a joined
// error if any closers fail. Safe to call multiple times — subsequent
// calls return the same result without re-closing.
func (a *App) Close() error {
	a.once.Do(func() {
		a.mu.Lock()
		closers := a.closers
		a.mu.Unlock()

		var errs []error
		for i := len(closers) - 1; i >= 0; i-- {
			if err := closers[i].Close(); err != nil {
				errs = append(errs, err)
			}
		}
		a.err = errors.Join(errs...)
	})
	return a.err
}

// SignalContext returns a context that is cancelled when SIGINT or SIGTERM
// is received, or when the returned cancel func is called.
func SignalContext(parent context.Context) (context.Context, context.CancelFunc) {
	return signal.NotifyContext(parent, syscall.SIGINT, syscall.SIGTERM)
}

// CloserFunc adapts a func() error into an io.Closer.
type CloserFunc func() error

// Close calls the underlying function.
func (f CloserFunc) Close() error { return f() }
