package generic

import (
	"fmt"
	"testing"
)

func TestThreadSafeMap(t *testing.T) {
	m := SyncMap[int, int]{}
	m.Store(1, 1)
	val, ok := m.Load(1)
	if !ok {
		t.Fatal("key not found")
	}
	if val != 1 {
		t.Fatal("key not found")
	}
}

type Custom struct {
	A int
	B bool
	C string
}

func TestThreadSafeMap2(t *testing.T) {
	m := SyncMap[string, Custom]{}
	m.Store("1", Custom{A: 1, B: true, C: "1"})
	val, ok := m.Load("1")
	if !ok {
		t.Fatal("key not found")
	}
	fmt.Println(val)
	m.Delete("1")
	_, ok = m.Load("1")
	if ok {
		t.Fatal("expected not found")
	}
}
