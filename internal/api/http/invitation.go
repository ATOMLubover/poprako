package http

import (
	"labelplus-next-web-be/internal/state"
	"labelplus-next-web-be/internal/value"

	"github.com/kataras/iris/v12"
)

// ListInvitations godoc
// @Summary 	获取邀请列表（已测试）
// @Description 获取指定汉化组的邀请列表，注意当列表为空，会返回 null 而不是空数组
//
// @Tags 		invitation
// @Security 	ApiKeyAuth
// @Produce 	json
// @Param 		team_id query string true "汉化组 ID"
//
// @Success 	200 {object} []value.InvitationInfo
//
// @Router 		/invitations [get]
func ListInvitations(appState *state.AppState) iris.Handler {
	invitationApplication := appState.InvitationApplication

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrentUserID(ctx)
		if !ok {
			return
		}

		// 从 query param 中获取 team_id 参数
		teamID := ctx.URLParam("team_id")
		if teamID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 team_id 查询参数")
			return
		}

		result, err := invitationApplication.ListInvitations(
			*buildTraceScope(ctx),
			currentUserID,
			teamID,
		)
		if err != nil {
			reject(ctx, iris.StatusForbidden, err.Error())
			return
		}

		accept(ctx, "获取邀请列表成功", result)
	}
}

// CreateInvitation godoc
// @Summary 	创建邀请（已测试）
// @Description 在指定汉化组中创建一个新的邀请
//
// @Tags 		invitation
// @Security 	ApiKeyAuth
// @Accept 		json
// @Produce 	json
// @Param 		body body value.CreateInvitationArgs true "创建邀请参数"
//
// @Success 	201 {object} value.InvitationInfo
//
// @Router 		/invitations [post]
func CreateInvitation(appState *state.AppState) iris.Handler {
	invitationApplication := appState.InvitationApplication

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrentUserID(ctx)
		if !ok {
			return
		}

		var args value.CreateInvitationArgs

		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误: "+err.Error())
			return
		}

		result, err := invitationApplication.CreateInvitation(
			*buildTraceScope(ctx),
			currentUserID,
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
// @Param 		body body value.UpdateInvitationArgs true "更新邀请参数"
//
// @Success 	200
//
// @Router 		/invitations/{invitation_id} [patch]
func PatchInvitation(appState *state.AppState) iris.Handler {
	invitationApplication := appState.InvitationApplication

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrentUserID(ctx)
		if !ok {
			return
		}

		var args value.UpdateInvitationArgs

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

		if err := invitationApplication.UpdateInvitation(
			*buildTraceScope(ctx),
			currentUserID,
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
	invitationApplication := appState.InvitationApplication

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrentUserID(ctx)
		if !ok {
			return
		}

		// 从路径参数中获取 invitation_id
		invitationID := ctx.Params().Get("invitation_id")
		if invitationID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 invitation_id 路径参数")
			return
		}

		if err := invitationApplication.DeleteInvitation(*buildTraceScope(ctx), currentUserID, invitationID); err != nil {
			reject(ctx, iris.StatusForbidden, err.Error())
			return
		}

		accept(ctx, "删除邀请成功", nil)
	}
}
