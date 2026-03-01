package util

type Option[T any] struct {
	value   *T
	isValid bool
}

func NewSomeOption[T any](value T) Option[T] {
	return Option[T]{
		value:   &value,
		isValid: true,
	}
}

func NewNoneOption[T any]() Option[T] {
	return Option[T]{
		value:   nil,
		isValid: false,
	}
}

func (o *Option[T]) IsSome() bool {
	return o.isValid
}

func (o *Option[T]) IsNone() bool {
	return !o.isValid
}

func (o *Option[T]) TryUnwrap() *T {
	if o.IsSome() {
		return o.value
	}

	return nil
}

func (o *Option[T]) UnwrapOr(defaultValue T) *T {
	if o.IsSome() {
		return o.value
	}

	return &defaultValue
}
