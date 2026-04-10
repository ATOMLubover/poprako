package http

import (
	"poprako-s/internal/app/val"
	"poprako-s/internal/state"

	"github.com/kataras/iris/v12"
)

// ImportChapter godoc
// @Summary 		导入章节数据
// @Description 以 Poprako JSON 或 LabelPlus 文本格式导入章节内容
//
// @Tags 		chapter
// @Security 	ApiKeyAuth
// @Accept 		json
// @Produce 	json
// @Param 		chapter_id path string true "章节 ID"
// @Param 		body body val.ImportChapterArgs true "导入参数"
//
// @Success 	200 {object} val.ImportChapterRes
//
// @Router 		/chapters/{chapter_id}/import [post]
func ImportChapter(appState *state.AppState) iris.Handler {
	chapterImportApp := appState.ChapterImportApp

	return func(ctx iris.Context) {
		currUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		chapterID := ctx.Params().Get("chapter_id")
		if chapterID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 chapter_id 参数")
			return
		}

		var args val.ImportChapterArgs

		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误: "+err.Error())
			return
		}

		args.ChapterID = chapterID

		result, err := chapterImportApp.ImportChapter(
			buildReqCx(ctx),
			currUserID,
			&args,
		)
		if err != nil {
			reject(ctx, iris.StatusForbidden, err.Error())
			return
		}

		accept(ctx, "导入章节数据成功", result)
	}
}
