package util

import (
	"bytes"
	"encoding/json"
)

type Option[T any] struct {
	value   T
	isValid bool
}

type OptionState string

const (
	OptionNone OptionState = "none"
	OptionSome OptionState = "some"
)

func NewSomeOption[T any](value T) Option[T] {
	return Option[T]{
		value:   value,
		isValid: true,
	}
}

func NewNoneOption[T any]() Option[T] {
	return Option[T]{
		isValid: false,
	}
}

func (o *Option[T]) State() OptionState {
	if o.isValid {
		return OptionSome
	}

	return OptionNone
}

func (o *Option[T]) Unwrap() T {
	if o.State() == OptionSome {
		return o.value
	}

	return *new(T)
}

func (o *Option[T]) UnwrapOr(defaultValue T) T {
	if o.State() == OptionSome {
		return o.value
	}

	return defaultValue
}

// 实现 JSON Marshaller 接口
func (o Option[T]) MarshalJSON() ([]byte, error) {
	if o.State() == OptionSome {
		return json.Marshal(o.value)
	}

	// None 的 JSON 表示为 null
	return []byte("null"), nil
}

func (o *Option[T]) UnmarshalJSON(data []byte) error {
	// 先检查是否为 null
	if bytes.Equal(data, []byte("null")) {
		o.isValid = false
		return nil
	}

	// 尝试解析为 T 类型
	if err := json.Unmarshal(data, &o.value); err != nil {
		return err
	}

	o.isValid = true

	return nil
}
