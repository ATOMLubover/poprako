package http

import (
	"poprako-s/internal/app/val"
	"poprako-s/internal/state"

	"github.com/kataras/iris/v12"
)

// GetChapterByID godoc
// @Summary 	根据 ID 获取章节详情
// @Description 根据章节 ID 获取单个章节的详细信息
//
// @Tags 		chapter
// @Security 	ApiKeyAuth
// @Produce 	json
// @Param 		chapter_id path string true "章节 ID"
//
// @Success 	200 {object} val.ChapterInfo
//
// @Router 		/chapters/{chapter_id} [get]
func GetChapterByID(appState *state.AppState) iris.Handler {
	chapterApp := appState.ChapterApp

	return func(ctx iris.Context) {
		currUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		chapterID := ctx.Params().Get("chapter_id")
		if chapterID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 chapter_id 路径参数")
			return
		}

		result, err := chapterApp.Get(
			buildReqCx(ctx),
			currUserID,
			chapterID,
		)
		if err != nil {
			reject(ctx, iris.StatusForbidden, err.Error())
			return
		}

		accept(ctx, "获取章节详情成功", result)
	}
}

// GetComicPinnedChapter godoc
// @Summary 	获取漫画置顶章节
// @Description 获取指定漫画的置顶章节信息；若尚无置顶章节则返回 null
//
// @Tags 		chapter
// @Security 	ApiKeyAuth
// @Produce 	json
// @Param 		comic_id path string true "漫画 ID"
//
// @Success 	200 {object} val.ChapterInfo
//
// @Router 		/comics/{comic_id}/pinned-chapter [get]
func GetComicPinnedChapter(appState *state.AppState) iris.Handler {
	chapterApp := appState.ChapterApp

	return func(ctx iris.Context) {
		currUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		comicID := ctx.Params().Get("comic_id")
		if comicID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 comic_id 路径参数")
			return
		}

		result, err := chapterApp.GetComicPinned(
			buildReqCx(ctx),
			currUserID,
			comicID,
		)
		if err != nil {
			reject(ctx, iris.StatusForbidden, err.Error())
			return
		}

		accept(ctx, "获取置顶章节成功", result)
	}
}

// ListComicChapters godoc
// @Summary 	获取漫画章节列表
// @Description 获取指定漫画的章节列表，支持分页和 includes 嵌套信息查询
//
// @Tags 		chapter
// @Security 	ApiKeyAuth
// @Produce 	json
// @Param 		comic_id query string true "漫画 ID"
// @Param 		offset query int true "偏移量"
// @Param 		limit query int true "每页数量"
// @Param 		"includes" query []string false "include 关联信息，可选值：creator"
//
// @Success 	200 {object} []val.ChapterInfo
//
// @Router 		/chapters [get]
func ListComicChapters(appState *state.AppState) iris.Handler {
	chapterApp := appState.ChapterApp

	return func(ctx iris.Context) {
		currUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		var args val.ListChapterArgs

		if err := ctx.ReadQuery(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "查询参数格式错误: "+err.Error())
			return
		}

		result, err := chapterApp.List(
			buildReqCx(ctx),
			currUserID,
			&args,
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
// @Description 在指定漫画中创建章节，并写入章节副标题
//
// @Tags 		chapter
// @Security 	ApiKeyAuth
// @Accept 		json
// @Produce 	json
// @Param 		body body val.CreateChapterArgs true "创建章节参数"
//
// @Success 	201 {object} val.CreateChapterRes
//
// @Router 		/chapters [post]
func CreateComicChapter(appState *state.AppState) iris.Handler {
	chapterApp := appState.ChapterApp

	return func(ctx iris.Context) {
		currUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		var args val.CreateChapterArgs

		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误: "+err.Error())
			return
		}

		result, err := chapterApp.Create(
			buildReqCx(ctx),
			currUserID,
			&args,
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
// @Summary 	更新章节
// @Description 局部更新指定章节的信息，包括 subtitle 与工作流状态；未传的字段不会被修改；除了 reviewer 以外，其他任何角色都只能修改自己对应的 workflow 的状态
//
// @Tags 		chapter
// @Security 	ApiKeyAuth
// @Accept 		json
// @Produce 	json
// @Param 		chapter_id path string true "章节 ID"
// @Param 		body body val.UpdateChapterArgs true "更新章节参数"
//
// @Success 	200
//
// @Router 		/chapters/{chapter_id} [patch]
func UpdateChapter(appState *state.AppState) iris.Handler {
	chapterApp := appState.ChapterApp

	return func(ctx iris.Context) {
		currUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		chapterID := ctx.Params().Get("chapter_id")
		if chapterID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 chapter_id 路径参数")
			return
		}

		var args val.UpdateChapterArgs

		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误: "+err.Error())
			return
		}

		args.ChapterID = chapterID

		if err := chapterApp.Update(
			buildReqCx(ctx),
			currUserID,
			&args,
		); err != nil {
			reject(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		accept(ctx, "更新章节成功", nil)
	}
}

// InviteChapterAssignee godoc
// @Summary 	创建章节协作邀请
// @Description 在指定章节下创建协作邀请并返回邀请码
//
// @Tags 		chapter
// @Security 	ApiKeyAuth
// @Accept 		json
// @Produce 	json
// @Param 		chapter_id path string true "章节 ID"
// @Param 		body body val.InviteChapterAssigneeArgs true "邀请参数"
//
// @Success 	200 {object} val.InviteChapterAssigneeRes
//
// @Router 		/chapters/{chapter_id}/invitations [post]
func InviteChapterAssignee(appState *state.AppState) iris.Handler {
	chapterApp := appState.ChapterApp

	return func(ctx iris.Context) {
		currUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		chapterID := ctx.Params().Get("chapter_id")
		if chapterID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 chapter_id 路径参数")
			return
		}

		var args val.InviteChapterAssigneeArgs

		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误: "+err.Error())
			return
		}

		args.ChapterID = chapterID

		res, err := chapterApp.InviteAssignee(
			buildReqCx(ctx),
			currUserID,
			&args,
		)
		if err != nil {
			reject(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		accept(ctx, "创建章节邀请成功", res)
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
	chapterApp := appState.ChapterApp

	return func(ctx iris.Context) {
		currUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		chapterID := ctx.Params().Get("chapter_id")
		if chapterID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 chapter_id 路径参数")
			return
		}

		if err := chapterApp.Remove(
			buildReqCx(ctx),
			currUserID,
			chapterID,
		); err != nil {
			reject(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		accept(ctx, "删除章节成功", nil)
	}
}
