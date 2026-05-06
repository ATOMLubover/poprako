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
// @Param chapter_id query string true "chapter id"
// @Param pending query bool false "pending filter"
// @Param offset query int false "pagination offset"
// @Param limit query int false "pagination limit"
// @Success 200 {object} res.HttpRes[[]val.AssignmentInvVal]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /api/v1/assignment-invitations [get]
func ListAssignmentInvitations(st *state.AppState) iris.Handler {
	app := st.AssignmentInvApp

	return func(cx iris.Context) {
		currUid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		var args val.ListAssignmentInvArgs
		if err := cx.ReadQuery(&args); err != nil {
			res.Reject(cx, iris.StatusBadRequest, "请求参数解析失败")
			return
		}

		if args.ChapterId == "" {
			res.Reject(cx, iris.StatusBadRequest, "缺少 chapter_id 参数")
			return
		}

		re := app.ListByChapter(newReqCx(cx), currUid, &args)
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
// @Router /api/v1/assignment-invitations [post]
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
// @Router /api/v1/assignment-invitations/{invitation_id} [delete]
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
// @Router /api/v1/assignment-invitations/join [post]
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
