package http

import (
	"poprako-s/internal/app/val"
	"poprako-s/internal/state"

	"github.com/kataras/iris/v12"
)

// CreateTeam godoc
// @Summary 	创建汉化组
// @Description 创建一个新的汉化组，仅超级管理员有权限
//
// @Tags 		team
// @Security 	ApiKeyAuth
// @Accept 		json
// @Produce 	json
// @Param 		body body val.CreateTeamArgs true "创建汉化组参数"
//
// @Success 	201 {object} val.CreateTeamRes
//
// @Router 		/teams [post]
func CreateTeam(appState *state.AppState) iris.Handler {
	teamApp := appState.TeamApp

	return func(ctx iris.Context) {
		currUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		var args val.CreateTeamArgs

		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误: "+err.Error())
			return
		}

		result, err := teamApp.Create(buildReqCx(ctx), currUserID, &args)
		if err != nil {
			reject(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		ctx.StatusCode(iris.StatusCreated)
		accept(ctx, "创建汉化组成功", result)
	}
}

// ListTeams godoc
// @Summary 	获取所有汉化组列表
// @Description 获取所有汉化组列表，仅超级管理员有权限，注意当列表为空，会返回 null 而不是空数组
//
// @Tags 		team
// @Security 	ApiKeyAuth
// @Produce 	json
// @Param 		offset query int true "偏移量"
// @Param 		limit query int true "每页数量"
//
// @Success 	200 {object} []val.TeamInfo
//
// @Router 		/teams [get]
func ListTeams(appState *state.AppState) iris.Handler {
	teamApp := appState.TeamApp

	return func(ctx iris.Context) {
		currUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		var args val.ListTeamArgs

		if err := ctx.ReadQuery(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "查询参数格式错误: "+err.Error())
			return
		}

		result, err := teamApp.List(buildReqCx(ctx), currUserID, &args)
		if err != nil {
			reject(ctx, iris.StatusForbidden, err.Error())
			return
		}

		accept(ctx, "获取所有汉化组列表成功", result)
	}
}

// ListMyTeams godoc
// @Summary 	获取当前用户所在的汉化组列表
// @Description 获取当前用户所在的汉化组列表，注意当列表为空，会返回 null 而不是空数组
//
// @Tags 		team
// @Security 	ApiKeyAuth
// @Produce 	json
// @Param 		offset query int true "偏移量"
// @Param 		limit query int true "每页数量"
//
// @Success 	200 {object} []val.TeamInfo
//
// @Router 		/teams/mine [get]
func ListMyTeams(appState *state.AppState) iris.Handler {
	teamApp := appState.TeamApp

	return func(ctx iris.Context) {
		currUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		var args val.ListMyTeamArgs

		if err := ctx.ReadQuery(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "查询参数格式错误: "+err.Error())
			return
		}

		result, err := teamApp.ListMy(buildReqCx(ctx), currUserID, &args)
		if err != nil {
			reject(ctx, iris.StatusForbidden, err.Error())
			return
		}

		accept(ctx, "获取我的汉化组列表成功", result)
	}
}

// UpdateTeam godoc
// @Summary 	更新汉化组信息
// @Description 更新指定汉化组的信息，超级管理员或汉化组管理员有权限
//
// @Tags 		team
// @Security 	ApiKeyAuth
// @Accept 		json
// @Produce 	json
// @Param 		team_id path string true "汉化组 ID"
// @Param 		body body val.UpdateTeamArgs true "更新汉化组参数"
//
// @Success 	200
//
// @Router 		/teams/{team_id} [put]
func UpdateTeam(appState *state.AppState) iris.Handler {
	teamApp := appState.TeamApp

	return func(ctx iris.Context) {
		currUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		teamID := ctx.Params().Get("team_id")
		if teamID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 team_id 路径参数")
			return
		}

		var args val.UpdateTeamArgs

		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误: "+err.Error())
			return
		}

		args.ID = teamID

		if err := teamApp.Update(buildReqCx(ctx), currUserID, &args); err != nil {
			reject(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		accept(ctx, "更新汉化组成功", nil)
	}
}

// ReserveTeamAvatar godoc
// @Summary 	预留汉化组头像上传
// @Description 为指定汉化组头像生成预签名 PUT URL，并预留 avatar_oss_key
//
// @Tags 		team
// @Security 	ApiKeyAuth
// @Produce 	json
// @Param 		team_id path string true "汉化组 ID"
//
// @Success 	200 {object} val.ReserveTeamAvatarRes
//
// @Router 		/teams/{team_id}/avatar [post]
func ReserveTeamAvatar(appState *state.AppState) iris.Handler {
	teamApp := appState.TeamApp

	return func(ctx iris.Context) {
		currUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		teamID := ctx.Params().Get("team_id")
		if teamID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 team_id 路径参数")
			return
		}

		result, err := teamApp.ReserveAvatar(buildReqCx(ctx), currUserID, teamID)
		if err != nil {
			reject(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		accept(ctx, "预留汉化组头像成功", result)
	}
}

// ConfirmTeamAvatarUploaded godoc
// @Summary 	确认汉化组头像已上传
// @Description 在客户端上传头像后，确认汉化组头像上传状态
//
// @Tags 		team
// @Security 	ApiKeyAuth
// @Produce 	json
// @Param 		team_id path string true "汉化组 ID"
//
// @Success 	200
//
// @Router 		/teams/{team_id}/avatar/confirm [post]
func ConfirmTeamAvatarUploaded(appState *state.AppState) iris.Handler {
	teamApp := appState.TeamApp

	return func(ctx iris.Context) {
		currUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		teamID := ctx.Params().Get("team_id")
		if teamID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 team_id 路径参数")
			return
		}

		if err := teamApp.ConfirmAvatarUploaded(buildReqCx(ctx), currUserID, teamID); err != nil {
			reject(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		accept(ctx, "确认汉化组头像上传成功", nil)
	}
}

// DeleteTeam godoc
// @Summary 	删除汉化组
// @Description 删除指定汉化组，超级管理员或汉化组管理员有权限
//
// @Tags 		team
// @Security 	ApiKeyAuth
// @Produce 	json
// @Param 		team_id path string true "汉化组 ID"
//
// @Success 	200
//
// @Router 		/teams/{team_id} [delete]
func DeleteTeam(appState *state.AppState) iris.Handler {
	teamApp := appState.TeamApp

	return func(ctx iris.Context) {
		currUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		teamID := ctx.Params().Get("team_id")
		if teamID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 team_id 路径参数")
			return
		}

		if err := teamApp.Remove(buildReqCx(ctx), currUserID, teamID); err != nil {
			reject(ctx, iris.StatusForbidden, err.Error())
			return
		}

		accept(ctx, "删除汉化组成功", nil)
	}
}
