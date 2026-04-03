package http

import (
	"context"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/requestid"
	"go.uber.org/zap"
)

func reject(
	ctx iris.Context,
	code int,
	message string,
) {
	ctx.StatusCode(code)

	_ = ctx.JSON(FormatResponse{
		Code:    code,
		Message: message,
	})
}

func accept(
	ctx iris.Context,
	message string,
	data any,
) {
	// 考虑到兼容性（如中间件挤占 status code）
	// 统一使用 200 状态码，错误信息通过 code 字段传递
	ctx.StatusCode(iris.StatusOK)

	_ = ctx.JSON(FormatResponse{
		// 暂时使用统一的 200 HTTP status code，不携带具体业务状态码
		Code:    iris.StatusOK,
		Message: message,
		Data:    data,
	})
}

// buildReqCx 从 HTTP 请求上下文构建一个新的 context.Context，并注入请求 ID 以供后续处理使用
func buildReqCx(ctx iris.Context) context.Context {
	requestID := requestid.Get(ctx)

	lgr := zap.L().With(zap.String("request_id", requestID))

	return context.WithValue(context.Background(), "lgr", lgr)
}

// extractCurrUserID 从请求上下文中取出认证用户 ID
// 若未找到则直接返回 401，并返回 false，调用方应立即 return
func extractCurrUserID(ctx iris.Context) (string, bool) {
	userID := ctx.Values().GetString("user_id")
	if userID == "" {
		reject(ctx, iris.StatusUnauthorized, "未授权，请先登录")
		return "", false
	}

	return userID, true
}
