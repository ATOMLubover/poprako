package middleware

import (
	"strings"

	"poprako-s/internal/api/http/res"
	token_iface "poprako-s/internal/domain/ext/token"
	token_impl "poprako-s/internal/infra/ext/token"

	"github.com/kataras/iris/v12"
)

// `UtkKey` is the context key under which the parsed `UserToken` is stored
// Both `Auth` middleware (writer) and `takeCurrUid` (reader) use this key
const UtkKey = "user_token"

const AuthCookieName = "authorization"

// `Auth` is an Iris middleware that extracts and validates a JWT token
// It reads the token from the `authorization` cookie first, falling back
// to the `Authorization` header when the cookie is absent
// On success the parsed `UserToken` is stored in the Iris context under `UtkKey`
// On failure a 401 `res.HttpRes` rejection is written and the handler chain stops
func Auth() iris.Handler {
	parser := token_impl.NewJwtParser()

	return func(cx iris.Context) {
		authHeader := cx.GetHeader("Authorization")
		authCookie := cx.GetCookie("authorization")

		if authHeader == "" && authCookie == "" {
			res.Reject(cx, 401, "缺少授权信息，请尝试重新登录")
			return
		}

		var tk string

		// NOTE: Cookie is preferred over header.
		if authCookie != "" {
			tk = authCookie
		} else {
			tk = authHeader
		}

		// Strip "Bearer " prefix from Authorization header
		tk = strings.TrimPrefix(tk, "Bearer ")

		utk, err := parser.ParseToken(token_iface.SignedToken(tk))
		if err != nil {
			res.Reject(cx, 401, "无效的授权信息，请尝试重新登录")
			return
		}

		cx.Values().Set(UtkKey, utk)

		cx.Next()
	}
}
