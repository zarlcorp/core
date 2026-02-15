package zsync

import (
	"cmp"
	"slices"
)

// ZSet is a thread-safe generic set implementation.
// It provides concurrent access to a set of unique values built on top of ZMap.
type ZSet[T comparable] struct {
	m *ZMap[T, struct{}]
}

// NewZSet creates a new thread-safe set with the specified value type.
func NewZSet[T comparable]() *ZSet[T] {
	return &ZSet[T]{
		m: NewZMap[T, struct{}](),
	}
}

// Add inserts a value into the set.
// If the value already exists, this operation has no effect.
func (s *ZSet[T]) Add(value T) {
	s.m.Set(value, struct{}{})
}

// Contains checks if a value exists in the set.
// Returns true if the value is present, false otherwise.
func (s *ZSet[T]) Contains(value T) bool {
	_, err := s.m.Get(value)
	return err == nil
}

// Remove deletes a value from the set.
// Returns true if the value existed and was removed, false otherwise.
func (s *ZSet[T]) Remove(value T) bool {
	return s.m.Delete(value)
}

// Len returns the number of unique values in the set.
func (s *ZSet[T]) Len() int {
	return s.m.Len()
}

// Values returns a slice containing all values in the set.
// The order of values is not guaranteed.
func (s *ZSet[T]) Values() []T {
	return s.m.Keys()
}

// Clear removes all values from the set.
func (s *ZSet[T]) Clear() {
	s.m.Clear()
}

// Ordered returns a slice containing all values in the set sorted in ascending order.
// Only works with types that implement cmp.Ordered (strings, numbers, etc).
//
// Example:
//
//	values := Ordered(s) // automatically sorts strings, ints, etc.
func Ordered[T cmp.Ordered](s *ZSet[T]) []T {
	values := s.Values()
	slices.Sort(values)
	return values
}

// Ordered returns a slice containing all values in the set sorted using the provided
// comparison function. Use the package-level Ordered for cmp.Ordered types.
func (s *ZSet[T]) Ordered(compare func(a, b T) int) []T {
	values := s.m.Keys()
	slices.SortFunc(values, compare)
	return values
}
