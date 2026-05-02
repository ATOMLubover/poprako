package res

import "github.com/kataras/iris/v12"

// `HttpRes` is the standard HTTP JSON response envelope
// It wraps every API response body so that clients can uniformly
// inspect `code` and `message` before consuming `data`
type HttpRes[T any] struct {
	// `Code` is the HTTP status code mirrored in JSON for client convenience
	Code int `json:"code"`

	// `Msg` carries a human-readable error message; empty on success
	Msg string `json:"message,omitempty"`

	// `Data` holds the response payload; nil when no body is returned
	Data T `json:"data,omitempty"`
}

// `Accept` writes a success response with the given HTTP status code
// When `data` is non-nil it is wrapped inside `HttpRes`
// When `data` is nil only the status code is written without a JSON body
func Accept[T any](cx iris.Context, code int, data T) {
	cx.StatusCode(code)

	_ = cx.JSON(HttpRes[T]{
		Code: code,
		Data: data,
	})
}

// `Reject` writes a failure response with the given HTTP status code
// When `msg` is non-empty it is wrapped inside `HttpRes[any]`
// When `msg` is empty only the status code is written without a JSON body
func Reject(cx iris.Context, code int, msg string) {
	cx.StatusCode(code)

	if msg != "" {
		_ = cx.JSON(HttpRes[any]{
			Code: code,
			Msg:  msg,
		})
	}
}
