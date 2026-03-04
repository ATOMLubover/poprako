package http

import (
	"labelplus-next-web-be/internal/state"
	"labelplus-next-web-be/internal/value"

	"github.com/kataras/iris/v12"
)

// ListMembers godoc
// @Summary 	获取指定汉化组的成员列表
// @Description 获取指定汉化组的成员列表，注意当列表为空，会返回 null 而不是空数组
//
// @Tags 		member
// @Produce 	json
// @Param 		team_id query string true "汉化组 ID"
// @Param 		offset query int false "偏移量，默认值为 0"
// @Param 		limit query int false "每页数量，默认值为 10"
//
// @Success 	200 {object} []value.MemberProfile
//
// @Router 		/api/v1/members [get]
func ListMembers(appState *state.AppState) iris.Handler {
	memberApplication := appState.MemberApplication

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrentUserID(ctx)
		if !ok {
			return
		}

		teamID := ctx.URLParam("team_id")
		if teamID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 team_id 查询参数")
			return
		}

		var paginationParams value.PaginationParams
		if err := ctx.ReadQuery(&paginationParams); err != nil {
			reject(ctx, iris.StatusBadRequest, "查询参数格式错误: "+err.Error())
			return
		}

		result, err := memberApplication.ListMembers(
			*buildTraceScope(ctx),
			currentUserID,
			teamID,
			paginationParams,
		)
		if err != nil {
			reject(ctx, iris.StatusForbidden, err.Error())
			return
		}

		accept(ctx, "获取成员列表成功", result)
	}
}

// UpdateMemberRole godoc
// @Summary 	更新成员角色
// @Description 更新指定成员的分工角色
//
// @Tags 		member
// @Accept 		json
// @Produce 	json
// @Param 		member_id path string true "成员 ID"
// @Param 		body body value.UpdateMemberRoleArgs true "更新成员角色参数"
//
// @Success 	200
//
// @Router 		/api/v1/members/{member_id} [patch]
func UpdateMemberRole(appState *state.AppState) iris.Handler {
	memberApplication := appState.MemberApplication

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrentUserID(ctx)
		if !ok {
			return
		}

		memberID := ctx.Params().Get("member_id")
		if memberID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 member_id 路径参数")
			return
		}

		var args value.UpdateMemberRoleArgs
		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误: "+err.Error())
			return
		}

		if args.ID != memberID {
			reject(ctx, iris.StatusBadRequest, "路径参数 member_id 与请求体中的 ID 不匹配")
			return
		}

		if err := memberApplication.UpdateMemberRole(
			*buildTraceScope(ctx),
			currentUserID,
			&args,
		); err != nil {
			reject(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		accept(ctx, "更新成员角色成功", nil)
	}
}

// RemoveMember godoc
// @Summary 	移除成员
// @Description 从汉化组中移除指定成员
//
// @Tags 		member
// @Produce 	json
// @Param 		member_id path string true "成员 ID"
//
// @Success 	200
//
// @Router 		/api/v1/members/{member_id} [delete]
func RemoveMember(appState *state.AppState) iris.Handler {
	memberApplication := appState.MemberApplication

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrentUserID(ctx)
		if !ok {
			return
		}

		memberID := ctx.Params().Get("member_id")
		if memberID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 member_id 路径参数")
			return
		}

		if err := memberApplication.RemoveMember(
			*buildTraceScope(ctx),
			currentUserID,
			memberID,
		); err != nil {
			reject(ctx, iris.StatusForbidden, err.Error())
			return
		}

		accept(ctx, "移除成员成功", nil)
	}
}
