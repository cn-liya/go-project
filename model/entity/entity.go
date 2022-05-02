package entity

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
)

const (
	StatusOff = -1
	StatusOn  = 1
)

type JsonSlice[T any] []T

func (v JsonSlice[T]) Value() (driver.Value, error) {
	if v == nil {
		return []byte{'[', ']'}, nil
	}
	return json.Marshal(v)
}

func (v *JsonSlice[T]) Scan(value any) error {
	if value == nil {
		return nil
	}
	b := value.([]byte)
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	return dec.Decode(v)
}

type JsonMap[K comparable, V any] map[K]V

func (v JsonMap[K, V]) Value() (driver.Value, error) {
	if v == nil {
		return []byte{'{', '}'}, nil
	}
	return json.Marshal(v)
}

func (v *JsonMap[K, V]) Scan(value any) error {
	if value == nil {
		return nil
	}
	b := value.([]byte)
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	return dec.Decode(v)
}

func JsonString(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}
