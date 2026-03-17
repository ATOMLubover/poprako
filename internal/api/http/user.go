package http

import (
	"labelplus-next-web-be/internal/state"
	"labelplus-next-web-be/internal/value"

	"github.com/kataras/iris/v12"
)

// GetUserByID godoc
// @Summary 	根据 ID 获取用户信息（已测试）
// @Description 根据用户 ID 获取用户详细信息
//
// @Tags 		user
// @Security 	ApiKeyAuth
// @Accept 		json
// @Produce 	json
// @Param 		user_id path string true "用户 ID"
// @Param 		body body value.ReserveUserAvatarArgs true "预留用户头像上传参数"
//
// @Success 	200 {object} value.UserInfo
//
// @Router 		/users/{user_id} [get]
func GetUserByID(appState *state.AppState) iris.Handler {
	userApplication := appState.UserApplication

	return func(ctx iris.Context) {
		// 从 URL 路径参数中获取用户 ID
		userID := ctx.Params().Get("user_id")
		if userID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 user_id 路径参数")
			return
		}

		var args value.GetUserArgs
		if err := ctx.ReadQuery(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "查询参数格式错误: "+err.Error())
			return
		}

		// 调用 UserApplication 获取用户信息
		result, err := userApplication.GetUser(buildTraceScope(ctx), userID, args)
		if err != nil {
			reject(ctx, iris.StatusInternalServerError, err.Error())
			return
		}

		// 返回成功响应
		accept(ctx, "获取用户信息成功", result)
	}
}

/* // ListUsers godoc
// @Summary 	获取用户列表（已测试）
// @Description 根据查询条件获取用户列表，支持按 QQ、模糊名称筛选，注意当列表为空，会返回 null 而不是空数组
//
// @Tags 		user
// @Security 	ApiKeyAuth
// @Produce 	json
// @Param 		qq query string false "QQ 号"
// @Param 		fuzzy_name query string false "模糊名称"
// @Param 		offset query int true "偏移量"
// @Param 		limit query int true "每页数量"
//
// @Success 	200 {object} []value.UserInfo
//
// @Router 		/users [get]
func ListUsers(appState *state.AppState) iris.Handler {
	userApplication := appState.UserApplication

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrentUserID(ctx)
		if !ok {
			return
		}

		var args value.ListUserArgs

		if err := ctx.ReadQuery(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "查询参数格式错误: "+err.Error())
			return
		}

		result, err := userApplication.ListUsers(
			buildTraceScope(ctx),
			currentUserID,
			args,
		)
		if err != nil {
			reject(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		accept(ctx, "获取用户列表成功", result)
	}
} */

// ReserveUserAvatar godoc
// @Summary 	预留用户头像上传
// @Description 为指定用户头像生成预签名 PUT URL，并预留 avatar_oss_key
//
// @Tags 		user
// @Security 	ApiKeyAuth
// @Produce 	json
// @Param 		user_id path string true "用户 ID"
//
// @Success 	200 {object} value.ReserveUserAvatarResult
//
// @Router 		/users/{user_id}/avatar [post]
func ReserveUserAvatar(appState *state.AppState) iris.Handler {
	userApplication := appState.UserApplication

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrentUserID(ctx)
		if !ok {
			return
		}

		targetUserID := ctx.Params().Get("user_id")
		if targetUserID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 user_id 路径参数")
			return
		}

		var args value.ReserveUserAvatarArgs
		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误: "+err.Error())
			return
		}

		result, err := userApplication.ReserveUserAvatar(
			buildTraceScope(ctx),
			currentUserID,
			targetUserID,
			args,
		)
		if err != nil {
			reject(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		accept(ctx, "预留用户头像成功", result)
	}
}

// ConfirmUserAvatarUploaded godoc
// @Summary 	确认用户头像已上传
// @Description 在客户端上传头像后，确认用户头像上传状态
//
// @Tags 		user
// @Security 	ApiKeyAuth
// @Produce 	json
// @Param 		user_id path string true "用户 ID"
//
// @Success 	200
//
// @Router 		/users/{user_id}/avatar/confirm [post]
func ConfirmUserAvatarUploaded(appState *state.AppState) iris.Handler {
	userApplication := appState.UserApplication

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrentUserID(ctx)
		if !ok {
			return
		}

		targetUserID := ctx.Params().Get("user_id")
		if targetUserID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 user_id 路径参数")
			return
		}

		if err := userApplication.ConfirmUserAvatarUploaded(
			buildTraceScope(ctx),
			currentUserID,
			targetUserID,
		); err != nil {
			reject(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		accept(ctx, "确认用户头像上传成功", nil)
	}
}

// UpdateUserByID godoc
// @Summary 	更新用户
// @Description 更新指定用户的信息
//
// @Tags 		user
// @Security 	ApiKeyAuth
// @Accept 		json
// @Produce 	json
// @Param 		user_id path string true "用户 ID"
// @Param 		body body value.UpdateUserArgs true "更新用户参数"
//
// @Success 	200
//
// @Router 		/users/{user_id} [put]
func UpdateUserByID(appState *state.AppState) iris.Handler {
	userApplication := appState.UserApplication

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrentUserID(ctx)
		if !ok {
			return
		}

		targetUserID := ctx.Params().Get("user_id")
		if targetUserID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 user_id 路径参数")
			return
		}

		var args value.UpdateUserArgs
		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误: "+err.Error())
			return
		}

		args.UserID = targetUserID

		if err := userApplication.UpdateUser(
			buildTraceScope(ctx),
			currentUserID,
			args,
		); err != nil {
			reject(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		accept(ctx, "更新用户成功", nil)
	}
}

// RemoveUserByID godoc
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
func RemoveUserByID(appState *state.AppState) iris.Handler {
	userApplication := appState.UserApplication

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrentUserID(ctx)
		if !ok {
			return
		}

		targetUserID := ctx.Params().Get("user_id")
		if targetUserID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 user_id 路径参数")
			return
		}

		if err := userApplication.RemoveUser(buildTraceScope(ctx), currentUserID, targetUserID); err != nil {
			reject(ctx, iris.StatusForbidden, err.Error())
			return
		}

		accept(ctx, "删除用户成功", nil)
	}
}

// GetMyUser godoc
// @Summary  获取当前登录用户信息
// @Description 获取当前登录用户的详细信息，用于保持登录状态
//
// @Tags    user
// @Security    ApiKeyAuth
// @Produce json
//
// @Success 200 {object} value.UserInfo
//
// @Router  /users/mine [get]
func GetMyUser(appState *state.AppState) iris.Handler {
	userApplication := appState.UserApplication

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrentUserID(ctx)
		if !ok {
			return
		}

		result, err := userApplication.GetMyUser(buildTraceScope(ctx), currentUserID)
		if err != nil {
			reject(ctx, iris.StatusForbidden, err.Error())
			return
		}

		accept(ctx, "获取当前用户信息成功", result)
	}
}
