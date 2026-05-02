package http

import (
	"poprako-s/internal/api/http/res"
	"poprako-s/internal/api/state"
	"poprako-s/internal/app/val"

	"github.com/kataras/iris/v12"
)

// `ListWorksets` godoc
// @Summary List Worksets
//
//	List all active worksets for a team.
//	The caller must be a member of the specified team.
//	Auth: `authorization` cookie is preferred over `Authorization` header when both are present.
//
// @Tags workset
// @Security ApiKeyAuth
// @Produce json
// @Param team_id path string true "team id"
// @Param offset query int false "pagination offset"
// @Param limit query int false "pagination limit"
// @Success 200 {object} res.HttpRes "res.HttpRes{data=[]val.WorksetVal}"
// @Failure 400 {object} res.HttpRes
// @Failure 401 {object} res.HttpRes
// @Failure 500 {object} res.HttpRes
// @Router /workset/team/{team_id} [get]
func ListWorksets(st *state.AppState) iris.Handler {
	worksetApp := st.WorksetApp

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

		args := &val.ListWorksetArgs{
			TeamId: teamId,
			Offset: offset,
			Limit:  limit,
		}

		re := worksetApp.List(newReqCx(cx), currUid, args)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), "获取作品集列表失败")
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

// `CreateWorkset` godoc
// @Summary Create Workset
//
//	Create a new workset inside a team.
//	The caller must be an admin of the specified team.
//	Auth: `authorization` cookie is preferred over `Authorization` header when both are present.
//
// @Tags workset
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param body body val.CreateWorksetArgs true "create workset args"
// @Success 201 {object} res.HttpRes "res.HttpRes{data=val.WorksetCreatedRes}"
// @Failure 400 {object} res.HttpRes
// @Failure 401 {object} res.HttpRes
// @Failure 500 {object} res.HttpRes
// @Router /workset [post]
func CreateWorkset(st *state.AppState) iris.Handler {
	worksetApp := st.WorksetApp

	return func(cx iris.Context) {
		currUid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		var args val.CreateWorksetArgs

		if err := cx.ReadJSON(&args); err != nil {
			res.Reject(cx, iris.StatusBadRequest, "请求参数解析失败")
			return
		}

		re := worksetApp.Create(newReqCx(cx), currUid, &args)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusCreated, re.Data())
	}
}

// `UpdateWorkset` godoc
// @Summary Update Workset
//
//	Update the name and/or description of an existing workset.
//	The caller must be an admin of the workset's owning team.
//	Auth: `authorization` cookie is preferred over `Authorization` header when both are present.
//
// @Tags workset
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param workset_id path string true "workset id"
// @Param body body val.WorksetUpdArgs true "update workset args"
// @Success 200 {object} res.HttpRes
// @Failure 400 {object} res.HttpRes
// @Failure 401 {object} res.HttpRes
// @Failure 500 {object} res.HttpRes
// @Router /workset/{workset_id} [put]
func UpdateWorkset(st *state.AppState) iris.Handler {
	worksetApp := st.WorksetApp

	return func(cx iris.Context) {
		currUid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		worksetId := cx.Params().Get("workset_id")
		if worksetId == "" {
			res.Reject(cx, iris.StatusBadRequest, "缺少 workset_id 参数")
			return
		}

		var args val.WorksetUpdArgs

		if err := cx.ReadJSON(&args); err != nil {
			res.Reject(cx, iris.StatusBadRequest, "请求参数解析失败")
			return
		}

		// Bind path id into args so the inner app only needs one field.
		args.Id = worksetId

		re := worksetApp.Update(newReqCx(cx), currUid, &args)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

// `DeleteWorkset` godoc
// @Summary Delete Workset
//
//	Hard-delete a workset by id.
//	The caller must be an admin of the workset's owning team.
//	Auth: `authorization` cookie is preferred over `Authorization` header when both are present.
//
// @Tags workset
// @Security ApiKeyAuth
// @Produce json
// @Param workset_id path string true "workset id"
// @Success 200 {object} res.HttpRes
// @Failure 400 {object} res.HttpRes
// @Failure 401 {object} res.HttpRes
// @Failure 500 {object} res.HttpRes
// @Router /workset/{workset_id} [delete]
func DeleteWorkset(st *state.AppState) iris.Handler {
	worksetApp := st.WorksetApp

	return func(cx iris.Context) {
		currUid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		worksetId := cx.Params().Get("workset_id")
		if worksetId == "" {
			res.Reject(cx, iris.StatusBadRequest, "缺少 workset_id 参数")
			return
		}

		re := worksetApp.Delete(newReqCx(cx), currUid, worksetId)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}
