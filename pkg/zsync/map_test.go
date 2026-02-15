package zsync_test

import (
	"sync"
	"testing"

	"github.com/zarlcorp/core/pkg/zsync"
)

func TestZMap_Set(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value int
	}{
		{
			name:  "set string key with int value",
			key:   "test",
			value: 42,
		},
		{
			name:  "set empty string key",
			key:   "",
			value: 0,
		},
		{
			name:  "overwrite existing key",
			key:   "existing",
			value: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := zsync.NewZMap[string, int]()
			m.Set(tt.key, tt.value)

			got, ok := m.Get(tt.key)
			if !ok {
				t.Errorf("Get() ok = false, want true")
				return
			}
			if got != tt.value {
				t.Errorf("Get() = %v, want %v", got, tt.value)
			}
		})
	}
}

func TestZMap_Get(t *testing.T) {
	tests := []struct {
		name   string
		setup  map[string]int
		key    string
		want   int
		wantOK bool
	}{
		{
			name:   "get existing key",
			setup:  map[string]int{"test": 42},
			key:    "test",
			want:   42,
			wantOK: true,
		},
		{
			name:   "get non-existent key",
			setup:  map[string]int{},
			key:    "missing",
			want:   0,
			wantOK: false,
		},
		{
			name:   "get from empty map",
			setup:  nil,
			key:    "any",
			want:   0,
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := zsync.NewZMap[string, int]()

			// setup
			for k, v := range tt.setup {
				m.Set(k, v)
			}

			got, ok := m.Get(tt.key)
			if ok != tt.wantOK {
				t.Errorf("Get() ok = %v, want %v", ok, tt.wantOK)
				return
			}
			if got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZMap_Delete(t *testing.T) {
	tests := []struct {
		name   string
		setup  map[string]int
		key    string
		want   bool
		length int
	}{
		{
			name:   "delete existing key",
			setup:  map[string]int{"test": 42, "other": 1},
			key:    "test",
			want:   true,
			length: 1,
		},
		{
			name:   "delete non-existent key",
			setup:  map[string]int{"test": 42},
			key:    "missing",
			want:   false,
			length: 1,
		},
		{
			name:   "delete from empty map",
			setup:  nil,
			key:    "any",
			want:   false,
			length: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := zsync.NewZMap[string, int]()

			// setup
			for k, v := range tt.setup {
				m.Set(k, v)
			}

			got := m.Delete(tt.key)
			if got != tt.want {
				t.Errorf("Delete() = %v, want %v", got, tt.want)
			}

			if m.Len() != tt.length {
				t.Errorf("Len() after delete = %v, want %v", m.Len(), tt.length)
			}
		})
	}
}

func TestZMap_Len(t *testing.T) {
	tests := []struct {
		name  string
		setup map[string]int
		want  int
	}{
		{
			name:  "empty map",
			setup: nil,
			want:  0,
		},
		{
			name:  "single item",
			setup: map[string]int{"test": 42},
			want:  1,
		},
		{
			name:  "multiple items",
			setup: map[string]int{"a": 1, "b": 2, "c": 3},
			want:  3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := zsync.NewZMap[string, int]()

			for k, v := range tt.setup {
				m.Set(k, v)
			}

			if got := m.Len(); got != tt.want {
				t.Errorf("Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZMap_Keys(t *testing.T) {
	tests := []struct {
		name  string
		setup map[string]int
		want  []string
	}{
		{
			name:  "empty map",
			setup: nil,
			want:  []string{},
		},
		{
			name:  "single key",
			setup: map[string]int{"test": 42},
			want:  []string{"test"},
		},
		{
			name:  "multiple keys",
			setup: map[string]int{"a": 1, "b": 2},
			want:  []string{"a", "b"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := zsync.NewZMap[string, int]()

			for k, v := range tt.setup {
				m.Set(k, v)
			}

			got := m.Keys()
			if len(got) != len(tt.want) {
				t.Errorf("Keys() length = %v, want %v", len(got), len(tt.want))
				return
			}

			// unordered comparison
			gotMap := make(map[string]bool)
			for _, k := range got {
				gotMap[k] = true
			}

			for _, wantKey := range tt.want {
				if !gotMap[wantKey] {
					t.Errorf("Keys() missing key %v", wantKey)
				}
			}
		})
	}
}

func TestZMap_Clear(t *testing.T) {
	m := zsync.NewZMap[string, int]()
	m.Set("a", 1)
	m.Set("b", 2)

	if m.Len() != 2 {
		t.Errorf("Len() before clear = %v, want 2", m.Len())
	}

	m.Clear()

	if m.Len() != 0 {
		t.Errorf("Len() after clear = %v, want 0", m.Len())
	}

	_, ok := m.Get("a")
	if ok {
		t.Errorf("Get() after clear ok = true, want false")
	}
}

// concurrent access test
func TestZMap_Concurrent(t *testing.T) {
	m := zsync.NewZMap[int, string]()

	const numGoroutines = 100
	const numOperations = 1000

	var wg sync.WaitGroup

	// concurrent writes
	for i := range numGoroutines {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := range numOperations {
				key := id*numOperations + j
				m.Set(key, "value")
			}
		}(i)
	}

	// concurrent reads
	for i := range numGoroutines {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := range numOperations {
				key := id*numOperations + j
				m.Get(key) // ignore result, just testing for races
			}
		}(i)
	}

	wg.Wait()

	expectedLen := numGoroutines * numOperations
	if m.Len() != expectedLen {
		t.Errorf("Len() after concurrent operations = %v, want %v", m.Len(), expectedLen)
	}
}
