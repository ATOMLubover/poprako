package middleware

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/cors"
	"github.com/kataras/iris/v12/x/errors"
)

func CorsMiddleware() iris.Handler {
	return cors.New().
		HandleErrorFunc(func(ctx iris.Context, err error) {
			errors.FailedPrecondition.Err(ctx, err)
		}).
		ExtractOriginFunc(cors.StrictOriginExtractor).
		ReferrerPolicy(cors.NoReferrerWhenDowngrade).
		AllowOrigin("localhost:5173").
		Handler()
}
