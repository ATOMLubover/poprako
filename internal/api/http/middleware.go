package http

import (
	"strings"
	"time"

	"labelplus-next-web-be/internal/domain/service"
	"labelplus-next-web-be/internal/state"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/requestid"
	"go.uber.org/zap"
)

func LogMiddleware(appState *state.AppState) iris.Handler {
	isDevelopment := appState.AppConfig.IsDevelopment()

	return func(ctx iris.Context) {
		// 记录请求处理时间
		startTime := time.Now()

		ctx.Next()

		duration := time.Since(startTime)

		// 根据运行环境，收集请求信息
		if isDevelopment {
			method := ctx.Method()
			path := ctx.Path()
			statusCode := ctx.GetStatusCode()
			remoteAddr := ctx.RemoteAddr()
			requestID := requestid.Get(ctx) // 必须在 requestid 中间件之后使用

			zap.L().Debug(
				"HTTP 请求完成",
				zap.String("method", method),
				zap.String("path", path),
				zap.Int("status_code", statusCode),
				zap.String("remote_addr", remoteAddr),
				zap.String("request_id", requestID),
				zap.Duration("duration", duration),
			)
		}
	}
}

func AuthorizeMiddleware(appState *state.AppState) iris.Handler {
	secretKey := []byte(appState.AppConfig.AuthConfig.JWTSecretKey)

	return func(ctx iris.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			reject(ctx, iris.StatusUnauthorized, "未提供 Authorization 头部")
			return
		}

		parts := strings.Fields(authHeader)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			reject(ctx, iris.StatusUnauthorized, "Authorization 头部格式必须为 Bearer <token>")
			return
		}

		tokenString := parts[1]

		claims, err := service.ParseAccessToken(tokenString, secretKey)
		if err != nil {
			reject(ctx, iris.StatusUnauthorized, "无效的访问令牌")
			return
		}

		if claims == nil || claims.UserID == "" {
			reject(ctx, iris.StatusUnauthorized, "访问令牌不包含用户信息")
			return
		}

		ctx.Values().Set("user_id", claims.UserID)
		ctx.Next()
	}
}
