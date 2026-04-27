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

func takeCurrUid(cx iris.Context) (string, bool) {
	utk, ok := cx.Values().Get(middleware.UtkKey).(*aggr.UserToken)
	if !ok || utk == nil {
		return "", false
	}

	return utk.UserId, true
}

func newReqCx(cx iris.Context) context.Context {
	reqId := requestid.Get(cx)
	if reqId == "" {
		reqId = util.GenId("request")
	}

	lgr := zap.L().With(zap.String("request_id", reqId))

	return app_util.SaveLgr(cx.Request().Context(), lgr)
}
