package bimap

import (
	"fmt"
	"testing"
)

func TestBiMap_BasicOperations(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    int
		expected bool
	}{
		{"添加新键值对", "key1", 1, true},
		{"添加重复键", "key1", 2, true},
		{"添加重复值", "key2", 1, true},
	}

	bimap := NewBiMap[string, int]()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bimap.Put(tt.key, tt.value)
			value, exists := bimap.GetByKey(tt.key)
			if !exists || value != tt.value {
				t.Errorf("GetByKey() = (%v, %v), want (%v, %v)", value, exists, tt.value, tt.expected)
			}
		})
	}
}

func TestBiMap_GetOperations(t *testing.T) {
	bimap := NewBiMap[string, int]()
	bimap.Put("key1", 1)
	bimap.Put("key2", 2)

	tests := []struct {
		name   string
		key    string
		value  int
		exists bool
	}{
		{"获取存在的键", "key1", 1, true},
		{"获取不存在的键", "key3", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, exists := bimap.GetByKey(tt.key)
			if exists != tt.exists || value != tt.value {
				t.Errorf("GetByKey() = (%v, %v), want (%v, %v)", value, exists, tt.value, tt.exists)
			}
		})
	}
}

func TestBiMap_GetByValue(t *testing.T) {
	bimap := NewBiMap[string, int]()
	bimap.Put("key1", 1)
	bimap.Put("key2", 2)

	tests := []struct {
		name   string
		value  int
		key    string
		exists bool
	}{
		{"获取存在的值", 1, "key1", true},
		{"获取不存在的值", 3, "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, exists := bimap.GetByValue(tt.value)
			if exists != tt.exists || key != tt.key {
				t.Errorf("GetByValue() = (%v, %v), want (%v, %v)", key, exists, tt.key, tt.exists)
			}
		})
	}
}

func TestBiMap_DeleteOperations(t *testing.T) {
	bimap := NewBiMap[string, int]()
	bimap.Put("key1", 1)
	bimap.Put("key2", 2)

	tests := []struct {
		name   string
		key    string
		value  int
		exists bool
	}{
		{"删除存在的键", "key1", 1, false},
		{"删除不存在的键", "key3", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bimap.DeleteByKey(tt.key)
			_, exists := bimap.GetByKey(tt.key)
			if exists != tt.exists {
				t.Errorf("DeleteByKey() 后 GetByKey() = %v, want %v", exists, tt.exists)
			}
		})
	}
}

func TestBiMap_DeleteByValue(t *testing.T) {
	bimap := NewBiMap[string, int]()
	bimap.Put("key1", 1)
	bimap.Put("key2", 2)

	tests := []struct {
		name   string
		value  int
		key    string
		exists bool
	}{
		{"删除存在的值", 1, "key1", false},
		{"删除不存在的值", 3, "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bimap.DeleteByValue(tt.value)
			_, exists := bimap.GetByValue(tt.value)
			if exists != tt.exists {
				t.Errorf("DeleteByValue() 后 GetByValue() = %v, want %v", exists, tt.exists)
			}
		})
	}
}

func TestBiMap_Len(t *testing.T) {
	bimap := NewBiMap[string, int]()
	bimap.Put("key1", 1)
	bimap.Put("key2", 2)

	if bimap.Len() != 2 {
		t.Errorf("Len() = %v, want %v", bimap.Len(), 2)
	}

	bimap.DeleteByKey("key1")
	if bimap.Len() != 1 {
		t.Errorf("Len() = %v, want %v", bimap.Len(), 1)
	}
}

func TestBiMap_Clear(t *testing.T) {
	bimap := NewBiMap[string, int]()
	bimap.Put("key1", 1)
	bimap.Put("key2", 2)

	bimap.Clear()
	if bimap.Len() != 0 {
		t.Errorf("Clear() 后 Len() = %v, want %v", bimap.Len(), 0)
	}
}

