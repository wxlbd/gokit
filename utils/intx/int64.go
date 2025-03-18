package intx

import (
	"bytes"
	"fmt"
	"strconv"
)

// Int64String 是一个自定义的 int64 类型，用于解决 JSON 序列化时 JavaScript 数值精度丢失的问题。
// 它在序列化时会将数值转换为字符串格式，反序列化时再转回 int64。
type Int64String int64

func (i *Int64String) UnmarshalJSON(b []byte) error {
	b = bytes.Trim(b, "\"")
	parsed, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return fmt.Errorf("parse Int64String: %w", err)
	}
	*i = Int64String(parsed)
	return nil
}

func (i Int64String) MarshalJSON() ([]byte, error) {
	return []byte(`"` + strconv.FormatInt(int64(i), 10) + `"`), nil
}

func (i Int64String) Int64() int64 {
	return int64(i)
}

func FromInt64(v int64) Int64String {
	return Int64String(v)
}

// FromString 将数字字符串转换为 Int64String
func FromString(s string) (Int64String, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse Int64String: %w", err)
	}
	return Int64String(i), nil
}
