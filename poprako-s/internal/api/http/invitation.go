package http

import (
	"poprako-s/internal/app/val"
	"poprako-s/internal/state"

	"github.com/kataras/iris/v12"
)

// ListInvitations godoc
// @Summary 	获取邀请列表
// @Description 获取指定汉化组的邀请列表，注意当列表为空，会返回 null 而不是空数组
//
// @Tags 		invitation
// @Security 	ApiKeyAuth
// @Produce 	json
// @Param 		team_id query string true "汉化组 ID"
// @Param 		offset query int true "偏移量"
// @Param 		limit query int true "每页数量"
// @Param 		"includes[]" query []string false "include 关联信息，可选值：invitor"
//
// @Success 	200 {object} []val.InvitationInfo
//
// @Router 		/invitations [get]
func ListInvitations(appState *state.AppState) iris.Handler {
	invitationApp := appState.InvitationApp

	return func(ctx iris.Context) {
		currUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		var args val.ListTeamInvitationArgs

		if err := ctx.ReadQuery(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "查询参数格式错误: "+err.Error())
			return
		}

		result, err := invitationApp.List(
			buildReqCx(ctx),
			currUserID,
			&args,
		)
		if err != nil {
			reject(ctx, iris.StatusForbidden, err.Error())
			return
		}

		accept(ctx, "获取邀请列表成功", result)
	}
}

// CreateInvitation godoc
// @Summary 	创建邀请
// @Description 在指定汉化组中创建一个新的邀请
//
// @Tags 		invitation
// @Security 	ApiKeyAuth
// @Accept 		json
// @Produce 	json
// @Param 		body body val.CreateInvitationArgs true "创建邀请参数"
//
// @Success 	201 {object} val.InvitationInfo
//
// @Router 		/invitations [post]
func CreateInvitation(appState *state.AppState) iris.Handler {
	invitationApp := appState.InvitationApp

	return func(ctx iris.Context) {
		currUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		var args val.CreateInvitationArgs

		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误: "+err.Error())
			return
		}

		result, err := invitationApp.Create(
			buildReqCx(ctx),
			currUserID,
			&args,
		)
		if err != nil {
			reject(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		ctx.StatusCode(iris.StatusCreated)
		accept(ctx, "创建邀请成功", result)
	}
}

// PatchInvitation godoc
// @Summary 	更新未被使用的邀请
// @Description 更新指定待处理邀请的信息，无法更新已被使用或已失效的邀请
//
// @Tags 		invitation
// @Security 	ApiKeyAuth
// @Accept 		json
// @Produce 	json
// @Param 		invitation_id path string true "邀请 ID"
// @Param 		body body val.UpdateInvitationArgs true "更新邀请参数"
//
// @Success 	200
//
// @Router 		/invitations/{invitation_id} [put]
func PatchInvitation(appState *state.AppState) iris.Handler {
	invitationApp := appState.InvitationApp

	return func(ctx iris.Context) {
		currUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		var args val.UpdateInvitationArgs

		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误: "+err.Error())
			return
		}

		// 从路径参数中获取 invitation_id
		invitationID := ctx.Params().Get("invitation_id")
		if invitationID == "" || args.ID != invitationID {
			reject(ctx, iris.StatusBadRequest, "缺少 invitation_id 路径参数，或与请求体中的 ID 不匹配")
			return
		}

		if err := invitationApp.Update(
			buildReqCx(ctx),
			currUserID,
			&args,
		); err != nil {
			reject(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		accept(ctx, "更新邀请成功", nil)
	}
}

// DeleteInvitation godoc
// @Summary 	删除邀请
// @Description 删除指定的邀请
//
// @Tags 		invitation
// @Security 	ApiKeyAuth
// @Produce 	json
// @Param 		invitation_id path string true "邀请 ID"
//
// @Success 	200
//
// @Router 		/invitations/{invitation_id} [delete]
func DeleteInvitation(appState *state.AppState) iris.Handler {
	invitationApp := appState.InvitationApp

	return func(ctx iris.Context) {
		currUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		// 从路径参数中获取 invitation_id
		invitationID := ctx.Params().Get("invitation_id")
		if invitationID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 invitation_id 路径参数")
			return
		}

		if err := invitationApp.Remove(buildReqCx(ctx), currUserID, invitationID); err != nil {
			reject(ctx, iris.StatusForbidden, err.Error())
			return
		}

		accept(ctx, "删除邀请成功", nil)
	}
}
