package http

import (
	"context"

	app_impl "poprako-s/internal/app/impl"
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/pkg/util"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/requestid"
	"go.uber.org/zap"
)

func takeCurrUid(cx iris.Context) (string, bool) {
	utk, ok := cx.Values().Get("").(*aggr.UserToken)
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

	return context.WithValue(context.Background(), app_impl.LgrKey, lgr)
}
