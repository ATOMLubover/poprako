package http

import (
	"poprako-s/internal/api/http/res"
	"poprako-s/internal/api/state"
	"poprako-s/internal/app/val"

	"github.com/kataras/iris/v12"
)

// `ListAssignmentsByChapter` godoc
// @Summary List Assignments By Chapter
// @Description List assignments for one chapter
// @Description The caller must be reviewer of the target chapter
// @Description Auth: `authorization` cookie is preferred over `Authorization` header when both are present
// @Tags assignment
// @Security ApiKeyAuth
// @Produce json
// @Param chapter_id path string true "chapter id"
// @Param offset query int false "pagination offset"
// @Param limit query int false "pagination limit"
// @Success 200 {object} res.HttpRes[[]val.AssignmentVal]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /assignments/chapters/{chapter_id} [get]
func ListAssignmentsByChapter(st *state.AppState) iris.Handler {
	app := st.AssignmentApp

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

		re := app.ListByChapter(newReqCx(cx), currUid, &val.ListAssignmentByChapterArgs{ChapterId: chapterId, Offset: offset, Limit: limit})
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

// `ListMyAssignments` godoc
// @Summary List My Assignments
// @Description List assignments of current user
// @Description Auth: `authorization` cookie is preferred over `Authorization` header when both are present
// @Tags assignment
// @Security ApiKeyAuth
// @Produce json
// @Param offset query int false "pagination offset"
// @Param limit query int false "pagination limit"
// @Success 200 {object} res.HttpRes[[]val.AssignmentVal]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /assignments/mine [get]
func ListMyAssignments(st *state.AppState) iris.Handler {
	app := st.AssignmentApp

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

		re := app.ListByUser(newReqCx(cx), currUid, &val.ListAssignmentByUserArgs{Offset: offset, Limit: limit})
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

// `UpsertAssignment` godoc
// @Summary Upsert Assignment
// @Description Upsert assignment by PUT semantics
// @Description If `role_mask` is zero the request will redirect to delete semantics
// @Description Auth: `authorization` cookie is preferred over `Authorization` header when both are present
// @Tags assignment
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param body body val.UpsertAssignmentArgs true "upsert args"
// @Success 200 {object} res.HttpRes[any]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /assignments [put]
func UpsertAssignment(st *state.AppState) iris.Handler {
	app := st.AssignmentApp

	return func(cx iris.Context) {
		currUid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		var args val.UpsertAssignmentArgs
		if err := cx.ReadJSON(&args); err != nil {
			res.Reject(cx, iris.StatusBadRequest, "请求参数解析失败")
			return
		}

		re := app.Upsert(newReqCx(cx), currUid, &args)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

// `DeleteAssignment` godoc
// @Summary Delete Assignment
// @Description Delete assignment by id
// @Description The caller must be reviewer of target chapter
// @Description Auth: `authorization` cookie is preferred over `Authorization` header when both are present
// @Tags assignment
// @Security ApiKeyAuth
// @Produce json
// @Param assignment_id path string true "assignment id"
// @Success 200 {object} res.HttpRes[any]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /assignments/{assignment_id} [delete]
func DeleteAssignment(st *state.AppState) iris.Handler {
	app := st.AssignmentApp

	return func(cx iris.Context) {
		currUid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		assignmentId := cx.Params().Get("assignment_id")
		if assignmentId == "" {
			res.Reject(cx, iris.StatusBadRequest, "缺少 assignment_id 参数")
			return
		}

		re := app.Delete(newReqCx(cx), currUid, assignmentId)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}
