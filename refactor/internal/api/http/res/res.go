package res

import "github.com/kataras/iris/v12"

type HttpRes struct {
	Code int    `json:"code"`
	Msg  string `json:"message,omitempty"`
	Data any    `json:"data,omitempty"`
}

func Accept(cx iris.Context, code int, data any) {
	cx.StatusCode(code)

	if data != nil {
		_ = cx.JSON(HttpRes{
			Code: code,
			Data: data,
		})
	}
}

func Reject(cx iris.Context, code int, msg string) {
	cx.StatusCode(code)

	if msg != "" {
		_ = cx.JSON(HttpRes{
			Code: code,
			Msg:  msg,
		})
	}
}
