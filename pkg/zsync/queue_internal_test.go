package zsync

import "testing"

// verifies that popped elements are zeroed in the backing array so the queue
// doesn't retain references to values that have already been consumed.
func TestZQueue_popZeroesBackingArray(t *testing.T) {
	q := NewZQueue[*string]()

	s := "hello"
	q.Push(&s)

	// grab reference to backing array before pop
	backing := q.items[:1]

	got, err := q.Pop()
	if err != nil {
		t.Fatalf("Pop() error = %v", err)
	}
	if *got != "hello" {
		t.Fatalf("Pop() = %v, want hello", *got)
	}

	// the element at index 0 of the original backing array should be nil
	if backing[0] != nil {
		t.Errorf("backing array still holds reference after Pop")
	}
}

func TestZQueue_tryPopZeroesBackingArray(t *testing.T) {
	q := NewZQueue[*string]()

	s := "world"
	q.Push(&s)

	backing := q.items[:1]

	got, err := q.TryPop()
	if err != nil {
		t.Fatalf("TryPop() error = %v", err)
	}
	if *got != "world" {
		t.Fatalf("TryPop() = %v, want world", *got)
	}

	if backing[0] != nil {
		t.Errorf("backing array still holds reference after TryPop")
	}
}

func TestZQueue_popContextZeroesBackingArray(t *testing.T) {
	q := NewZQueue[*string]()

	s := "ctx"
	q.Push(&s)

	backing := q.items[:1]

	got, err := q.PopContext(t.Context())
	if err != nil {
		t.Fatalf("PopContext() error = %v", err)
	}
	if *got != "ctx" {
		t.Fatalf("PopContext() = %v, want ctx", *got)
	}

	if backing[0] != nil {
		t.Errorf("backing array still holds reference after PopContext")
	}
}
