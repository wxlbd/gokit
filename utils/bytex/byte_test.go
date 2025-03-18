package bytex

import (
	"testing"
)

func TestBytesToString(t *testing.T) {
	tests := []struct {
		name string
		b    []byte
		want string
	}{
		{
			name: "empty bytes",
			b:    []byte{},
			want: "",
		},
		{
			name: "ascii bytes",
			b:    []byte("Hello, World!"),
			want: "Hello, World!",
		},
		{
			name: "unicode bytes",
			b:    []byte("你好，世界！"),
			want: "你好，世界！",
		},
		{
			name: "special characters",
			b:    []byte("!@#$%^&*()_+"),
			want: "!@#$%^&*()_+",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BytesToString(tt.b); got != tt.want {
				t.Errorf("BytesToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringToBytes(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want []byte
	}{
		{
			name: "empty string",
			s:    "",
			want: []byte{},
		},
		{
			name: "ascii string",
			s:    "Hello, World!",
			want: []byte("Hello, World!"),
		},
		{
			name: "unicode string",
			s:    "你好，世界！",
			want: []byte("你好，世界！"),
		},
		{
			name: "special characters",
			s:    "!@#$%^&*()_+",
			want: []byte("!@#$%^&*()_+"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StringToBytes(tt.s)
			if len(got) != len(tt.want) {
				t.Errorf("StringToBytes() length = %v, want %v", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("StringToBytes() at index %d = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func BenchmarkBytesToString(b *testing.B) {
	sample := []byte("Hello, World! 你好，世界！")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BytesToString(sample)
	}
}

func BenchmarkStringToBytes(b *testing.B) {
	sample := "Hello, World! 你好，世界！"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StringToBytes(sample)
	}
}
