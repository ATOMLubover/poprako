package http

import (
	"labelplus-next-web-be/internal/state"
	"labelplus-next-web-be/internal/value"

	"github.com/kataras/iris/v12"
)

// Login godoc
// @Summary 	用户登录（已测试）
// @Description 使用 QQ 和密码进行登录，成功返回访问令牌
//
// @Tags 		auth
// @Accept 		json
// @Produce 	json
// @Param 		body body value.LoginUserArgs true "登录参数"
//
// @Success 	200 {object} value.LoginUserResult
//
// @Router 		/auth/login [post]
func Login(appState *state.AppState) iris.Handler {
	userApplication := appState.UserApplication

	return func(ctx iris.Context) {
		var args value.LoginUserArgs
		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误: "+err.Error())
			return
		}

		result, err := userApplication.LoginUser(*buildTraceScope(ctx), &args)
		if err != nil {
			reject(ctx, iris.StatusUnauthorized, err.Error())
			return
		}

		accept(ctx, "登录成功", result)
	}
}

// Register godoc
// @Summary 	用户注册
// @Description 使用 QQ、密码、名字和邀请码进行注册，成功返回访问令牌
//
// @Tags 		auth
// @Accept 		json
// @Produce 	json
// @Param 		body body value.RegisterUserArgs true "注册参数"
//
// @Success 	200 {object} value.RegisterUserResult
//
// @Router 		/auth/register [post]
func Register(appState *state.AppState) iris.Handler {
	userApplication := appState.UserApplication

	return func(ctx iris.Context) {
		var args value.RegisterUserArgs
		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误: "+err.Error())
			return
		}

		result, err := userApplication.RegisterUser(*buildTraceScope(ctx), &args)
		if err != nil {
			reject(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		accept(ctx, "注册成功", result)
	}
}
