package res

type AppRes[T any] struct {
	success bool
	msg     string
	code    int
	data    *T
}

func Reject[T any](code ErrCode, msg string) AppRes[T] {
	return AppRes[T]{
		success: false,
		code:    int(code),
		msg:     msg,
	}
}

func Accept[T any](data *T) AppRes[T] {
	return AppRes[T]{
		success: true,
		data:    data,
	}
}

func (r *AppRes[T]) Code() ErrCode {
	return ErrCode(r.code)
}

func (r *AppRes[T]) Msg() string {
	return r.msg
}

func (r *AppRes[T]) Data() *T {
	return r.data
}

func (r *AppRes[T]) IsAccept() bool {
	return r.success
}

func (r *AppRes[T]) IsReject() bool {
	return !r.success
}

func (r *AppRes[T]) WithCode(code ErrCode) *AppRes[T] {
	r.code = int(code)
	return r
}

type None struct{}
