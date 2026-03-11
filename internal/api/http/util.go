package http

import (
	"labelplus-next-web-be/internal/util"

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

	ctx.JSON(FormatResponse{
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

	ctx.JSON(FormatResponse{
		// 暂时使用统一的 200 HTTP status code，不携带具体业务状态码
		Code:    iris.StatusOK,
		Message: message,
		Data:    data,
	})
}

func buildTraceScope(ctx iris.Context) util.TraceScope {
	requestID := requestid.Get(ctx)

	return *util.NewTraceScope(zap.L()).
		WithFields(zap.String("request_id", requestID))
}

// extractCurrentUserID 从请求上下文中取出认证用户 ID
// 若未找到则直接返回 401，并返回 false，调用方应立即 return
func extractCurrentUserID(ctx iris.Context) (string, bool) {
	userID := ctx.Values().GetString("user_id")
	if userID == "" {
		reject(ctx, iris.StatusUnauthorized, "未授权，请先登录")
		return "", false
	}

	return userID, true
}