func TestBiMap_Keys(t *testing.T) {
	bimap := NewBiMap[string, int]()
	bimap.Put("key1", 1)
	bimap.Put("key2", 2)

	keys := bimap.Keys()
	if len(keys) != 2 {
		t.Errorf("Keys() 长度 = %v, want %v", len(keys), 2)
	}

	// 检查是否包含所有键
	expectedKeys := map[string]bool{"key1": true, "key2": true}
	for _, key := range keys {
		if !expectedKeys[key] {
			t.Errorf("Keys() 包含意外的键: %v", key)
		}
	}
}

func TestBiMap_Values(t *testing.T) {
	bimap := NewBiMap[string, int]()
	bimap.Put("key1", 1)
	bimap.Put("key2", 2)

	values := bimap.Values()
	if len(values) != 2 {
		t.Errorf("Values() 长度 = %v, want %v", len(values), 2)
	}

	// 检查是否包含所有值
	expectedValues := map[int]bool{1: true, 2: true}
	for _, value := range values {
		if !expectedValues[value] {
			t.Errorf("Values() 包含意外的值: %v", value)
		}
	}
}

func TestBiMap_ForEach(t *testing.T) {
	bimap := NewBiMap[string, int]()
	bimap.Put("key1", 1)
	bimap.Put("key2", 2)

	visited := make(map[string]int)
	bimap.ForEach(func(key string, value int) {
		visited[key] = value
	})

	if len(visited) != 2 {
		t.Errorf("ForEach() 访问的键值对数量 = %v, want %v", len(visited), 2)
	}

	// 检查访问的键值对是否正确
	if visited["key1"] != 1 || visited["key2"] != 2 {
		t.Errorf("ForEach() 访问的键值对不正确: %v", visited)
	}
}

func TestBiMap_GetOrDefault(t *testing.T) {
	bimap := NewBiMap[string, int]()
	bimap.Put("key1", 1)

	tests := []struct {
		name          string
		key           string
		defaultValue  int
		expectedValue int
	}{
		{"获取存在的键", "key1", 0, 1},
		{"获取不存在的键", "key2", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value := bimap.GetOrDefault(tt.key, tt.defaultValue)
			if value != tt.expectedValue {
				t.Errorf("GetOrDefault() = %v, want %v", value, tt.expectedValue)
			}
		})
	}
}

func TestBiMap_Concurrent(t *testing.T) {
	bimap := NewBiMap[string, int]()
	done := make(chan bool)
	const numGoroutines = 100

	// 启动多个 goroutine 进行并发操作
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			key := fmt.Sprintf("key%d", id)
			value := id

			// 写入操作
			bimap.Put(key, value)

			// 读取操作
			gotValue, exists := bimap.GetByKey(key)
			if !exists || gotValue != value {
				t.Errorf("并发 GetByKey() = (%v, %v), want (%v, true)", gotValue, exists, value)
			}

			gotKey, exists := bimap.GetByValue(value)
			if !exists || gotKey != key {
				t.Errorf("并发 GetByValue() = (%v, %v), want (%v, true)", gotKey, exists, key)
			}

			// 删除操作
			bimap.DeleteByKey(key)

			// 验证删除
			_, exists = bimap.GetByKey(key)
			if exists {
				t.Errorf("并发 DeleteByKey() 后 GetByKey() 仍然存在键: %v", key)
			}

			done <- true
		}(i)
	}

	// 等待所有 goroutine 完成
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// 验证最终状态
	if bimap.Len() != 0 {
		t.Errorf("并发操作后 Len() = %v, want %v", bimap.Len(), 0)
	}

	// 验证所有键和值都已被删除
	keys := bimap.Keys()
	if len(keys) != 0 {
		t.Errorf("并发操作后 Keys() 长度 = %v, want %v", len(keys), 0)
	}

	values := bimap.Values()
	if len(values) != 0 {
		t.Errorf("并发操作后 Values() 长度 = %v, want %v", len(values), 0)
	}
}
