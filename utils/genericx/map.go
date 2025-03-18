package genericx

import (
	"reflect"
	"sync"
	"sync/atomic"
	"unsafe"
)

// SyncMap 是一个线程安全的 map[K]V
type SyncMap[K comparable, V any] struct {
	mu     sync.Mutex
	read   atomic.Pointer[readOnly[K, V]]
	dirty  map[K]*entry[V]
	misses int
}

func NewSyncMap[K comparable, V any]() *SyncMap[K, V] {
	return &SyncMap[K, V]{}
}

// readOnly 是一个存储在 SyncMap.read 字段中的不可变结构体
type readOnly[K comparable, V any] struct {
	m       map[K]*entry[V]
	amended bool // 如果 dirty map 包含一些不在 m 中的键，则为 true
}

// expunged 是一个任意指针，用于标记已从 dirty map 中删除的条目
var expunged = unsafe.Pointer(new(any))

// entry 是映射中对应特定键的槽位
type entry[V any] struct {
	p atomic.Pointer[V]
}

func newEntry[V any](i V) *entry[V] {
	e := &entry[V]{}
	e.p.Store(&i)
	return e
}

func (m *SyncMap[K, V]) loadReadOnly() readOnly[K, V] {
	if p := m.read.Load(); p != nil {
		return *p
	}
	return readOnly[K, V]{}
}

// Load 返回存储在映射中的键对应的值，如果没有值则返回 nil
// ok 结果表示是否在映射中找到了值
func (m *SyncMap[K, V]) Load(key K) (value V, ok bool) {
	read := m.loadReadOnly()
	e, ok := read.m[key]
	if !ok && read.amended {
		m.mu.Lock()
		read = m.loadReadOnly()
		e, ok = read.m[key]
		if !ok && read.amended {
			e, ok = m.dirty[key]
			m.missLocked()
		}
		m.mu.Unlock()
	}
	if !ok {
		var zero V
		return zero, false
	}
	return e.load()
}

func (e *entry[V]) load() (value V, ok bool) {
	p := e.p.Load()
	if p == nil || unsafe.Pointer(p) == expunged {
		var zero V
		return zero, false
	}
	return *p, true
}

// Store 设置键对应的值
func (m *SyncMap[K, V]) Store(key K, value V) {
	_, _ = m.Swap(key, value)
}

// tryCompareAndSwap 比较条目与给定的旧值，如果相等且条目未被删除，则与新值进行交换
func (e *entry[V]) tryCompareAndSwap(old, new V) bool {
	p := e.p.Load()
	if p == nil || unsafe.Pointer(p) == expunged {
		return false
	}
	if !reflect.DeepEqual(*p, old) {
		return false
	}
	nc := new
	return e.p.CompareAndSwap(p, &nc)
}

// unexpungeLocked 确保条目未被标记为已删除
func (e *entry[V]) unexpungeLocked() (wasExpunged bool) {
	return e.p.CompareAndSwap((*V)(expunged), nil)
}

// swapLocked 无条件地将值交换到条目中
func (e *entry[V]) swapLocked(i *V) *V {
	return e.p.Swap(i)
}

