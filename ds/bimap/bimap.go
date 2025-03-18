package bimap

import "sync"

// BiMap 双向映射结构，使用泛型K和V分别表示key和value的类型
type BiMap[K comparable, V comparable] struct {
	forward map[K]V      // 正向映射
	reverse map[V]K      // 反向映射
	mu      sync.RWMutex // 读写锁
}

// NewBiMap 创建新的双向映射
func NewBiMap[K comparable, V comparable]() *BiMap[K, V] {
	return &BiMap[K, V]{
		forward: make(map[K]V),
		reverse: make(map[V]K),
	}
}

// Put 添加键值对
func (b *BiMap[K, V]) Put(key K, value V) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// 如果key已存在，先删除旧的映射
	if oldValue, exists := b.forward[key]; exists {
		delete(b.reverse, oldValue)
	}

	// 如果value已存在，先删除旧的映射
	if oldKey, exists := b.reverse[value]; exists {
		delete(b.forward, oldKey)
	}

	b.forward[key] = value
	b.reverse[value] = key
}

// GetByKey 通过key获取value
func (b *BiMap[K, V]) GetByKey(key K) (V, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	value, exists := b.forward[key]
	return value, exists
}

// GetByValue 通过value获取key
func (b *BiMap[K, V]) GetByValue(value V) (K, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	key, exists := b.reverse[value]
	return key, exists
}

// DeleteByKey 通过key删除映射
func (b *BiMap[K, V]) DeleteByKey(key K) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if value, exists := b.forward[key]; exists {
		delete(b.forward, key)
		delete(b.reverse, value)
	}
}

// DeleteByValue 通过value删除映射
func (b *BiMap[K, V]) DeleteByValue(value V) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if key, exists := b.reverse[value]; exists {
		delete(b.forward, key)
		delete(b.reverse, value)
	}
}

// Len 获取映射数量
func (b *BiMap[K, V]) Len() int {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return len(b.forward)
}

// Clear 清空映射
func (b *BiMap[K, V]) Clear() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.forward = make(map[K]V)
	b.reverse = make(map[V]K)
}

// Keys 获取所有key
func (b *BiMap[K, V]) Keys() []K {
	b.mu.RLock()
	defer b.mu.RUnlock()

	keys := make([]K, 0, len(b.forward))
	for k := range b.forward {
		keys = append(keys, k)
	}
	return keys
}

// Values 获取所有value
func (b *BiMap[K, V]) Values() []V {
	b.mu.RLock()
	defer b.mu.RUnlock()

	values := make([]V, 0, len(b.reverse))
	for v := range b.reverse {
		values = append(values, v)
	}
	return values
}

// ForEach 遍历所有键值对
func (b *BiMap[K, V]) ForEach(fn func(K, V)) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for k, v := range b.forward {
		fn(k, v)
	}
}

// GetOrDefault 通过key获取value，如果不存在则返回默认值
func (b *BiMap[K, V]) GetOrDefault(key K, defaultValue V) V {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if value, exists := b.forward[key]; exists {
		return value
	}
	return defaultValue
}
