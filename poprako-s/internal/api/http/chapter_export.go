package http

import (
	"fmt"

	"poprako-s/internal/state"

	"github.com/kataras/iris/v12"
)

// ExportChapter godoc
// @Summary		导出章节数据（JSON 格式）
// @Description 导出指定章节的完整数据，包含页面与翻译单元信息
//
// @Tags		chapter
// @Security	ApiKeyAuth
// @Produce		json
// @Param		chapter_id path string true "章节 ID"
//
// @Success		200 {object} val.ChapterExport
//
// @Router		/chapters/{chapter_id}/export [get]
func ExportChapter(appState *state.AppState) iris.Handler {
	chapterExportApp := appState.ChapterExportApp

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

		result, err := chapterExportApp.ExportChapter(
			buildReqCx(ctx),
			currUserID,
			chapterID,
		)
		if err != nil {
			reject(ctx, iris.StatusForbidden, err.Error())
			return
		}

		accept(ctx, "导出章节数据成功", result)
	}
}

// ExportChapterLp godoc
// @Summary		导出章节数据（LabelPlus 格式）
// @Description 导出指定章节的 LabelPlus 格式文本文件，可直接用于 LabelPlus 工具导入
//
// @Tags		chapter
// @Security	ApiKeyAuth
// @Produce		text/plain
// @Param		chapter_id path string true "章节 ID"
//
// @Success		200 {string} string
//
// @Router		/chapters/{chapter_id}/export/lp [get]
func ExportChapterLp(appState *state.AppState) iris.Handler {
	chapterExportApp := appState.ChapterExportApp

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

		content, err := chapterExportApp.ExportChapterLp(
			buildReqCx(ctx),
			currUserID,
			chapterID,
		)
		if err != nil {
			reject(ctx, iris.StatusForbidden, err.Error())
			return
		}

		fileName := fmt.Sprintf("chapter-%s.lp.txt", chapterID)

		ctx.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
		ctx.ContentType("text/plain; charset=utf-8")
		ctx.StatusCode(iris.StatusOK)

		_, _ = ctx.WriteString(content)
	}
}
