package http

import (
	"poprako-s/internal/api/http/res"
	"poprako-s/internal/api/state"
	"poprako-s/internal/app/val"

	"github.com/kataras/iris/v12"
)

// `ListChapters` godoc
// @Summary List Chapters
//
//	List chapters for one comic.
//	The caller must be a member of the target comic team.
//	Auth: `authorization` cookie is preferred over `Authorization` header when both are present.
//
// @Tags chapter
// @Security ApiKeyAuth
// @Produce json
// @Param comic_id path string true "comic id"
// @Param offset query int false "pagination offset"
// @Param limit query int false "pagination limit"
// @Success 200 {object} res.HttpRes "res.HttpRes{data=[]val.ChapterVal}"
// @Failure 400 {object} res.HttpRes
// @Failure 401 {object} res.HttpRes
// @Failure 500 {object} res.HttpRes
// @Router /chapter/comic/{comic_id} [get]
func ListChapters(st *state.AppState) iris.Handler {
	chapterApp := st.ChapterApp

	return func(cx iris.Context) {
		currUid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		comicId := cx.Params().Get("comic_id")
		if comicId == "" {
			res.Reject(cx, iris.StatusBadRequest, "缺少 comic_id 参数")
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

		args := &val.ListChapterArgs{ComicId: comicId, Offset: offset, Limit: limit}

		re := chapterApp.List(newReqCx(cx), currUid, args)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

// `GetChapterById` godoc
// @Summary Get Chapter By Id
//
//	Get one chapter by id.
//	The caller must be a member of the owning team.
//	Auth: `authorization` cookie is preferred over `Authorization` header when both are present.
//
// @Tags chapter
// @Security ApiKeyAuth
// @Produce json
// @Param chapter_id path string true "chapter id"
// @Success 200 {object} res.HttpRes "res.HttpRes{data=val.ChapterVal}"
// @Failure 400 {object} res.HttpRes
// @Failure 404 {object} res.HttpRes
// @Failure 401 {object} res.HttpRes
// @Failure 500 {object} res.HttpRes
// @Router /chapter/{chapter_id} [get]
func GetChapterById(st *state.AppState) iris.Handler {
	chapterApp := st.ChapterApp

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

		re := chapterApp.GetById(newReqCx(cx), currUid, chapterId)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

// `GetPinnedChapter` godoc
// @Summary Get Pinned Chapter
//
//	Get pinned chapter for one comic.
//	The caller must be a member of the target comic team.
//	Auth: `authorization` cookie is preferred over `Authorization` header when both are present.
//
// @Tags chapter
// @Security ApiKeyAuth
// @Produce json
// @Param comic_id path string true "comic id"
// @Success 200 {object} res.HttpRes "res.HttpRes{data=val.ChapterVal}"
// @Failure 400 {object} res.HttpRes
// @Failure 401 {object} res.HttpRes
// @Failure 500 {object} res.HttpRes
// @Router /chapter/comic/{comic_id}/pinned [get]
func GetPinnedChapter(st *state.AppState) iris.Handler {
	chapterApp := st.ChapterApp

	return func(cx iris.Context) {
		currUid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		comicId := cx.Params().Get("comic_id")
		if comicId == "" {
			res.Reject(cx, iris.StatusBadRequest, "缺少 comic_id 参数")
			return
		}

		re := chapterApp.GetPinned(newReqCx(cx), currUid, comicId)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

// `CreateChapter` godoc
// @Summary Create Chapter
//
//	Create chapter under one comic.
//	The caller must be an admin of the target comic team.
//	Auth: `authorization` cookie is preferred over `Authorization` header when both are present.
//
// @Tags chapter
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param body body val.CreateChapterArgs true "create chapter args"
// @Success 201 {object} res.HttpRes "res.HttpRes{data=val.ChapterCreatedRes}"
// @Failure 400 {object} res.HttpRes
// @Failure 401 {object} res.HttpRes
// @Failure 500 {object} res.HttpRes
// @Router /chapter [post]
func CreateChapter(st *state.AppState) iris.Handler {
	chapterApp := st.ChapterApp

	return func(cx iris.Context) {
		currUid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		var args val.CreateChapterArgs
		if err := cx.ReadJSON(&args); err != nil {
			res.Reject(cx, iris.StatusBadRequest, "请求参数解析失败")
			return
		}

		re := chapterApp.Create(newReqCx(cx), currUid, &args)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusCreated, re.Data())
	}
}

// `UpdateChapter` godoc
// @Summary Update Chapter
//
//	Update one chapter.
//	The caller must be an admin of the target chapter team.
//	Auth: `authorization` cookie is preferred over `Authorization` header when both are present.
//
// @Tags chapter
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param chapter_id path string true "chapter id"
// @Param body body val.ChapterUpdArgs true "update chapter args"
// @Success 200 {object} res.HttpRes
// @Failure 400 {object} res.HttpRes
// @Failure 401 {object} res.HttpRes
// @Failure 500 {object} res.HttpRes
// @Router /chapter/{chapter_id} [put]
func UpdateChapter(st *state.AppState) iris.Handler {
	chapterApp := st.ChapterApp

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

		var args val.ChapterUpdArgs
		if err := cx.ReadJSON(&args); err != nil {
			res.Reject(cx, iris.StatusBadRequest, "请求参数解析失败")
			return
		}

		args.Id = chapterId

		re := chapterApp.Update(newReqCx(cx), currUid, &args)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

// `RemoveChapter` godoc
// @Summary Remove Chapter
//
//	Soft-delete one chapter.
//	The caller must be an admin of the target chapter team.
//	Auth: `authorization` cookie is preferred over `Authorization` header when both are present.
//
// @Tags chapter
// @Security ApiKeyAuth
// @Produce json
// @Param chapter_id path string true "chapter id"
// @Success 200 {object} res.HttpRes
// @Failure 400 {object} res.HttpRes
// @Failure 401 {object} res.HttpRes
// @Failure 500 {object} res.HttpRes
// @Router /chapter/{chapter_id} [delete]
func RemoveChapter(st *state.AppState) iris.Handler {
	chapterApp := st.ChapterApp

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

		re := chapterApp.Remove(newReqCx(cx), currUid, chapterId)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}
