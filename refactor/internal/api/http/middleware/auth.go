package middleware

import (
	"poprako-s/internal/api/http/res"
	token_iface "poprako-s/internal/domain/ext/token"
	token_impl "poprako-s/internal/infra/ext/token"

	"github.com/kataras/iris/v12"
)

func Auth() iris.Handler {
	parser := token_impl.NewJwtParser()

	return func(cx iris.Context) {
		authHeader := cx.GetHeader("Authorization")
		authCookie := cx.GetCookie("authorization")

		if authHeader == "" && authCookie == "" {
			res.Reject(cx, 401, "缺少授权信息")
			return
		}

		var tk string

		// NOTE: Cookie is preferred over header.
		if authCookie != "" {
			tk = authCookie
		} else {
			tk = authHeader
		}

		utk, err := parser.ParseToken(token_iface.SignedToken(tk))
		if err != nil {
			res.Reject(cx, 401, "无效的授权信息")
			return
		}

		cx.Values().Set("user_token", utk)

		cx.Next()
	}
}
