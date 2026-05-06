package http

import (
	"poprako-s/internal/api/http/res"
	"poprako-s/internal/api/state"
	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/model/enum"

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
// @Param chapter_id query string true "chapter id"
// @Param includes query []string false "include related fields, optional: user, chapter, chapter.comic, chapter.comic.workset, chapter.comic.workset.team, chapter.creator, chapter.comic.creator"
// @Param offset query int false "pagination offset"
// @Param limit query int false "pagination limit"
// @Success 200 {object} res.HttpRes[[]val.AssignmentVal]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /api/v1/assignments [get]
func ListAssignmentsByChapter(st *state.AppState) iris.Handler {
	app := st.AssignmentApp

	return func(cx iris.Context) {
		currUid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		var args val.ListAssignmentByChapterArgs
		if err := cx.ReadQuery(&args); err != nil {
			res.Reject(cx, iris.StatusBadRequest, "请求参数解析失败")
			return
		}

		if args.ChapterId == "" {
			res.Reject(cx, iris.StatusBadRequest, "缺少 chapter_id 参数")
			return
		}

		if !parseAssignmentIncludes(args.Includes) {
			res.Reject(cx, iris.StatusBadRequest, "includes 参数格式错误")
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

// `ListMyAssignments` godoc
// @Summary List My Assignments
// @Description List assignments of current user
// @Description Auth: `authorization` cookie is preferred over `Authorization` header when both are present
// @Tags assignment
// @Security ApiKeyAuth
// @Produce json
// @Param includes query []string false "include related fields, optional: user, chapter, chapter.comic, chapter.comic.workset, chapter.comic.workset.team, chapter.creator, chapter.comic.creator"
// @Param offset query int false "pagination offset"
// @Param limit query int false "pagination limit"
// @Success 200 {object} res.HttpRes[[]val.AssignmentVal]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /api/v1/assignments/mine [get]
func ListMyAssignments(st *state.AppState) iris.Handler {
	app := st.AssignmentApp

	return func(cx iris.Context) {
		currUid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		var args val.ListAssignmentByUserArgs
		if err := cx.ReadQuery(&args); err != nil {
			res.Reject(cx, iris.StatusBadRequest, "请求参数解析失败")
			return
		}

		if !parseAssignmentIncludes(args.Includes) {
			res.Reject(cx, iris.StatusBadRequest, "includes 参数格式错误")
			return
		}

		re := app.ListByUser(newReqCx(cx), currUid, &args)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

// `parseAssignmentIncludes` validates include query values.
func parseAssignmentIncludes(includes []enum.AssignmentIncl) bool {
	if len(includes) == 0 {
		return true
	}

	for i := range includes {
		switch includes[i] {
		case enum.AssignmentInclUser,
			enum.AssignmentInclChapter,
			enum.AssignmentInclChapterComic,
			enum.AssignmentInclChapterComicWorkset,
			enum.AssignmentInclChapterComicWorksetTeam,
			enum.AssignmentInclChapterCreator,
			enum.AssignmentInclChapterComicCreator:
			continue

		default:
			return false
		}
	}

	return true
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
// @Router /api/v1/assignments [put]
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
// @Router /api/v1/assignments/{assignment_id} [delete]
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
