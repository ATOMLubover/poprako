package http

import (
	"poprako-s/internal/app/val"
	"poprako-s/internal/state"

	"github.com/kataras/iris/v12"
)

// Login godoc
// @Summary 	用户登录
// @Description 使用 QQ 和密码进行登录，成功返回访问令牌
//
// @Tags 		auth
// @Accept 		json
// @Produce 	json
// @Param 		body body val.LoginUserArgs true "登录参数"
//
// @Success 	200 {object} val.LoginUserRes
//
// @Router 		/auth/login [post]
func Login(appState *state.AppState) iris.Handler {
	userApp := appState.UserApp

	return func(ctx iris.Context) {
		var args val.LoginUserArgs
		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误: "+err.Error())
			return
		}

		result, err := userApp.Login(buildReqCx(ctx), &args)
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
// @Param 		body body val.RegUserArgs true "注册参数"
//
// @Success 	200 {object} val.RegUserRes
//
// @Router 		/auth/register [post]
func Register(appState *state.AppState) iris.Handler {
	userApp := appState.UserApp

	return func(ctx iris.Context) {
		var args val.RegUserArgs
		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误: "+err.Error())
			return
		}

		result, err := userApp.Reg(buildReqCx(ctx), &args)
		if err != nil {
			reject(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		accept(ctx, "注册成功", result)
	}
}
