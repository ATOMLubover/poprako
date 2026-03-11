package http

import (
	"labelplus-next-web-be/internal/state"
	"labelplus-next-web-be/internal/value"

	"github.com/kataras/iris/v12"
)

// ListComicChapters godoc
// @Summary 	获取漫画章节列表
// @Description 获取指定漫画的章节列表，支持分页，注意当列表为空，会返回 null 而不是空数组
//
// @Tags 		chapter
// @Security 	ApiKeyAuth
// @Produce 	json
// @Param 		comic_id query string true "漫画 ID"
// @Param 		offset query int true "偏移量"
// @Param 		limit query int true "每页数量"
//
// @Success 	200 {object} []value.ChapterDetail
//
// @Router 		/chapters [get]
func ListComicChapters(appState *state.AppState) iris.Handler {
	chapterApplication := appState.ChapterApplication

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrentUserID(ctx)
		if !ok {
			return
		}

		var args value.ListComicChapterArgs
		if err := ctx.ReadQuery(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "查询参数格式错误: "+err.Error())
			return
		}

		result, err := chapterApplication.ListComicChapters(
			buildTraceScope(ctx),
			currentUserID,
			args,
		)
		if err != nil {
			reject(ctx, iris.StatusForbidden, err.Error())
			return
		}

		accept(ctx, "获取章节列表成功", result)
	}
}

// CreateComicChapter godoc
// @Summary 	创建漫画章节
// @Description 在指定漫画中创建章节
//
// @Tags 		chapter
// @Security 	ApiKeyAuth
// @Accept 	json
// @Produce 	json
// @Param 		body body value.CreateChapterArgs true "创建章节参数"
//
// @Success 	201 {object} value.CreateChapterResult
//
// @Router 		/chapters [post]
func CreateComicChapter(appState *state.AppState) iris.Handler {
	chapterApplication := appState.ChapterApplication

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrentUserID(ctx)
		if !ok {
			return
		}

		var args value.CreateChapterArgs
		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误: "+err.Error())
			return
		}

		result, err := chapterApplication.CreateComicChapter(
			buildTraceScope(ctx),
			currentUserID,
			args,
		)
		if err != nil {
			reject(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		ctx.StatusCode(iris.StatusCreated)
		accept(ctx, "创建章节成功", result)
	}
}

// UpdateChapter godoc
// @Summary 	更新章节（PATCH 语义，仅传需要更新的字段）
// @Description 局部更新指定章节的信息，未传的字段不会被修改
//
// @Tags 		chapter
// @Security 	ApiKeyAuth
// @Accept 	json
// @Produce 	json
// @Param 		chapter_id path string true "章节 ID"
// @Param 		body body value.UpdateChapterArgs true "更新章节参数"
//
// @Success 	200
//
// @Router 		/chapters/{chapter_id} [patch]
func UpdateChapter(appState *state.AppState) iris.Handler {
	chapterApplication := appState.ChapterApplication

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrentUserID(ctx)
		if !ok {
			return
		}

		chapterID := ctx.Params().Get("chapter_id")
		if chapterID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 chapter_id 路径参数")
			return
		}

		var args value.UpdateChapterArgs
		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误: "+err.Error())
			return
		}
		args.ChapterID = chapterID

		if err := chapterApplication.UpdateChapter(
			buildTraceScope(ctx),
			currentUserID,
			args,
		); err != nil {
			reject(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		accept(ctx, "更新章节成功", nil)
	}
}

// DeleteComicChapter godoc
// @Summary 	删除章节
// @Description 删除指定章节
//
// @Tags 		chapter
// @Security 	ApiKeyAuth
// @Produce 	json
// @Param 		chapter_id path string true "章节 ID"
//
// @Success 	200
//
// @Router 		/chapters/{chapter_id} [delete]
func DeleteComicChapter(appState *state.AppState) iris.Handler {
	chapterApplication := appState.ChapterApplication

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrentUserID(ctx)
		if !ok {
			return
		}

		chapterID := ctx.Params().Get("chapter_id")
		if chapterID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 chapter_id 路径参数")
			return
		}

		if err := chapterApplication.DeleteComicChapter(
			buildTraceScope(ctx),
			currentUserID,
			chapterID,
		); err != nil {
			reject(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		accept(ctx, "删除章节成功", nil)
	}
}
