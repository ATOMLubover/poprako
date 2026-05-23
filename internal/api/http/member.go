package http

import (
	"poprako-s/internal/api/http/res"
	"poprako-s/internal/api/state"
	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/model/enum"

	"github.com/kataras/iris/v12"
)

// `CreateMember` godoc
// @Summary Create Member
// @Description Create one member under one team
// @Description The caller must be team admin
// @Description Auth: `authorization` cookie is preferred over `Authorization` header when both are present
// @Tags member
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param body body val.CreateMemberArgs true "create member args"
// @Success 201 {object} res.HttpRes[val.CreateMemberRes]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /api/v1/members [post]
func CreateMember(st *state.AppState) iris.Handler {
	memberApp := st.MemberApp

	return func(cx iris.Context) {
		currUid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		var args val.CreateMemberArgs
		if err := cx.ReadJSON(&args); err != nil {
			res.Reject(cx, iris.StatusBadRequest, "请求参数解析失败")
			return
		}

		re := memberApp.Create(newReqCx(cx), currUid, &args)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusCreated, re.Data())
	}
}

// `ListTeamMembers` godoc
// @Summary List Team Members
// @Description List members under one team
// @Description The caller must be a member of the target team
// @Description Auth: `authorization` cookie is preferred over `Authorization` header when both are present
// @Tags member
// @Security ApiKeyAuth
// @Produce json
// @Param team_id query string true "team id"
// @Param includes query []string false "include related fields, optional: user, team"
// @Param user_nickname_keyword query string false "fuzzy keyword for member user nickname"
// @Param role query int false "single role mask value, at most one bit set"
// @Param offset query int false "pagination offset"
// @Param limit query int false "pagination limit"
// @Success 200 {object} res.HttpRes[[]val.MemberVal]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /api/v1/members [get]
func ListTeamMembers(st *state.AppState) iris.Handler {
	memberApp := st.MemberApp

	return func(cx iris.Context) {
		currUid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		var args val.ListMemberByTeamArgs
		if err := cx.ReadQuery(&args); err != nil {
			res.Reject(cx, iris.StatusBadRequest, "请求参数解析失败")
			return
		}

		if args.TeamId == "" {
			res.Reject(cx, iris.StatusBadRequest, "缺少 team_id 参数")
			return
		}

		if !parseMemberIncludes(args.Includes) {
			res.Reject(cx, iris.StatusBadRequest, "includes 参数格式错误")
			return
		}

		re := memberApp.ListByTeam(newReqCx(cx), currUid, &args)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

// `ListMyMembers` godoc
// @Summary List My Members
// @Description List all memberships of current user
// @Description Auth: `authorization` cookie is preferred over `Authorization` header when both are present
// @Tags member
// @Security ApiKeyAuth
// @Produce json
// @Param includes query []string false "include related fields, optional: user, team"
// @Param offset query int false "pagination offset"
// @Param limit query int false "pagination limit"
// @Success 200 {object} res.HttpRes[[]val.MemberVal]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /api/v1/members/mine [get]
func ListMyMembers(st *state.AppState) iris.Handler {
	memberApp := st.MemberApp

	return func(cx iris.Context) {
		currUid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		var args val.ListMyMemberArgs
		if err := cx.ReadQuery(&args); err != nil {
			res.Reject(cx, iris.StatusBadRequest, "请求参数解析失败")
			return
		}

		if !parseMemberIncludes(args.Includes) {
			res.Reject(cx, iris.StatusBadRequest, "includes 参数格式错误")
			return
		}

		re := memberApp.ListMine(newReqCx(cx), currUid, &args)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

// `parseMemberIncludes` validates include query values.
func parseMemberIncludes(includes []enum.MemberIncl) bool {
	if len(includes) == 0 {
		return true
	}

	for i := range includes {
		switch includes[i] {
		case enum.MemberInclUser, enum.MemberInclTeam:
			continue

		default:
			return false
		}
	}

	return true
}

// `UpdateMemberRole` godoc
// @Summary Update Member Role
// @Description Update one member role mask by put semantics
// @Description The caller must be team admin
// @Description Auth: `authorization` cookie is preferred over `Authorization` header when both are present
// @Tags member
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param member_id path string true "member id"
// @Param body body val.MemberRoleUpdArgs true "update member role args"
// @Success 200 {object} res.HttpRes[any]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /api/v1/members/{member_id} [put]
func UpdateMemberRole(st *state.AppState) iris.Handler {
	memberApp := st.MemberApp

	return func(cx iris.Context) {
		currUid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		memberId := cx.Params().Get("member_id")
		if memberId == "" {
			res.Reject(cx, iris.StatusBadRequest, "缺少 member_id 参数")
			return
		}

		var args val.MemberRoleUpdArgs
		if err := cx.ReadJSON(&args); err != nil {
			res.Reject(cx, iris.StatusBadRequest, "请求参数解析失败")
			return
		}

		args.Id = memberId

		re := memberApp.UpdateRole(newReqCx(cx), currUid, &args)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

// `DeleteMember` godoc
// @Summary Delete Member
// @Description Hard delete one member by id
// @Description The caller must be team admin
// @Description Auth: `authorization` cookie is preferred over `Authorization` header when both are present
// @Tags member
// @Security ApiKeyAuth
// @Produce json
// @Param member_id path string true "member id"
// @Success 200 {object} res.HttpRes[any]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /api/v1/members/{member_id} [delete]
func DeleteMember(st *state.AppState) iris.Handler {
	memberApp := st.MemberApp

	return func(cx iris.Context) {
		currUid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		memberId := cx.Params().Get("member_id")
		if memberId == "" {
			res.Reject(cx, iris.StatusBadRequest, "缺少 member_id 参数")
			return
		}

		re := memberApp.Delete(newReqCx(cx), currUid, memberId)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

// `GetMemberByUserTeam` godoc
// @Summary Get Member By User And Team
// @Description Get one member record by `user_id` and `team_id`
// @Description The caller must be a member of the target team
// @Description Auth: `authorization` cookie is preferred over `Authorization` header when both are present
// @Tags member
// @Security ApiKeyAuth
// @Produce json
// @Param user_id query string true "user id"
// @Param team_id query string true "team id"
// @Param includes query []string false "include related fields, optional: user"
// @Success 200 {object} res.HttpRes[val.MemberVal]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 404 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /api/v1/members/detail [get]
func GetMemberByUserTeam(st *state.AppState) iris.Handler {
	memberApp := st.MemberApp

	return func(cx iris.Context) {
		currUid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		var args val.GetMemberByUserTeamIdArgs
		if err := cx.ReadQuery(&args); err != nil {
			res.Reject(cx, iris.StatusBadRequest, "请求参数解析失败")
			return
		}

		// Validate required query params.
		if args.UserId == "" {
			res.Reject(cx, iris.StatusBadRequest, "缺少 user_id 参数")
			return
		}

		if args.TeamId == "" {
			res.Reject(cx, iris.StatusBadRequest, "缺少 team_id 参数")
			return
		}

		// Validate includes.
		if !parseMemberIncludes(args.Includes) {
			res.Reject(cx, iris.StatusBadRequest, "includes 参数格式错误")
			return
		}

		re := memberApp.GetByUserTeamId(newReqCx(cx), currUid, &args)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

// `JoinTeamByInvitation` godoc
// @Summary Join Team By Invitation
// @Description Join one team by member invitation code
// @Description Auth: `authorization` cookie is preferred over `Authorization` header when both are present
// @Tags member
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param body body val.JoinTeamArgs true "join team args"
// @Success 200 {object} res.HttpRes[any]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /api/v1/members/join [post]
func JoinTeamByInvitation(st *state.AppState) iris.Handler {
	memberApp := st.MemberApp

	return func(cx iris.Context) {
		currUid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		var args val.JoinTeamArgs

		if err := cx.ReadJSON(&args); err != nil {
			res.Reject(cx, iris.StatusBadRequest, "请求参数解析失败")
			return
		}

		re := memberApp.JoinTeam(newReqCx(cx), currUid, &args)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}
