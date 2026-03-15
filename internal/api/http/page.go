package http

import (
	"labelplus-next-web-be/internal/state"
	"labelplus-next-web-be/internal/value"

	"github.com/kataras/iris/v12"
)

// ListChapterPages godoc
// @Summary 	获取章节页面列表
// @Description 获取指定章节的所有页面，注意当列表为空，会返回 null 而不是空数组
//
// @Tags 		page
// @Security 	ApiKeyAuth
// @Produce 	json
// @Param 		chapter_id query string true "章节 ID"
// @Param 		offset query int true "偏移量"
// @Param 		limit query int true "每页数量"
// @Param 		"includes[]" query []string false "include 关联信息，可选值：creator"
//
// @Success 	200 {object} []value.PageInfo
//
// @Router 		/pages [get]
func ListChapterPages(appState *state.AppState) iris.Handler {
	pageApplication := appState.PageApplication

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrentUserID(ctx)
		if !ok {
			return
		}

		var args value.ListChapterPageArgs
		if err := ctx.ReadQuery(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "查询参数格式错误: "+err.Error())
			return
		}

		result, err := pageApplication.ListChapterPages(
			buildTraceScope(ctx),
			currentUserID,
			args,
		)
		if err != nil {
			reject(ctx, iris.StatusForbidden, err.Error())
			return
		}

		accept(ctx, "获取页面列表成功", result)
	}
}

// ReserveChapterPages godoc
// @Summary 	预留页面记录，并生成每个页面的预签名上传地址
// @Description 为指定章节批量创建页面，返回每个页面的预签名上传地址
//
// @Tags 		page
// @Security 	ApiKeyAuth
// @Accept 	json
// @Produce 	json
// @Param 		body body value.ReserveChapterPagesArgs true "预留页面参数"
//
// @Success 	201 {object} value.ReserveChapterPagesResult
//
// @Router 		/pages [post]
func ReserveChapterPages(appState *state.AppState) iris.Handler {
	pageApplication := appState.PageApplication

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrentUserID(ctx)
		if !ok {
			return
		}

		var args value.ReserveChapterPagesArgs
		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误: "+err.Error())
			return
		}

		result, err := pageApplication.ReserveChapterPages(
			buildTraceScope(ctx),
			currentUserID,
			args,
		)
		if err != nil {
			reject(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		ctx.StatusCode(iris.StatusCreated)
		accept(ctx, "创建页面成功", result)
	}
}

// UpdatePage godoc
// @Summary 	更新页面
// @Description 更新指定页面的信息，例如标记页面已上传完成
//
// @Tags 		page
// @Security 	ApiKeyAuth
// @Accept 	json
// @Produce 	json
// @Param 		page_id path string true "页面 ID"
// @Param 		body body value.UpdatePageArgs true "更新页面参数"
//
// @Success 	200
//
// @Router 		/pages/{page_id} [put]
func UpdatePage(appState *state.AppState) iris.Handler {
	pageApplication := appState.PageApplication

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrentUserID(ctx)
		if !ok {
			return
		}

		pageID := ctx.Params().Get("page_id")
		if pageID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 page_id 路径参数")
			return
		}

		var args value.UpdatePageArgs
		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误: "+err.Error())
			return
		}

		if args.ID != pageID {
			reject(ctx, iris.StatusBadRequest, "路径参数 page_id 与请求体中的 id 不匹配")
			return
		}

		if err := pageApplication.UpdatePage(
			buildTraceScope(ctx),
			currentUserID,
			args,
		); err != nil {
			reject(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		accept(ctx, "更新页面成功", nil)
	}
}

// DeletePages godoc
// @Summary 	删除章节所有页面
// @Description 删除指定章节的所有页面
//
// @Tags 		page
// @Security 	ApiKeyAuth
// @Produce 	json
// @Param 		chapter_id path string true "章节 ID"
//
// @Success 	200
//
// @Router 		/pages/{chapter_id} [delete]
func DeletePages(appState *state.AppState) iris.Handler {
	pageApplication := appState.PageApplication

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

		if err := pageApplication.DeletePages(
			buildTraceScope(ctx),
			currentUserID,
			chapterID,
		); err != nil {
			reject(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		accept(ctx, "删除页面成功", nil)
	}
}
