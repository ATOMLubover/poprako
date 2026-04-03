package http

import (
	"poprako-s/internal/app/val"
	"poprako-s/internal/state"

	"github.com/kataras/iris/v12"
)

// CreateMember godoc
// @Summary 	创建成员
// @Description 由超级管理员直接创建成员记录
//
// @Tags 		member
// @Security 	ApiKeyAuth
// @Accept 		json
// @Produce 	json
// @Param 		body body val.CreateMemberArgs true "创建成员参数"
//
// @Success 	201 {object} val.CreateMemberRes
//
// @Router 		/members [post]
func CreateMember(appState *state.AppState) iris.Handler {
	memberApp := appState.MemberApp

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		var args val.CreateMemberArgs

		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误: "+err.Error())
			return
		}

		result, err := memberApp.Create(
			buildReqCx(ctx),
			currentUserID,
			&args,
		)
		if err != nil {
			reject(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		accept(ctx, "创建成员成功", result)
	}
}

// ListMembers godoc
// @Summary 	获取指定汉化组的成员列表
// @Description 获取指定汉化组的成员列表，注意当列表为空，会返回 null 而不是空数组
//
// @Tags 		member
// @Security 	ApiKeyAuth
// @Produce 	json
// @Param 		team_id query string true "汉化组 ID"
// @Param 		"includes[]" query []string false "include 关联信息，可选值：user"
// @Param 		offset query int true "偏移量"
// @Param 		limit query int true "每页数量"
//
// @Success 	200 {object} []val.MemberInfo
//
// @Router 		/members [get]
func ListMembers(appState *state.AppState) iris.Handler {
	memberApp := appState.MemberApp

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		var args val.ListTeamMemberArgs

		if err := ctx.ReadQuery(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "查询参数格式错误: "+err.Error())
			return
		}

		result, err := memberApp.ListByTeam(
			buildReqCx(ctx),
			currentUserID,
			&args,
		)
		if err != nil {
			reject(ctx, iris.StatusForbidden, err.Error())
			return
		}

		accept(ctx, "获取成员列表成功", result)
	}
}

// ListMyMembers godoc
// @Summary 	获取当前用户的成员身份列表
// @Description 获取当前用户在各汉化组中的成员信息，注意当列表为空，会返回 null 而不是空数组
//
// @Tags 		member
// @Security 	ApiKeyAuth
// @Produce 	json
// @Param 		"includes[]" query []string false "include 关联信息，可选值：team"
// @Param 		offset query int true "偏移量"
// @Param 		limit query int true "每页数量"
//
// @Success 	200 {object} []val.MemberInfo
//
// @Router 		/members/mine [get]
func ListMyMembers(appState *state.AppState) iris.Handler {
	memberApp := appState.MemberApp

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		var args val.ListMyMemberArgs

		if err := ctx.ReadQuery(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "查询参数格式错误: "+err.Error())
			return
		}

		result, err := memberApp.ListMy(buildReqCx(ctx), currentUserID, &args)
		if err != nil {
			reject(ctx, iris.StatusForbidden, err.Error())
			return
		}

		accept(ctx, "获取我的成员身份列表成功", result)
	}
}

// UpdateMemberRole godoc
// @Summary 	更新成员角色
// @Description 更新指定成员的分工角色
//
// @Tags 		member
// @Security 	ApiKeyAuth
// @Accept 		json
// @Produce 	json
// @Param 		member_id path string true "成员 ID"
// @Param 		body body val.UpdateMemberRoleArgs true "更新成员角色参数"
//
// @Success 	200
//
// @Router 		/members/{member_id} [put]
func UpdateMemberRole(appState *state.AppState) iris.Handler {
	memberApp := appState.MemberApp

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		memberID := ctx.Params().Get("member_id")
		if memberID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 member_id 路径参数")
			return
		}

		var args val.UpdateMemberRoleArgs

		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误: "+err.Error())
			return
		}

		if args.ID != memberID {
			reject(ctx, iris.StatusBadRequest, "路径参数 member_id 与请求体中的 ID 不匹配")
			return
		}

		if err := memberApp.UpdateRole(
			buildReqCx(ctx),
			currentUserID,
			&args,
		); err != nil {
			reject(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		accept(ctx, "更新成员角色成功", nil)
	}
}

// JoinTeam godoc
// @Summary 	通过邀请码加入汉化组
// @Description 已登录用户使用邀请码加入对应汉化组
//
// @Tags 		member
// @Security 	ApiKeyAuth
// @Accept 		json
// @Produce 	json
// @Param 		body body val.JoinTeamArgs true "加入汉化组参数"
//
// @Success 	200
//
// @Router 		/members/join [post]
func JoinTeam(appState *state.AppState) iris.Handler {
	memberApp := appState.MemberApp

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		var args val.JoinTeamArgs

		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误: "+err.Error())
			return
		}

		if err := memberApp.JoinTeam(
			buildReqCx(ctx),
			currentUserID,
			&args,
		); err != nil {
			reject(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		accept(ctx, "加入汉化组成功", nil)
	}
}

// RemoveMember godoc
// @Summary 	移除成员
// @Description 从汉化组中移除指定成员
//
// @Tags 		member
// @Security 	ApiKeyAuth
// @Produce 	json
// @Param 		member_id path string true "成员 ID"
//
// @Success 	200
//
// @Router 		/members/{member_id} [delete]
func RemoveMember(appState *state.AppState) iris.Handler {
	memberApp := appState.MemberApp

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		memberID := ctx.Params().Get("member_id")
		if memberID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 member_id 路径参数")
			return
		}

		if err := memberApp.Remove(
			buildReqCx(ctx),
			currentUserID,
			memberID,
		); err != nil {
			reject(ctx, iris.StatusForbidden, err.Error())
			return
		}

		accept(ctx, "移除成员成功", nil)
	}
}
