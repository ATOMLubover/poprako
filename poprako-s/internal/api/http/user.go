package http

import (
	"poprako-s/internal/app/val"
	"poprako-s/internal/state"

	"github.com/kataras/iris/v12"
)

// GetUserByID godoc
// @Summary 	根据 ID 获取用户信息
// @Description 根据用户 ID 获取用户详细信息
//
// @Tags 		user
// @Security 	ApiKeyAuth
// @Produce 	json
// @Param 		user_id path string true "用户 ID"
//
// @Success 	200 {object} val.UserInfo
//
// @Router 		/users/{user_id} [get]
func GetUserByID(appState *state.AppState) iris.Handler {
	userApp := appState.UserApp

	return func(ctx iris.Context) {
		userID := ctx.Params().Get("user_id")
		if userID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 user_id 路径参数")
			return
		}

		result, err := userApp.GetInfo(buildReqCx(ctx), userID)
		if err != nil {
			reject(ctx, iris.StatusInternalServerError, err.Error())
			return
		}

		accept(ctx, "获取用户信息成功", result)
	}
}

// GetMyUser godoc
// @Summary 	获取当前登录用户信息
// @Description 获取当前登录用户的详细信息，用于保持登录状态
//
// @Tags 		user
// @Security 	ApiKeyAuth
// @Produce 	json
//
// @Success 	200 {object} val.UserInfo
//
// @Router 		/users/mine [get]
func GetMyUser(appState *state.AppState) iris.Handler {
	userApp := appState.UserApp

	return func(ctx iris.Context) {
		currUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		result, err := userApp.GetMyInfo(buildReqCx(ctx), currUserID)
		if err != nil {
			reject(ctx, iris.StatusForbidden, err.Error())
			return
		}

		accept(ctx, "获取当前用户信息成功", result)
	}
}

// GetMyUserStats godoc
// @Summary 	获取当前登录用户统计信息
// @Description 获取当前登录用户的任务统计信息
//
// @Tags 		user
// @Security 	ApiKeyAuth
// @Produce 	json
//
// @Success 	200 {object} val.UserStatsInfo
//
// @Router 		/users/mine/stats [get]
func GetMyUserStats(appState *state.AppState) iris.Handler {
	userApp := appState.UserApp

	return func(ctx iris.Context) {
		currUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		result, err := userApp.GetMyStats(buildReqCx(ctx), currUserID)
		if err != nil {
			reject(ctx, iris.StatusForbidden, err.Error())
			return
		}

		accept(ctx, "获取当前用户统计信息成功", result)
	}
}

// ReserveMyAvatar godoc
// @Summary 	预留当前用户头像上传
// @Description 为当前用户头像生成预签名 PUT URL，并预留 avatar_oss_key
//
// @Tags 		user
// @Security 	ApiKeyAuth
// @Accept 		json
// @Produce 	json
// @Param 		body body val.ReserveUserAvatarArgs true "预留用户头像参数"
//
// @Success 	200 {object} val.ReserveUserAvatarRes
//
// @Router 		/users/mine/avatar [post]
func ReserveMyAvatar(appState *state.AppState) iris.Handler {
	userApp := appState.UserApp

	return func(ctx iris.Context) {
		currUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		var args val.ReserveUserAvatarArgs

		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误: "+err.Error())
			return
		}

		result, err := userApp.ReserveMyAvatar(buildReqCx(ctx), currUserID, &args)
		if err != nil {
			reject(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		accept(ctx, "预留用户头像成功", result)
	}
}

// ConfirmMyAvatarUploaded godoc
// @Summary 	确认当前用户头像已上传
// @Description 在客户端上传头像后，确认当前用户头像上传状态
//
// @Tags 		user
// @Security 	ApiKeyAuth
// @Produce 	json
//
// @Success 	200
//
// @Router 		/users/mine/avatar/confirm [post]
func ConfirmMyAvatarUploaded(appState *state.AppState) iris.Handler {
	userApp := appState.UserApp

	return func(ctx iris.Context) {
		currUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		if err := userApp.ConfirmMyAvatarUploaded(buildReqCx(ctx), currUserID); err != nil {
			reject(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		accept(ctx, "确认用户头像上传成功", nil)
	}
}

// UpdateMyUser godoc
// @Summary 	更新当前用户信息
// @Description 更新当前登录用户的基本资料
//
// @Tags 		user
// @Security 	ApiKeyAuth
// @Accept 		json
// @Produce 	json
// @Param 		body body val.UpdateUserArgs true "更新用户参数"
//
// @Success 	200
//
// @Router 		/users/mine [put]
func UpdateMyUser(appState *state.AppState) iris.Handler {
	userApp := appState.UserApp

	return func(ctx iris.Context) {
		currUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		var args val.UpdateUserArgs

		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误: "+err.Error())
			return
		}

		if err := userApp.UpdateMyInfo(buildReqCx(ctx), currUserID, &args); err != nil {
			reject(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		accept(ctx, "更新用户成功", nil)
	}
}

// RemoveUser godoc
// @Summary 	删除用户
// @Description 根据用户 ID 删除用户
//
// @Tags 		user
// @Security 	ApiKeyAuth
// @Produce 	json
// @Param 		user_id path string true "用户 ID"
//
// @Success 	200
//
// @Router 		/users/{user_id} [delete]
func RemoveUser(appState *state.AppState) iris.Handler {
	userApp := appState.UserApp

	return func(ctx iris.Context) {
		currUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		targetUserID := ctx.Params().Get("user_id")
		if targetUserID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 user_id 路径参数")
			return
		}

		if err := userApp.Remove(buildReqCx(ctx), currUserID, targetUserID); err != nil {
			reject(ctx, iris.StatusForbidden, err.Error())
			return
		}

		accept(ctx, "删除用户成功", nil)
	}
}