// LoadOrStore 如果键存在则返回现有值
// 否则，存储并返回给定值
// loaded 结果表示值是加载的还是存储的
func (m *SyncMap[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	// 如果是干净命中，避免加锁
	read := m.loadReadOnly()
	if e, ok := read.m[key]; ok {
		actual, loaded, ok := e.tryLoadOrStore(value)
		if ok {
			return actual, loaded
		}
	}

	m.mu.Lock()
	read = m.loadReadOnly()
	if e, ok := read.m[key]; ok {
		if e.unexpungeLocked() {
			m.dirty[key] = e
		}
		actual, loaded, _ = e.tryLoadOrStore(value)
	} else if e, ok := m.dirty[key]; ok {
		actual, loaded, _ = e.tryLoadOrStore(value)
		m.missLocked()
	} else {
		if !read.amended {
			m.dirtyLocked()
			m.read.Store(&readOnly[K, V]{m: read.m, amended: true})
		}
		m.dirty[key] = newEntry[V](value)
		actual, loaded = value, false
	}
	m.mu.Unlock()

	return actual, loaded
}

// tryLoadOrStore 如果条目未被删除，则原子地加载或存储值
func (e *entry[V]) tryLoadOrStore(i V) (actual V, loaded, ok bool) {
	p := e.p.Load()
	if unsafe.Pointer(p) == expunged {
		var zero V
		return zero, false, false
	}
	if p != nil {
		return *p, true, true
	}
	ic := i
	for {
		if e.p.CompareAndSwap(nil, &ic) {
			return i, false, true
		}
		p = e.p.Load()
		if unsafe.Pointer(p) == expunged {
			var zero V
			return zero, false, false
		}
		if p != nil {
			return *p, true, true
		}
	}
}

// LoadAndDelete 删除键对应的值，如果有则返回之前的值
// loaded 结果报告键是否存在
func (m *SyncMap[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	read := m.loadReadOnly()
	e, ok := read.m[key]
	if !ok && read.amended {
		m.mu.Lock()
		read = m.loadReadOnly()
		e, ok = read.m[key]
		if !ok && read.amended {
			e, ok = m.dirty[key]
			delete(m.dirty, key)
			m.missLocked()
		}
		m.mu.Unlock()
	}
	if ok {
		return e.delete()
	}
	var zero V
	return zero, false
}

// Delete 删除键对应的值
func (m *SyncMap[K, V]) Delete(key K) {
	m.LoadAndDelete(key)
}

func (e *entry[V]) delete() (value V, ok bool) {
	for {
		p := e.p.Load()
		if p == nil || unsafe.Pointer(p) == expunged {
			var zero V
			return zero, false
		}
		if e.p.CompareAndSwap(p, nil) {
			return *p, true
		}
	}
}

// trySwap 如果条目未被删除，则交换值
func (e *entry[V]) trySwap(i *V) (*V, bool) {
	for {
		p := e.p.Load()
		if unsafe.Pointer(p) == expunged {
			return nil, false
		}
		if e.p.CompareAndSwap(p, i) {
			return p, true
		}
	}
}

// Swap 交换键对应的值并返回之前的值（如果有）
// loaded 结果报告键是否存在
func (m *SyncMap[K, V]) Swap(key K, value V) (previous any, loaded bool) {
	read := m.loadReadOnly()
	if e, ok := read.m[key]; ok {
		if v, ok := e.trySwap(&value); ok {
			if v == nil {
				return nil, false
			}
			return *v, true
		}
	}

	m.mu.Lock()
	read = m.loadReadOnly()
	if e, ok := read.m[key]; ok {
		if e.unexpungeLocked() {
			m.dirty[key] = e
		}
		if v := e.swapLocked(&value); v != nil {
			loaded = true
			previous = *v
		}
	} else if e, ok := m.dirty[key]; ok {
		if v := e.swapLocked(&value); v != nil {
			loaded = true
			previous = *v
		}
	} else {
		if !read.amended {
			m.dirtyLocked()
			m.read.Store(&readOnly[K, V]{m: read.m, amended: true})
		}
		m.dirty[key] = newEntry(value)
	}
	m.mu.Unlock()
	return previous, loaded
}

// CompareAndSwap 如果映射中存储的值等于旧值，则交换键的旧值和新值
// 旧值必须是可比较类型
func (m *SyncMap[K, V]) CompareAndSwap(key K, old, new V) bool {
	read := m.loadReadOnly()
	if e, ok := read.m[key]; ok {
		return e.tryCompareAndSwap(old, new)
	} else if !read.amended {
		return false // 键没有现有值
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	read = m.loadReadOnly()
	swapped := false
	if e, ok := read.m[key]; ok {
		swapped = e.tryCompareAndSwap(old, new)
	} else if e, ok := m.dirty[key]; ok {
		swapped = e.tryCompareAndSwap(old, new)
		// 我们需要锁定 mu 来加载键的条目，
		// 并且操作没有改变映射中的键集
		// （所以通过将 dirty map 提升为只读可以提高效率）。
		// 将其计为未命中，以便最终切换到更高效的稳定状态。
		m.missLocked()
	}
	return swapped
}

// CompareAndDelete 如果键的值等于旧值，则删除该条目
// 旧值必须是可比较类型
//
// 如果映射中没有键的当前值，CompareAndDelete 返回 false
// （即使旧值是 nil 接口值）
func (m *SyncMap[K, V]) CompareAndDelete(key K, old V) (deleted bool) {
	read := m.loadReadOnly()
	e, ok := read.m[key]
	if !ok && read.amended {
		m.mu.Lock()
		read = m.loadReadOnly()
		e, ok = read.m[key]
		if !ok && read.amended {
			e, ok = m.dirty[key]
			m.missLocked()
		}
		m.mu.Unlock()
	}
	if !ok {
		return false
	}
	for {
		p := e.p.Load()
		if p == nil || unsafe.Pointer(p) == expunged {
			return false
		}
		if !reflect.DeepEqual(*p, old) {
			return false
		}
		if e.p.CompareAndSwap(p, nil) {
			return true
		}
	}
}

// Range 按顺序为映射中的每个键和值调用 f
// 如果 f 返回 false，range 停止迭代
func (m *SyncMap[K, V]) Range(f func(key K, value V) bool) {
	read := m.loadReadOnly()
	if read.amended {
		m.mu.Lock()
		read = m.loadReadOnly()
		if read.amended {
			read = readOnly[K, V]{m: m.dirty}
			copyRead := read
			m.read.Store(&copyRead)
			m.dirty = nil
			m.misses = 0
		}
		m.mu.Unlock()
	}

	for k, e := range read.m {
		v, ok := e.load()
		if !ok {
			continue
		}
		if !f(k, v) {
			break
		}
	}
}

func (m *SyncMap[K, V]) missLocked() {
	m.misses++
	if m.misses < len(m.dirty) {
		return
	}
	m.read.Store(&readOnly[K, V]{m: m.dirty})
	m.dirty = nil
	m.misses = 0
}

func (m *SyncMap[K, V]) dirtyLocked() {
	if m.dirty != nil {
		return
	}

	read := m.loadReadOnly()
	m.dirty = make(map[K]*entry[V], len(read.m))
	for k, e := range read.m {
		if !e.tryExpungeLocked() {
			m.dirty[k] = e
		}
	}
}

func (e *entry[V]) tryExpungeLocked() (isExpunged bool) {
	p := e.p.Load()
	for p == nil {
		if e.p.CompareAndSwap(nil, (*V)(expunged)) {
			return true
		}
		p = e.p.Load()
	}
	return unsafe.Pointer(p) == expunged
}
