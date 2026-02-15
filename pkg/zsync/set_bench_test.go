package zsync_test

import (
	"strconv"
	"testing"

	"github.com/zarlcorp/core/pkg/zsync"
)

func BenchmarkZSet_Add(b *testing.B) {
	s := zsync.NewZSet[string]()

	b.ResetTimer()
	for i := range b.N {
		s.Add("item" + strconv.Itoa(i))
	}
}

func BenchmarkZSet_Contains(b *testing.B) {
	s := zsync.NewZSet[string]()

	// warmup
	for i := range 1000 {
		s.Add("item" + strconv.Itoa(i))
	}

	b.ResetTimer()
	for i := range b.N {
		s.Contains("item" + strconv.Itoa(i%1000))
	}
}

func BenchmarkZSet_Remove(b *testing.B) {
	s := zsync.NewZSet[string]()

	// warmup
	for i := range b.N {
		s.Add("item" + strconv.Itoa(i))
	}

	b.ResetTimer()
	for i := range b.N {
		s.Remove("item" + strconv.Itoa(i))
	}
}

func BenchmarkZSet_AddContains_Mixed(b *testing.B) {
	s := zsync.NewZSet[string]()

	b.ResetTimer()
	for i := range b.N {
		if i%2 == 0 {
			s.Add("item" + strconv.Itoa(i))
		} else {
			s.Contains("item" + strconv.Itoa(i-1))
		}
	}
}

func BenchmarkZSet_Concurrent(b *testing.B) {
	s := zsync.NewZSet[string]()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			item := "item" + strconv.Itoa(i)
			switch i % 3 {
			case 0:
				s.Add(item)
			case 1:
				s.Contains(item)
			default:
				s.Remove(item)
			}
			i++
		}
	})
}

func BenchmarkZSet_Values(b *testing.B) {
	s := zsync.NewZSet[string]()

	// warmup
	for i := range 1000 {
		s.Add("item" + strconv.Itoa(i))
	}

	b.ResetTimer()
	for range b.N {
		s.Values()
	}
}

func BenchmarkZSet_Ordered(b *testing.B) {
	s := zsync.NewZSet[int]()

	// warmup
	for i := range 1000 {
		s.Add(i)
	}

	b.ResetTimer()
	for range b.N {
		zsync.Ordered(s)
	}
}

func BenchmarkZSet_Clear(b *testing.B) {
	for range b.N {
		s := zsync.NewZSet[string]()

		// warmup
		for i := range 1000 {
			s.Add("item" + strconv.Itoa(i))
		}

		b.StartTimer()
		s.Clear()
		b.StopTimer()
	}
}

func BenchmarkZSet_Len(b *testing.B) {
	s := zsync.NewZSet[string]()

	// warmup
	for i := range 1000 {
		s.Add("item" + strconv.Itoa(i))
	}

	b.ResetTimer()
	for range b.N {
		s.Len()
	}
}

// comparison with stdlib map[T]struct{} + mutex
type stdSet struct {
	m *zsync.ZMap[string, struct{}]
}

func newStdSet() *stdSet {
	return &stdSet{
		m: zsync.NewZMap[string, struct{}](),
	}
}

func (ss *stdSet) Add(item string) {
	ss.m.Set(item, struct{}{})
}

func (ss *stdSet) Contains(item string) bool {
	_, err := ss.m.Get(item)
	return err == nil
}

func BenchmarkStdSet_Add(b *testing.B) {
	s := newStdSet()

	b.ResetTimer()
	for i := range b.N {
		s.Add("item" + strconv.Itoa(i))
	}
}

func BenchmarkStdSet_Contains(b *testing.B) {
	s := newStdSet()

	// warmup
	for i := range 1000 {
		s.Add("item" + strconv.Itoa(i))
	}

	b.ResetTimer()
	for i := range b.N {
		s.Contains("item" + strconv.Itoa(i%1000))
	}
}
