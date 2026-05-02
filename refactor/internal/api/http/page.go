package http

import (
	"poprako-s/internal/api/http/res"
	"poprako-s/internal/api/state"
	"poprako-s/internal/app/val"

	"github.com/kataras/iris/v12"
)

// `ResvChapterPages` godoc
// @Summary Reserve Chapter Pages Upload
// @Description Reserve signed upload URLs for all pages of one chapter
// @Description The caller must be an admin of the owning team
// @Description Auth: `authorization` cookie is preferred over `Authorization` header when both are present
// @Tags chapter
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param chapter_id path string true "chapter id"
// @Param body body val.ResvChapterPagesArgs true "reserve chapter pages args"
// @Success 200 {object} res.HttpRes[val.ResvChapterPagesRes]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /chapter/{chapter_id}/pages/reserve [post]
func ResvChapterPages(st *state.AppState) iris.Handler {
	pageApp := st.PageApp

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

		var args val.ResvChapterPagesArgs
		if err := cx.ReadJSON(&args); err != nil {
			res.Reject(cx, iris.StatusBadRequest, "请求参数解析失败")
			return
		}
		args.ChapterId = chapterId

		re := pageApp.ResvChapterPages(newReqCx(cx), currUid, &args)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

// `ListChapterPages` godoc
// @Summary List Chapter Pages
// @Description List pages under one chapter with pagination
// @Description The caller must be a member of the owning team
// @Description Auth: `authorization` cookie is preferred over `Authorization` header when both are present
// @Tags chapter
// @Security ApiKeyAuth
// @Produce json
// @Param chapter_id path string true "chapter id"
// @Param offset query int false "pagination offset"
// @Param limit query int false "pagination limit"
// @Success 200 {object} res.HttpRes[[]val.PageVal]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /chapter/{chapter_id}/pages [get]
func ListChapterPages(st *state.AppState) iris.Handler {
	pageApp := st.PageApp

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

		args := &val.ListChapterPageArgs{ChapterId: chapterId, Offset: offset, Limit: limit}
		re := pageApp.List(newReqCx(cx), currUid, args)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

// `MarkPageImageUploaded` godoc
// @Summary Confirm Page Image Uploaded
// @Description Confirm one page image upload after client upload completed
// @Description The caller must be assigned to the target chapter
// @Description Auth: `authorization` cookie is preferred over `Authorization` header when both are present
// @Tags page
// @Security ApiKeyAuth
// @Produce json
// @Param page_id path string true "page id"
// @Success 200 {object} res.HttpRes[any]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /page/{page_id}/image/uploaded [post]
func MarkPageImageUploaded(st *state.AppState) iris.Handler {
	pageApp := st.PageApp

	return func(cx iris.Context) {
		currUid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		pageId := cx.Params().Get("page_id")
		if pageId == "" {
			res.Reject(cx, iris.StatusBadRequest, "缺少 page_id 参数")
			return
		}

		re := pageApp.MarkImageUploaded(newReqCx(cx), currUid, &val.MarkPageImageUploadedArgs{PageId: pageId})
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

// `DeleteChapterPages` godoc
// @Summary Delete Chapter Pages
// @Description Hard-delete all pages under one chapter
// @Description The caller must be an admin of the owning team
// @Description Auth: `authorization` cookie is preferred over `Authorization` header when both are present
// @Tags chapter
// @Security ApiKeyAuth
// @Produce json
// @Param chapter_id path string true "chapter id"
// @Success 200 {object} res.HttpRes[any]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /chapter/{chapter_id}/pages [delete]
func DeleteChapterPages(st *state.AppState) iris.Handler {
	pageApp := st.PageApp

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

		re := pageApp.DeleteByChapterId(newReqCx(cx), currUid, chapterId)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}
