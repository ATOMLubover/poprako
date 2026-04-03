package http

import (
	"poprako-s/internal/app/val"
	"poprako-s/internal/state"

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
// @Success 	200 {object} []val.PageInfo
//
// @Router 		/pages [get]
func ListChapterPages(appState *state.AppState) iris.Handler {
	pageApp := appState.PageApp

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		var args val.ListChapterPageArgs

		if err := ctx.ReadQuery(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "查询参数格式错误: "+err.Error())
			return
		}

		result, err := pageApp.List(
			buildReqCx(ctx),
			currentUserID,
			&args,
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
// @Accept 		json
// @Produce 	json
// @Param 		body body val.ReserveChapterPagesArgs true "预留页面参数"
//
// @Success 	201 {object} val.ReserveChapterPagesRes
//
// @Router 		/pages [post]
func ReserveChapterPages(appState *state.AppState) iris.Handler {
	pageApp := appState.PageApp

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		var args val.ReserveChapterPagesArgs

		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误: "+err.Error())
			return
		}

		result, err := pageApp.Reserve(
			buildReqCx(ctx),
			currentUserID,
			&args,
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
// @Accept 		json
// @Produce 	json
// @Param 		page_id path string true "页面 ID"
// @Param 		body body val.UpdatePageArgs true "更新页面参数"
//
// @Success 	200
//
// @Router 		/pages/{page_id} [put]
func UpdatePage(appState *state.AppState) iris.Handler {
	pageApp := appState.PageApp

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		pageID := ctx.Params().Get("page_id")
		if pageID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 page_id 路径参数")
			return
		}

		var args val.UpdatePageArgs

		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误: "+err.Error())
			return
		}

		if args.ID != pageID {
			reject(ctx, iris.StatusBadRequest, "路径参数 page_id 与请求体中的 id 不匹配")
			return
		}

		if err := pageApp.Update(
			buildReqCx(ctx),
			currentUserID,
			&args,
		); err != nil {
			reject(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		accept(ctx, "更新页面成功", nil)
	}
}

// DeletePage godoc
// @Summary 	删除页面
// @Description 删除指定页面
//
// @Tags 		page
// @Security 	ApiKeyAuth
// @Produce 	json
// @Param 		page_id path string true "页面 ID"
//
// @Success 	200
//
// @Router 		/pages/{page_id} [delete]
func DeletePage(appState *state.AppState) iris.Handler {
	pageApp := appState.PageApp

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		pageID := ctx.Params().Get("page_id")
		if pageID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 page_id 路径参数")
			return
		}

		if err := pageApp.Remove(
			buildReqCx(ctx),
			currentUserID,
			pageID,
		); err != nil {
			reject(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		accept(ctx, "删除页面成功", nil)
	}
}
