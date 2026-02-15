// Package zapp provides application lifecycle management for zarlcorp tools.
//
// zapp is a toolkit, not a framework. The consumer owns main and wires
// things together explicitly. zapp handles resource tracking with ordered
// cleanup, signal-based context cancellation, and functional options.
//
// # Usage
//
//	func main() {
//		app := zapp.New()
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
	"os/signal"
	"sync"
	"syscall"

	"github.com/zarlcorp/core/pkg/zoptions"
)

// ErrClosed is returned by Track when the app has already been closed.
var ErrClosed = errors.New("app closed")

// Option configures an App.
type Option = zoptions.Option[App]

// App tracks resources and tears them down in LIFO order on close.
type App struct {
	mu      sync.Mutex
	once    sync.Once
	closed  bool
	closers []io.Closer
	err     error
}

// New creates an App with the given options.
func New(opts ...Option) *App {
	a := &App{}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

// Track registers a closer for cleanup. Returns ErrClosed if the app
// has already been closed. Safe for concurrent use.
func (a *App) Track(c io.Closer) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.closed {
		return ErrClosed
	}

	a.closers = append(a.closers, c)
	return nil
}

// Close tears down all tracked resources in LIFO order. Returns a joined
// error if any closers fail. Safe to call multiple times — subsequent
// calls return the same result without re-closing.
func (a *App) Close() error {
	a.once.Do(func() {
		a.mu.Lock()
		a.closed = true
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

// SignalContext returns a context that is canceled when SIGINT or SIGTERM
// is received, or when the returned cancel func is called.
func SignalContext(parent context.Context) (context.Context, context.CancelFunc) {
	return signal.NotifyContext(parent, syscall.SIGINT, syscall.SIGTERM)
}

// CloserFunc adapts a func() error into an io.Closer.
type CloserFunc func() error

// Close calls the underlying function.
func (f CloserFunc) Close() error { return f() }
