package genericx

import (
	"testing"
)

func TestSyncMap_Load(t *testing.T) {
	m := NewSyncMap[string, int]()
	m.Store("key1", 1)

	tests := []struct {
		name   string
		key    string
		want   int
		wantOk bool
	}{
		{
			name:   "existing key",
			key:    "key1",
			want:   1,
			wantOk: true,
		},
		{
			name:   "non-existent key",
			key:    "key2",
			want:   0,
			wantOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := m.Load(tt.key)
			if ok != tt.wantOk {
				t.Errorf("Load() ok = %v, want %v", ok, tt.wantOk)
				return
			}
			if tt.wantOk && got != tt.want {
				t.Errorf("Load() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSyncMap_Store(t *testing.T) {
	m := NewSyncMap[string, int]()
	m.Store("key1", 1)

	// Verify the value was stored
	if got, ok := m.Load("key1"); !ok || got != 1 {
		t.Errorf("Store() failed to store value")
	}
}

func TestSyncMap_LoadOrStore(t *testing.T) {
	m := NewSyncMap[string, int]()
	m.Store("key1", 1)

	tests := []struct {
		name     string
		key      string
		value    int
		want     int
		wantLoad bool
	}{
		{
			name:     "existing key",
			key:      "key1",
			value:    2,
			want:     1,
			wantLoad: true,
		},
		{
			name:     "non-existent key",
			key:      "key2",
			value:    2,
			want:     2,
			wantLoad: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, loaded := m.LoadOrStore(tt.key, tt.value)
			if loaded != tt.wantLoad {
				t.Errorf("LoadOrStore() loaded = %v, want %v", loaded, tt.wantLoad)
				return
			}
			if got != tt.want {
				t.Errorf("LoadOrStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSyncMap_Delete(t *testing.T) {
	m := NewSyncMap[string, int]()
	m.Store("key1", 1)

	// Delete existing key
	m.Delete("key1")
	if _, ok := m.Load("key1"); ok {
		t.Error("Delete() failed to delete existing key")
	}

	// Delete non-existent key (should not panic)
	m.Delete("key2")
}

func TestSyncMap_LoadAndDelete(t *testing.T) {
	m := NewSyncMap[string, int]()
	m.Store("key1", 1)

	tests := []struct {
		name     string
		key      string
		want     int
		wantLoad bool
	}{
		{
			name:     "existing key",
			key:      "key1",
			want:     1,
			wantLoad: true,
		},
		{
			name:     "non-existent key",
			key:      "key2",
			want:     0,
			wantLoad: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, loaded := m.LoadAndDelete(tt.key)
			if loaded != tt.wantLoad {
				t.Errorf("LoadAndDelete() loaded = %v, want %v", loaded, tt.wantLoad)
				return
			}
			if loaded && got != tt.want {
				t.Errorf("LoadAndDelete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSyncMap_Swap(t *testing.T) {
	m := NewSyncMap[string, int]()
	m.Store("key1", 1)

	tests := []struct {
		name     string
		key      string
		value    int
		want     int
		wantLoad bool
	}{
		{
			name:     "existing key",
			key:      "key1",
			value:    2,
			want:     1,
			wantLoad: true,
		},
		{
			name:     "non-existent key",
			key:      "key2",
			value:    2,
			want:     0,
			wantLoad: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, loaded := m.Swap(tt.key, tt.value)
			if loaded != tt.wantLoad {
				t.Errorf("Swap() loaded = %v, want %v", loaded, tt.wantLoad)
				return
			}
			if loaded && got != tt.want {
				t.Errorf("Swap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSyncMap_CompareAndSwap(t *testing.T) {
	m := NewSyncMap[string, int]()
	m.Store("key1", 1)

	tests := []struct {
		name string
		key  string
		old  int
		new  int
		want bool
	}{
		{
			name: "successful swap",
			key:  "key1",
			old:  1,
			new:  2,
			want: true,
		},
		{
			name: "failed swap - wrong old value",
			key:  "key1",
			old:  3,
			new:  4,
			want: false,
		},
		{
			name: "non-existent key",
			key:  "key2",
			old:  1,
			new:  2,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := m.CompareAndSwap(tt.key, tt.old, tt.new); got != tt.want {
				t.Errorf("CompareAndSwap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSyncMap_CompareAndDelete(t *testing.T) {
	m := NewSyncMap[string, int]()
	m.Store("key1", 1)

	tests := []struct {
		name string
		key  string
		old  int
		want bool
	}{
		{
			name: "successful delete",
			key:  "key1",
			old:  1,
			want: true,
		},
		{
			name: "failed delete - wrong old value",
			key:  "key1",
			old:  2,
			want: false,
		},
		{
			name: "non-existent key",
			key:  "key2",
			old:  1,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := m.CompareAndDelete(tt.key, tt.old); got != tt.want {
				t.Errorf("CompareAndDelete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSyncMap_Range(t *testing.T) {
	m := NewSyncMap[string, int]()
	m.Store("key1", 1)
	m.Store("key2", 2)

	got := make(map[string]int)
	m.Range(func(key string, value int) bool {
		got[key] = value
		return true
	})

	want := map[string]int{
		"key1": 1,
		"key2": 2,
	}

	if len(got) != len(want) {
		t.Errorf("Range() got %d items, want %d", len(got), len(want))
		return
	}

	for k, v := range want {
		if got[k] != v {
			t.Errorf("Range() got[%s] = %v, want %v", k, got[k], v)
		}
	}
}

func TestSyncMap_Concurrent(t *testing.T) {
	m := NewSyncMap[int, int]()
	const numGoroutines = 100
	const numOperations = 1000

	// Start multiple goroutines that perform operations on the map
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			for j := 0; j < numOperations; j++ {
				key := id*numOperations + j
				m.Store(key, key)
				if _, ok := m.Load(key); !ok {
					t.Errorf("Load() failed to find key %d", key)
				}
				m.Delete(key)
				if _, ok := m.Load(key); ok {
					t.Errorf("Delete() failed to delete key %d", key)
				}
			}
		}(i)
	}
}
