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
// @Router /members [post]
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
// @Param team_id path string true "team id"
// @Param includes query []string false "include related fields, optional: user, team"
// @Param offset query int false "pagination offset"
// @Param limit query int false "pagination limit"
// @Success 200 {object} res.HttpRes[[]val.MemberVal]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /members/team/{team_id} [get]
func ListTeamMembers(st *state.AppState) iris.Handler {
	memberApp := st.MemberApp

	return func(cx iris.Context) {
		currUid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		teamId := cx.Params().Get("team_id")
		if teamId == "" {
			res.Reject(cx, iris.StatusBadRequest, "缺少 team_id 参数")
			return
		}

		offset, err := cx.URLParamInt("offset")
		if err != nil {
			res.Reject(cx, iris.StatusBadRequest, "offset 参数格式错误")
			return
		}

		limit, err := cx.URLParamInt("limit")
		if err != nil {
			res.Reject(cx, iris.StatusBadRequest, "limit 参数格式错误")
			return
		}

		includes, ok := parseMemberIncludes(cx.URLParamSlice("includes"))
		if !ok {
			res.Reject(cx, iris.StatusBadRequest, "includes 参数格式错误")
			return
		}

		re := memberApp.ListByTeam(newReqCx(cx), currUid, &val.ListMemberByTeamArgs{
			TeamId:   teamId,
			Includes: includes,
			Offset:   offset,
			Limit:    limit,
		})
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
// @Router /members/mine [get]
func ListMyMembers(st *state.AppState) iris.Handler {
	memberApp := st.MemberApp

	return func(cx iris.Context) {
		currUid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		offset, err := cx.URLParamInt("offset")
		if err != nil {
			res.Reject(cx, iris.StatusBadRequest, "offset 参数格式错误")
			return
		}

		limit, err := cx.URLParamInt("limit")
		if err != nil {
			res.Reject(cx, iris.StatusBadRequest, "limit 参数格式错误")
			return
		}

		includes, ok := parseMemberIncludes(cx.URLParamSlice("includes"))
		if !ok {
			res.Reject(cx, iris.StatusBadRequest, "includes 参数格式错误")
			return
		}

		re := memberApp.ListMine(newReqCx(cx), currUid, &val.ListMyMemberArgs{
			Includes: includes,
			Offset:   offset,
			Limit:    limit,
		})
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
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
// @Router /members/{member_id} [put]
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
// @Router /members/{member_id} [delete]
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
// @Router /members/join [post]
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

// `parseMemberIncludes` parses include query values into typed includes.
func parseMemberIncludes(rawIncludes []string) ([]enum.MemberIncl, bool) {
	if len(rawIncludes) == 0 {
		return nil, true
	}

	includes := make([]enum.MemberIncl, 0, len(rawIncludes))
	for i := range rawIncludes {
		incl := enum.MemberIncl(rawIncludes[i])
		switch incl {
		case enum.MemberInclUser, enum.MemberInclTeam:
			includes = append(includes, incl)
		default:
			return nil, false
		}
	}

	return includes, true
}
