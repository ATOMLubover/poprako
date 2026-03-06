package http

import (
	"labelplus-next-web-be/internal/state"
	"labelplus-next-web-be/internal/value"

	"github.com/kataras/iris/v12"
)

// ListTeamComics godoc
// @Summary 	获取指定汉化组的漫画列表
// @Description 获取指定汉化组的漫画列表，支持分页，注意当列表为空，会返回 null 而不是空数组
//
// @Tags 		comic
// @Security 	ApiKeyAuth
// @Produce 	json
// @Param 		team_id query string true "汉化组 ID"
// @Param 		offset query int true "偏移量"
// @Param 		limit query int true "每页数量"
//
// @Success 	200 {object} []value.ComicInfo
//
// @Router 		/comics [get]
func ListTeamComics(appState *state.AppState) iris.Handler {
	comicApplication := appState.ComicApplication

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrentUserID(ctx)
		if !ok {
			return
		}

		var args value.ListTeamComicArgs
		if err := ctx.ReadQuery(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "查询参数格式错误: "+err.Error())
			return
		}

		result, err := comicApplication.ListTeamComics(
			*buildTraceScope(ctx),
			currentUserID,
			&args,
		)
		if err != nil {
			reject(ctx, iris.StatusForbidden, err.Error())
			return
		}

		accept(ctx, "获取漫画列表成功", result)
	}
}

// CreateComic godoc
// @Summary 	创建漫画
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
			*buildTraceScope(ctx),
			currentUserID,
			&args,
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
// @Summary 	更新漫画
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
// @Router 		/comics/{comic_id} [patch]
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
			*buildTraceScope(ctx),
			currentUserID,
			&args,
		); err != nil {
			reject(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		accept(ctx, "更新漫画成功", nil)
	}
}

// DeleteComic godoc
// @Summary 	删除漫画
// @Description 删除指定漫画
//
// @Tags 		comic
// @Security 	ApiKeyAuth
// @Produce 	json
// @Param 		comic_id path string true "漫画 ID"
// @Param 		team_id query string true "汉化组 ID"
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

		teamID := ctx.URLParam("team_id")
		if teamID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 team_id 查询参数")
			return
		}

		if err := comicApplication.DeleteComic(
			*buildTraceScope(ctx),
			currentUserID,
			teamID,
			comicID,
		); err != nil {
			reject(ctx, iris.StatusForbidden, err.Error())
			return
		}

		accept(ctx, "删除漫画成功", nil)
	}
}
