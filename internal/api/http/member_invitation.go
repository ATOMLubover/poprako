package http

import (
	"poprako-s/internal/api/http/res"
	"poprako-s/internal/api/state"
	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/model/enum"

	"github.com/kataras/iris/v12"
)

// `ListMemberInvitations` godoc
// @Summary List Member Invitations
// @Description List invitations under one team
// @Description The caller must be member of the target team
// @Description Auth: `authorization` cookie is preferred over `Authorization` header when both are present
// @Tags member-invitation
// @Security ApiKeyAuth
// @Produce json
// @Param team_id query string true "team id"
// @Param includes query []string false "include related fields, optional: invitor, invitee"
// @Param pending query bool false "pending filter"
// @Param offset query int false "pagination offset"
// @Param limit query int false "pagination limit"
// @Success 200 {object} res.HttpRes[[]val.MemberInvVal]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /api/v1/member-invitations [get]
func ListMemberInvitations(st *state.AppState) iris.Handler {
	memberInvApp := st.MemberInvApp

	return func(cx iris.Context) {
		currUid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		var args val.ListMemberInvArgs
		if err := cx.ReadQuery(&args); err != nil {
			res.Reject(cx, iris.StatusBadRequest, "请求参数解析失败")
			return
		}

		if args.TeamId == "" {
			res.Reject(cx, iris.StatusBadRequest, "缺少 team_id 参数")
			return
		}

		if !parseMemberInvIncludes(args.Includes) {
			res.Reject(cx, iris.StatusBadRequest, "includes 参数格式错误")
			return
		}

		re := memberInvApp.List(newReqCx(cx), currUid, &args)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

// `parseMemberInvIncludes` validates include query values.
func parseMemberInvIncludes(includes []enum.MemberInvIncl) bool {
	if len(includes) == 0 {
		return true
	}

	for i := range includes {
		switch includes[i] {
		case enum.MemberInvInclInvitor, enum.MemberInvInclInvitee:
			continue

		default:
			return false
		}
	}

	return true
}

// `CreateMemberInvitation` godoc
// @Summary Create Member Invitation
// @Description Create one invitation under one team
// @Description The caller must be team admin
// @Description Auth: `authorization` cookie is preferred over `Authorization` header when both are present
// @Tags member-invitation
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param body body val.CreateMemberInvArgs true "create invitation args"
// @Success 201 {object} res.HttpRes[val.CreateMemberInvRes]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /api/v1/member-invitations [post]
func CreateMemberInvitation(st *state.AppState) iris.Handler {
	memberInvApp := st.MemberInvApp

	return func(cx iris.Context) {
		currUid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		var args val.CreateMemberInvArgs
		if err := cx.ReadJSON(&args); err != nil {
			res.Reject(cx, iris.StatusBadRequest, "请求参数解析失败")
			return
		}

		re := memberInvApp.Create(newReqCx(cx), currUid, &args)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusCreated, re.Data())
	}
}

// `UpdateMemberInvitation` godoc
// @Summary Update Member Invitation
// @Description Update invitation role mask by put semantics
// @Description The caller must be team admin
// @Description Auth: `authorization` cookie is preferred over `Authorization` header when both are present
// @Tags member-invitation
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param invitation_id path string true "invitation id"
// @Param body body val.MemberInvUpdArgs true "update invitation args"
// @Success 200 {object} res.HttpRes[any]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /api/v1/member-invitations/{invitation_id} [put]
func UpdateMemberInvitation(st *state.AppState) iris.Handler {
	memberInvApp := st.MemberInvApp

	return func(cx iris.Context) {
		currUid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		invId := cx.Params().Get("invitation_id")
		if invId == "" {
			res.Reject(cx, iris.StatusBadRequest, "缺少 invitation_id 参数")
			return
		}

		var args val.MemberInvUpdArgs
		if err := cx.ReadJSON(&args); err != nil {
			res.Reject(cx, iris.StatusBadRequest, "请求参数解析失败")
			return
		}

		args.Id = invId

		re := memberInvApp.Update(newReqCx(cx), currUid, &args)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

// `DeleteMemberInvitation` godoc
// @Summary Delete Member Invitation
// @Description Hard delete one invitation by id
// @Description The caller must be team admin
// @Description Auth: `authorization` cookie is preferred over `Authorization` header when both are present
// @Tags member-invitation
// @Security ApiKeyAuth
// @Produce json
// @Param invitation_id path string true "invitation id"
// @Success 200 {object} res.HttpRes[any]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /api/v1/member-invitations/{invitation_id} [delete]
func DeleteMemberInvitation(st *state.AppState) iris.Handler {
	memberInvApp := st.MemberInvApp

	return func(cx iris.Context) {
		currUid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		invId := cx.Params().Get("invitation_id")
		if invId == "" {
			res.Reject(cx, iris.StatusBadRequest, "缺少 invitation_id 参数")
			return
		}

		re := memberInvApp.Delete(newReqCx(cx), currUid, invId)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}
