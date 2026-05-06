package http

import (
	"poprako-s/internal/api/http/res"
	"poprako-s/internal/api/state"
	"poprako-s/internal/app/val"

	"github.com/kataras/iris/v12"
)

// `CreateTeam` godoc
// @Summary Create Team
//
//	Create one team with super-admin permission.
//	Auth: `authorization` cookie is preferred over `Authorization` header when both are present
//
// @Tags team
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param body body val.TeamCreArgs true "create team args"
// @Success 201 {object} res.HttpRes[val.TeamCreRes]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /api/v1/teams [post]
func CreateTeam(st *state.AppState) iris.Handler {
	teamApp := st.TeamApp

	return func(cx iris.Context) {
		currUid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		var args val.TeamCreArgs
		if err := cx.ReadJSON(&args); err != nil {
			res.Reject(cx, iris.StatusBadRequest, "请求参数解析失败")
			return
		}

		re := teamApp.Create(newReqCx(cx), currUid, &args)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusCreated, re.Data())
	}
}

// `GetTeamInfo` godoc
// @Summary Get Team Info
//
//	Get team info by team id and return a `res.HttpRes` wrapper with `val.TeamVal`.
//	Auth: `authorization` cookie is preferred over `Authorization` header when both are present
//
// @Tags team
// @Security ApiKeyAuth
// @Produce json
// @Param team_id path string true "team id"
// @Success 200 {object} res.HttpRes[val.TeamVal]
// @Failure 400 {object} res.HttpRes[any]
// @Router /api/v1/teams/{team_id} [get]
func GetTeamInfo(st *state.AppState) iris.Handler {
	teamApp := st.TeamApp

	return func(cx iris.Context) {
		teamId := cx.Params().Get("team_id")
		if teamId == "" {
			res.Reject(cx, iris.StatusBadRequest, "缺少 team_id 参数")
			return
		}

		re := teamApp.GetInfo(newReqCx(cx), teamId)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

// `ListTeams` godoc
// @Summary List Teams
//
//	List all teams
//	Auth: `authorization` cookie is preferred over `Authorization` header when both are present
//
// @Tags team
// @Security ApiKeyAuth
// @Produce json
// @Param offset query int false "pagination offset"
// @Param limit query int false "pagination limit"
// @Success 200 {object} res.HttpRes[[]val.TeamVal]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /api/v1/teams [get]
func ListTeams(st *state.AppState) iris.Handler {
	teamApp := st.TeamApp

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

		re := teamApp.List(newReqCx(cx), currUid, &val.ListTeamArgs{Offset: offset, Limit: limit})
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

// `ListMyTeams` godoc
// @Summary List My Teams
//
//	List teams that current user belongs to
//	Auth: `authorization` cookie is preferred over `Authorization` header when both are present
//
// @Tags team
// @Security ApiKeyAuth
// @Produce json
// @Param offset query int false "pagination offset"
// @Param limit query int false "pagination limit"
// @Success 200 {object} res.HttpRes[[]val.TeamVal]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /api/v1/teams/mine [get]
func ListMyTeams(st *state.AppState) iris.Handler {
	teamApp := st.TeamApp

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

		re := teamApp.ListByUser(newReqCx(cx), currUid, &val.ListTeamArgs{Offset: offset, Limit: limit})
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

// `UpdateTeam` godoc
// @Summary Update Team
//
//	Update one team by put semantics
//	Auth: `authorization` cookie is preferred over `Authorization` header when both are present
//
// @Tags team
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param team_id path string true "team id"
// @Param body body val.TeamUpdArgs true "update team args"
// @Success 200 {object} res.HttpRes[any]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /api/v1/teams/{team_id} [put]
func UpdateTeam(st *state.AppState) iris.Handler {
	teamApp := st.TeamApp

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

		var args val.TeamUpdArgs
		if err := cx.ReadJSON(&args); err != nil {
			res.Reject(cx, iris.StatusBadRequest, "请求参数解析失败")
			return
		}

		args.Id = teamId

		re := teamApp.Update(newReqCx(cx), currUid, &args)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

// `ReserveTeamAvatar` godoc
// @Summary Reserve Team Avatar Upload
//
//	Reserve team avatar upload and return one signed put url
//	Auth: `authorization` cookie is preferred over `Authorization` header when both are present
//
// @Tags team
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param team_id path string true "team id"
// @Param body body val.ResvTeamAvatarArgs true "reserve team avatar args"
// @Success 200 {object} res.HttpRes[val.ResvTeamAvatarRes]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /api/v1/teams/{team_id}/avatar [post]
func ReserveTeamAvatar(st *state.AppState) iris.Handler {
	teamApp := st.TeamApp

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

		var args val.ResvTeamAvatarArgs
		if err := cx.ReadJSON(&args); err != nil {
			res.Reject(cx, iris.StatusBadRequest, "请求参数解析失败")
			return
		}

		args.TeamId = teamId

		re := teamApp.ResvAvatar(newReqCx(cx), currUid, &args)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

// `ConfirmTeamAvatarUploaded` godoc
// @Summary Confirm Team Avatar Uploaded
//
//	Confirm team avatar uploaded status
//	Auth: `authorization` cookie is preferred over `Authorization` header when both are present
//
// @Tags team
// @Security ApiKeyAuth
// @Produce json
// @Param team_id path string true "team id"
// @Success 200 {object} res.HttpRes[any]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /api/v1/teams/{team_id}/avatar/confirm [post]
func ConfirmTeamAvatarUploaded(st *state.AppState) iris.Handler {
	teamApp := st.TeamApp

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

		re := teamApp.MarkAvatarUploaded(newReqCx(cx), currUid, teamId)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}
