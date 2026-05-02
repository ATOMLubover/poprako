package http

import (
	"poprako-s/internal/api/http/res"
	"poprako-s/internal/api/state"
	"poprako-s/internal/app/val"

	"github.com/kataras/iris/v12"
)

// `ListAssignmentInvitations` godoc
// @Summary List Assignment Invitations
// @Description List assignment invitations for one chapter
// @Description The caller must be reviewer of the target chapter
// @Description Auth: `authorization` cookie is preferred over `Authorization` header when both are present
// @Tags assignment-invitation
// @Security ApiKeyAuth
// @Produce json
// @Param chapter_id path string true "chapter id"
// @Param pending query bool false "pending filter"
// @Param offset query int false "pagination offset"
// @Param limit query int false "pagination limit"
// @Success 200 {object} res.HttpRes[[]val.AssignmentInvVal]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /assignment-invitations/chapter/{chapter_id} [get]
func ListAssignmentInvitations(st *state.AppState) iris.Handler {
	app := st.AssignmentInvApp

	return func(cx iris.Context) {
		currUid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		chapterId := cx.Params().Get("chapter_id")
		if chapterId == "" {
			res.Reject(cx, iris.StatusBadRequest, "缺少 chapter_id 参数")
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

		var pending *bool
		if p := cx.URLParam("pending"); p != "" {
			pendingVal, boolErr := cx.URLParamBool("pending")
			if boolErr != nil {
				res.Reject(cx, iris.StatusBadRequest, "pending 参数格式错误")
				return
			}
			pending = &pendingVal
		}

		re := app.ListByChapter(newReqCx(cx), currUid, &val.ListAssignmentInvArgs{ChapterId: chapterId, Pending: pending, Offset: offset, Limit: limit})
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

// `CreateAssignmentInvitation` godoc
// @Summary Create Assignment Invitation
// @Description Create assignment invitation for one chapter
// @Description The caller must be reviewer of the target chapter
// @Description Auth: `authorization` cookie is preferred over `Authorization` header when both are present
// @Tags assignment-invitation
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param body body val.CreateAssignmentInvArgs true "create assignment invitation args"
// @Success 201 {object} res.HttpRes[val.CreateAssignmentInvRes]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /assignment-invitations [post]
func CreateAssignmentInvitation(st *state.AppState) iris.Handler {
	app := st.AssignmentInvApp

	return func(cx iris.Context) {
		currUid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		var args val.CreateAssignmentInvArgs
		if err := cx.ReadJSON(&args); err != nil {
			res.Reject(cx, iris.StatusBadRequest, "请求参数解析失败")
			return
		}

		re := app.Create(newReqCx(cx), currUid, &args)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusCreated, re.Data())
	}
}

// `DeleteAssignmentInvitation` godoc
// @Summary Delete Assignment Invitation
// @Description Delete assignment invitation by id
// @Description The caller must be reviewer of the target chapter
// @Description Auth: `authorization` cookie is preferred over `Authorization` header when both are present
// @Tags assignment-invitation
// @Security ApiKeyAuth
// @Produce json
// @Param invitation_id path string true "invitation id"
// @Success 200 {object} res.HttpRes[any]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /assignment-invitations/{invitation_id} [delete]
func DeleteAssignmentInvitation(st *state.AppState) iris.Handler {
	app := st.AssignmentInvApp

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

		re := app.Delete(newReqCx(cx), currUid, invId)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

// `JoinByAssignmentInvitation` godoc
// @Summary Join Chapter By Invitation
// @Description Join chapter collaboration by assignment invitation code
// @Description Auth: `authorization` cookie is preferred over `Authorization` header when both are present
// @Tags assignment-invitation
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param body body val.JoinAssignmentInvArgs true "join args"
// @Success 200 {object} res.HttpRes[any]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /assignment-invitations/join [post]
func JoinByAssignmentInvitation(st *state.AppState) iris.Handler {
	app := st.AssignmentInvApp

	return func(cx iris.Context) {
		currUid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		var args val.JoinAssignmentInvArgs
		if err := cx.ReadJSON(&args); err != nil {
			res.Reject(cx, iris.StatusBadRequest, "请求参数解析失败")
			return
		}

		re := app.JoinByInvCode(newReqCx(cx), currUid, &args)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}
