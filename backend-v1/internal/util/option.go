package util

import (
	"bytes"
	"encoding/json"
	"reflect"
)

// 无法在内部处理在序列化时 “不传字段” 的操作，只能外部通过指针 + omitempty 实现
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
	// 如果是 None，则返回 null
	if o.State() == OptionNone {
		return []byte("null"), nil
	}

	// 如果是 Some，需要反射检查是否为指针类型，如果是指针类型且值为 nil，则也返回 null
	v := reflect.ValueOf(o.value)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return []byte("null"), nil
	}

	// 否则正常序列化值
	return json.Marshal(o.value)
}

func (o *Option[T]) UnmarshalJSON(data []byte) error {
	// 只要这个函数被调用，就说明 JSON 中存在这个字段了，所以不用处理 None 的情况
	// 先检查是否为 null，需要对指针类型特判
	if bytes.Equal(data, []byte("null")) {
		// 为了防止 T 是一个接口，构造 T* 指针再取 T
		t := reflect.TypeOf((*T)(nil)).Elem()
		if t.Kind() == reflect.Ptr {
			o.isValid = true
			o.value = *new(T)
		} else {
			o.isValid = false
		}

		return nil
	}

	// 尝试解析为 T 类型
	if err := json.Unmarshal(data, &o.value); err != nil {
		return err
	}

	o.isValid = true

	return nil
}
