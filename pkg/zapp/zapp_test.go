package zapp_test

import (
	"context"
	"errors"
	"io"
	"sync"
	"testing"

	"github.com/zarlcorp/core/pkg/zapp"
)

func TestNew(t *testing.T) {
	app := zapp.New()
	if app == nil {
		t.Fatal("New returned nil")
	}
}

func TestCloseLIFO(t *testing.T) {
	var order []int

	app := zapp.New()
	app.Track(&closer{id: 1, order: &order})
	app.Track(&closer{id: 2, order: &order})
	app.Track(&closer{id: 3, order: &order})

	if err := app.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	want := []int{3, 2, 1}
	if len(order) != len(want) {
		t.Fatalf("close order = %v, want %v", order, want)
	}
	for i := range want {
		if order[i] != want[i] {
			t.Fatalf("close order = %v, want %v", order, want)
		}
	}
}

func TestCloseErrors(t *testing.T) {
	errA := errors.New("a broke")
	errB := errors.New("b broke")

	app := zapp.New()
	app.Track(&closer{order: new([]int), err: errA})
	app.Track(&closer{order: new([]int)})
	app.Track(&closer{order: new([]int), err: errB})

	err := app.Close()
	if err == nil {
		t.Fatal("Close returned nil, want error")
	}

	if !errors.Is(err, errA) {
		t.Errorf("error should contain errA")
	}
	if !errors.Is(err, errB) {
		t.Errorf("error should contain errB")
	}
}

func TestCloseNoResources(t *testing.T) {
	app := zapp.New()
	if err := app.Close(); err != nil {
		t.Fatalf("Close with no resources: %v", err)
	}
}

func TestCloseIdempotent(t *testing.T) {
	var count int
	app := zapp.New()
	app.Track(zapp.CloserFunc(func() error {
		count++
		return nil
	}))

	_ = app.Close()
	_ = app.Close()
	_ = app.Close()

	if count != 1 {
		t.Fatalf("closer called %d times, want 1", count)
	}
}

func TestCloseIdempotentError(t *testing.T) {
	want := errors.New("boom")
	app := zapp.New()
	app.Track(zapp.CloserFunc(func() error { return want }))

	err1 := app.Close()
	err2 := app.Close()

	if !errors.Is(err1, want) {
		t.Fatalf("first Close: got %v, want %v", err1, want)
	}
	if !errors.Is(err2, want) {
		t.Fatalf("second Close: got %v, want %v", err2, want)
	}
}

func TestTrackConcurrent(t *testing.T) {
	app := zapp.New()
	var wg sync.WaitGroup
	var order []int
	var mu sync.Mutex

	for i := range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			app.Track(zapp.CloserFunc(func() error {
				mu.Lock()
				order = append(order, i)
				mu.Unlock()
				return nil
			}))
		}()
	}

	wg.Wait()

	if err := app.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	if len(order) != 100 {
		t.Fatalf("closed %d resources, want 100", len(order))
	}
}

func TestTrackBeforeClose(t *testing.T) {
	app := zapp.New()

	err := app.Track(zapp.CloserFunc(func() error { return nil }))
	if err != nil {
		t.Fatalf("Track before Close: %v", err)
	}
}

func TestTrackAfterClose(t *testing.T) {
	app := zapp.New()
	app.Close()

	err := app.Track(zapp.CloserFunc(func() error { return nil }))
	if !errors.Is(err, zapp.ErrClosed) {
		t.Fatalf("Track after Close: got %v, want %v", err, zapp.ErrClosed)
	}
}

func TestTrackAfterCloseNotCalled(t *testing.T) {
	// the closer registered after Close must not be called on a subsequent Close
	app := zapp.New()
	app.Close()

	var called bool
	app.Track(zapp.CloserFunc(func() error {
		called = true
		return nil
	}))

	// second Close is a no-op due to sync.Once
	app.Close()

	if called {
		t.Fatal("closer registered after Close should not be called")
	}
}

func TestCloserFunc(t *testing.T) {
	var called bool
	var c io.Closer = zapp.CloserFunc(func() error {
		called = true
		return nil
	})

	if err := c.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	if !called {
		t.Fatal("CloserFunc was not called")
	}
}

func TestCloserFuncError(t *testing.T) {
	want := errors.New("cleanup failed")
	c := zapp.CloserFunc(func() error { return want })

	if got := c.Close(); !errors.Is(got, want) {
		t.Fatalf("Close = %v, want %v", got, want)
	}
}

func TestSignalContextCancel(t *testing.T) {
	ctx, cancel := zapp.SignalContext(context.Background())
	defer cancel()

	// calling cancel should cancel the context
	cancel()

	select {
	case <-ctx.Done():
		// expected
	default:
		t.Fatal("context not canceled after cancel()")
	}
}

func TestSignalContextInheritsParent(t *testing.T) {
	parent, parentCancel := context.WithCancel(context.Background())
	ctx, cancel := zapp.SignalContext(parent)
	defer cancel()

	// canceling parent should cancel the signal context
	parentCancel()

	select {
	case <-ctx.Done():
		// expected
	default:
		t.Fatal("context not canceled when parent canceled")
	}
}

// closer records whether Close was called and in what order.
type closer struct {
	id    int
	order *[]int
	err   error
}

func (c *closer) Close() error {
	*c.order = append(*c.order, c.id)
	return c.err
}
