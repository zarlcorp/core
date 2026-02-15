package zsync

import (
	"context"
	"io"
	"sync"
)

// ZQueue is a thread-safe generic FIFO queue implementation.
// It provides blocking operations for producer/consumer scenarios with context cancellation support.
type ZQueue[T any] struct {
	mu     sync.Mutex
	cond   *sync.Cond
	items  []T
	closed bool
}

// ensure ZQueue implements io.Closer
var _ io.Closer = (*ZQueue[any])(nil)

// NewZQueue creates a new thread-safe FIFO queue.
func NewZQueue[T any]() *ZQueue[T] {
	q := &ZQueue[T]{
		items: make([]T, 0),
	}
	q.cond = sync.NewCond(&q.mu)
	return q
}

// Push adds an item to the back of the queue.
// Returns ErrQueueClosed if the queue has been closed.
func (q *ZQueue[T]) Push(item T) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		return ErrQueueClosed
	}

	q.items = append(q.items, item)
	q.cond.Signal()
	return nil
}

// Pop removes and returns an item from the front of the queue.
// Returns ErrQueueClosed if the queue is empty and closed.
func (q *ZQueue[T]) Pop() (T, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for len(q.items) == 0 && !q.closed {
		q.cond.Wait()
	}

	if len(q.items) == 0 && q.closed {
		var zero T
		return zero, ErrQueueClosed
	}

	item := q.items[0]
	var zero T
	q.items[0] = zero
	q.items = q.items[1:]
	return item, nil
}

// PopContext removes and returns an item from the front of the queue with context cancellation.
// Returns ErrQueueClosed if the queue is empty and closed, or ErrCanceled if the context is canceled.
func (q *ZQueue[T]) PopContext(ctx context.Context) (T, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for len(q.items) == 0 && !q.closed {
		// check context before waiting
		select {
		case <-ctx.Done():
			var zero T
			return zero, ErrCanceled
		default:
		}

		// use goroutine to wait on context while waiting on condition
		contextDone := make(chan struct{})
		go func() {
			select {
			case <-ctx.Done():
				q.cond.Broadcast() // wake up waiting pop operations
			case <-contextDone:
			}
		}()

		q.cond.Wait()
		close(contextDone) // stop context watcher goroutine

		// check if context canceled while waiting
		select {
		case <-ctx.Done():
			var zero T
			return zero, ErrCanceled
		default:
		}
	}

	if len(q.items) == 0 && q.closed {
		var zero T
		return zero, ErrQueueClosed
	}

	item := q.items[0]
	var zero T
	q.items[0] = zero
	q.items = q.items[1:]
	return item, nil
}

// TryPop attempts to remove and return an item from the front of the queue without blocking.
// Returns ErrQueueEmpty if the queue is empty.
func (q *ZQueue[T]) TryPop() (T, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.items) == 0 {
		var zero T
		return zero, ErrQueueEmpty
	}

	item := q.items[0]
	var zero T
	q.items[0] = zero
	q.items = q.items[1:]
	return item, nil
}

// Len returns the current number of items in the queue.
func (q *ZQueue[T]) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items)
}

// Close closes the queue, waking up any blocked Pop operations.
// After closing, Push operations will return ErrQueueClosed.
func (q *ZQueue[T]) Close() error {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.closed = true
	q.cond.Broadcast()
	return nil
}

// IsClosed returns true if the queue has been closed.
func (q *ZQueue[T]) IsClosed() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.closed
}
