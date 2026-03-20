package http

import (
	"labelplus-next-web-be/internal/state"
	"labelplus-next-web-be/internal/value"

	"github.com/kataras/iris/v12"
)

// ListComics godoc
// @Summary 		获取指定工作集的漫画列表
// @Description 获取指定工作集的漫画列表，支持分页、includes 嵌套信息查询及基于最新章节状态的筛选。
// @Description 注意：所有状态筛选均针对每个漫画的 "最新章节"（即该漫画中 index 最大且未删除的那一章）。
//
// @Tags 		comic
// @Security 	ApiKeyAuth
// @Produce 	json
// @Param 		workset_id query string true "工作集 ID"
// @Param 		offset query int true "偏移量"
// @Param 		limit query int true "每页数量"
// @Param 		"includes[]" query []string false "include 关联信息，可选值：workset, creator"
// @Param 		fuzzy_title query string false "漫画标题模糊搜索（不区分大小写，空字符串时忽略）"
// @Param 		upload_status query string false "最新章节上传状态，可选值: pending, completed"
// @Param 		translate_status query string false "最新章节翻译状态，可选值: pending, in_progress, completed"
// @Param 		proofread_status query string false "最新章节校对状态，可选值: pending, in_progress, completed"
// @Param 		typeset_status query string false "最新章节嵌字状态，可选值: pending, in_progress, completed"
// @Param 		review_status query string false "最新章节审核状态，可选值: pending, completed"
// @Param 		publish_status query string false "最新章节发布状态，可选值: pending, completed"
//
// @Success 	200 {object} []value.ComicInfo
//
// @Router 		/comics [get]
func ListComics(appState *state.AppState) iris.Handler {
	comicApplication := appState.ComicApplication

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrentUserID(ctx)
		if !ok {
			return
		}

		var args value.ListComicArgs
		if err := ctx.ReadQuery(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "查询参数格式错误: "+err.Error())
			return
		}

		result, err := comicApplication.ListComics(
			buildTraceScope(ctx),
			currentUserID,
			args,
		)
		if err != nil {
			reject(ctx, iris.StatusForbidden, err.Error())
			return
		}

		accept(ctx, "获取漫画列表成功", result)
	}
}

// CreateComic godoc
// @Summary 	创建漫画（已测试）
// @Description 在指定汉化组中创建漫画
//
// @Tags 		comic
// @Security 	ApiKeyAuth
// @Accept 	json
// @Produce 	json
// @Param 		body body value.CreateComicArgs true "创建漫画参数"
//
// @Success 	201 {object} value.CreateComicResult
//
// @Router 		/comics [post]
func CreateComic(appState *state.AppState) iris.Handler {
	comicApplication := appState.ComicApplication

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrentUserID(ctx)
		if !ok {
			return
		}

		var args value.CreateComicArgs
		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误: "+err.Error())
			return
		}

		result, err := comicApplication.CreateComic(
			buildTraceScope(ctx),
			currentUserID,
			args,
		)
		if err != nil {
			reject(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		ctx.StatusCode(iris.StatusCreated)
		accept(ctx, "创建漫画成功", result)
	}
}

// PatchComic godoc
// @Summary 	更新漫画（已测试）
// @Description 更新指定漫画的信息
//
// @Tags 		comic
// @Security 	ApiKeyAuth
// @Accept 	json
// @Produce 	json
// @Param 		comic_id path string true "漫画 ID"
// @Param 		body body value.UpdateComicArgs true "更新漫画参数"
//
// @Success 	200
//
// @Router 		/comics/{comic_id} [put]
func PatchComic(appState *state.AppState) iris.Handler {
	comicApplication := appState.ComicApplication

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrentUserID(ctx)
		if !ok {
			return
		}

		comicID := ctx.Params().Get("comic_id")
		if comicID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 comic_id 路径参数")
			return
		}

		var args value.UpdateComicArgs
		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误: "+err.Error())
			return
		}

		if args.ID != comicID {
			reject(ctx, iris.StatusBadRequest, "路径参数 comic_id 与请求体中的 ID 不匹配")
			return
		}

		if err := comicApplication.UpdateComic(
			buildTraceScope(ctx),
			currentUserID,
			args,
		); err != nil {
			reject(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		accept(ctx, "更新漫画成功", nil)
	}
}

// DeleteComic godoc
// @Summary 	删除漫画（已测试）
// @Description 删除指定漫画
//
// @Tags 		comic
// @Security 	ApiKeyAuth
// @Produce 	json
// @Param 		comic_id path string true "漫画 ID"
//
// @Success 	200
//
// @Router 		/comics/{comic_id} [delete]
func DeleteComic(appState *state.AppState) iris.Handler {
	comicApplication := appState.ComicApplication

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrentUserID(ctx)
		if !ok {
			return
		}

		comicID := ctx.Params().Get("comic_id")
		if comicID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 comic_id 路径参数")
			return
		}

		if err := comicApplication.DeleteComic(
			buildTraceScope(ctx),
			currentUserID,
			comicID,
		); err != nil {
			reject(ctx, iris.StatusForbidden, err.Error())
			return
		}

		accept(ctx, "删除漫画成功", nil)
	}
}

// GetComicCover godoc
// @Summary 		获取漫画封面
// @Description 返回最新（index 最大且未删除）的章节第一页封面链接；若不存在则返回 null
//
// @Tags 		comic
// @Security 	ApiKeyAuth
// @Produce 	json
// @Param 		comic_id path string true "漫画 ID"
//
// @Success 	200 {object} value.ComicCoverResult
//
// @Router 		/comics/{comic_id}/cover [get]
func GetComicCover(appState *state.AppState) iris.Handler {
	comicApplication := appState.ComicApplication

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrentUserID(ctx)
		if !ok {
			return
		}

		comicID := ctx.Params().Get("comic_id")
		if comicID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 comic_id 路径参数")
			return
		}

		coverURL, err := comicApplication.GetComicCover(
			buildTraceScope(ctx),
			currentUserID,
			comicID,
		)
		if err != nil {
			reject(ctx, iris.StatusForbidden, err.Error())
			return
		}

		if coverURL == "" {
			accept(ctx, "获取漫画封面成功", nil)
			return
		}

		accept(ctx, "获取漫画封面成功", value.NewComicCoverResult(coverURL))
	}
}
