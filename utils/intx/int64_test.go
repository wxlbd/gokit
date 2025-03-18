package intx

import (
	"encoding/json"
	"testing"
)

func TestInt64String_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    Int64String
		wantErr bool
	}{
		{
			name:    "valid number string",
			json:    `"9223372036854775807"`,
			want:    Int64String(9223372036854775807),
			wantErr: false,
		},
		{
			name:    "valid small number",
			json:    `"123"`,
			want:    Int64String(123),
			wantErr: false,
		},
		{
			name:    "invalid format",
			json:    `abc`,
			want:    Int64String(0),
			wantErr: true,
		},
		{
			name:    "empty string",
			json:    `""`,
			want:    Int64String(0),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Int64String
			err := json.Unmarshal([]byte(tt.json), &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("UnmarshalJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64String_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		i       Int64String
		want    string
		wantErr bool
	}{
		{
			name:    "large number",
			i:       Int64String(9223372036854775807),
			want:    `"9223372036854775807"`,
			wantErr: false,
		},
		{
			name:    "small number",
			i:       Int64String(123),
			want:    `"123"`,
			wantErr: false,
		},
		{
			name:    "zero",
			i:       Int64String(0),
			want:    `"0"`,
			wantErr: false,
		},
		{
			name:    "negative number",
			i:       Int64String(-9223372036854775808),
			want:    `"-9223372036854775808"`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && string(got) != tt.want {
				t.Errorf("MarshalJSON() got = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestInt64String_Int64(t *testing.T) {
	tests := []struct {
		name string
		i    Int64String
		want int64
	}{
		{
			name: "large number",
			i:    Int64String(9223372036854775807),
			want: 9223372036854775807,
		},
		{
			name: "small number",
			i:    Int64String(123),
			want: 123,
		},
		{
			name: "zero",
			i:    Int64String(0),
			want: 0,
		},
		{
			name: "negative number",
			i:    Int64String(-9223372036854775808),
			want: -9223372036854775808,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.Int64(); got != tt.want {
				t.Errorf("Int64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromInt64(t *testing.T) {
	tests := []struct {
		name string
		v    int64
		want Int64String
	}{
		{
			name: "large number",
			v:    9223372036854775807,
			want: Int64String(9223372036854775807),
		},
		{
			name: "small number",
			v:    123,
			want: Int64String(123),
		},
		{
			name: "zero",
			v:    0,
			want: Int64String(0),
		},
		{
			name: "negative number",
			v:    -9223372036854775808,
			want: Int64String(-9223372036854775808),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromInt64(tt.v); got != tt.want {
				t.Errorf("FromInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromString(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    Int64String
		wantErr bool
	}{
		{
			name:    "valid large number",
			s:       "9223372036854775807",
			want:    Int64String(9223372036854775807),
			wantErr: false,
		},
		{
			name:    "valid small number",
			s:       "123",
			want:    Int64String(123),
			wantErr: false,
		},
		{
			name:    "invalid format",
			s:       "abc",
			want:    Int64String(0),
			wantErr: true,
		},
		{
			name:    "empty string",
			s:       "",
			want:    Int64String(0),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromString(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("FromString() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkInt64String_MarshalJSON(b *testing.B) {
	i := Int64String(9223372036854775807)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = i.MarshalJSON()
	}
}

func BenchmarkInt64String_UnmarshalJSON(b *testing.B) {
	data := []byte(`"9223372036854775807"`)
	var i Int64String
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = i.UnmarshalJSON(data)
	}
}
