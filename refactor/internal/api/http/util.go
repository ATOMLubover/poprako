package http

import (
	"context"

	"poprako-s/internal/api/http/middleware"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/pkg/util"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/requestid"
	"go.uber.org/zap"
)

// `takeCurrUid` extracts the current user id from the Iris context
// It reads the `UserToken` stored by `Auth` middleware under `UtkKey`
// Returns the user id and `true` on success, or an empty string and `false`
// when the token is missing or invalid
func takeCurrUid(cx iris.Context) (string, bool) {
	utk, ok := cx.Values().Get(middleware.UtkKey).(*aggr.UserToken)
	if !ok || utk == nil {
		return "", false
	}

	return utk.UserId, true
}

// `newReqCx` creates a request-scoped `context.Context` from an Iris context
// It injects the request id (from middleware or a newly generated one) into
// a zap logger attachment, which `app` layer constructors consume via `app_util`
func newReqCx(cx iris.Context) context.Context {
	reqId := requestid.Get(cx)
	if reqId == "" {
		reqId = util.GenId("request")
	}

	lgr := zap.L().With(zap.String("request_id", reqId))

	return app_util.SaveLgr(cx.Request().Context(), lgr)
}
