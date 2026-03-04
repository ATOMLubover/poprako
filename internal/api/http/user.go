package http

import (
	"labelplus-next-web-be/internal/state"
	"labelplus-next-web-be/internal/value"

	"github.com/kataras/iris/v12"
)

// GetUserByID godoc
// @Summary 	根据 ID 获取用户信息
// @Description 根据用户 ID 获取用户详细信息
//
// @Tags 		user
// @Produce 	json
// @Param 		user_id path string true "用户 ID"
//
// @Success 	200 {object} value.UserInfo
//
// @Router 		/api/v1/users/{user_id} [get]
func GetUserByID(appState *state.AppState) iris.Handler {
	userApplication := appState.UserApplication

	return func(ctx iris.Context) {
		// 从 URL 路径参数中获取用户 ID
		userID := ctx.Params().Get("user_id")
		if userID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 user_id 路径参数")
			return
		}

		// 调用 UserApplication 获取用户信息
		result, err := userApplication.GetUser(*buildTraceScope(ctx), userID)
		if err != nil {
			reject(ctx, iris.StatusInternalServerError, err.Error())
			return
		}

		// 返回成功响应
		accept(ctx, "获取用户信息成功", result)
	}
}

// ListUsers godoc
// @Summary 	获取用户列表
// @Description 根据查询条件获取用户列表，支持按 QQ、模糊名称筛选，注意当列表为空，会返回 null 而不是空数组
//
// @Tags 		user
// @Produce 	json
// @Param 		qq query string false "QQ 号"
// @Param 		fuzzy_name query string false "模糊名称"
// @Param 		offset query int false "偏移量，默认值为 0"
// @Param 		limit query int false "每页数量，默认值为 10"
//
// @Success 	200 {object} []value.UserInfo
//
// @Router 		/api/v1/users [get]
func ListUsers(appState *state.AppState) iris.Handler {
	userApplication := appState.UserApplication

	return func(ctx iris.Context) {
		var args value.ListUserArgs

		if err := ctx.ReadQuery(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "查询参数格式错误: "+err.Error())
			return
		}

		result, err := userApplication.ListUsers(*buildTraceScope(ctx), &args)
		if err != nil {
			reject(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		accept(ctx, "获取用户列表成功", result)
	}
}

// RemoveUserByID godoc
// @Summary 	删除用户
// @Description 根据用户 ID 删除用户
//
// @Tags 		user
// @Produce 	json
// @Param 		user_id path string true "用户 ID"
//
// @Success 	200
//
// @Router 		/api/v1/users/{user_id} [delete]
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

		if err := userApplication.RemoveUserByID(*buildTraceScope(ctx), currentUserID, targetUserID); err != nil {
			reject(ctx, iris.StatusForbidden, err.Error())
			return
		}

		accept(ctx, "删除用户成功", nil)
	}
}
