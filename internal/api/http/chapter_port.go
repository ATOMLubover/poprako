package http

import (
	"fmt"

	"poprako-s/internal/api/http/res"
	"poprako-s/internal/api/state"
	"poprako-s/internal/app/val"

	"github.com/kataras/iris/v12"
)

// `ExportChapter` godoc
// @Summary Export Chapter
// @Description Export one chapter in JSON format
// @Description The caller must have any assignment on the target chapter
// @Description Auth: `authorization` cookie is preferred over `Authorization` header when both are present
// @Tags chapter
// @Security ApiKeyAuth
// @Produce json
// @Param chapter_id path string true "chapter id"
// @Success 200 {object} res.HttpRes[val.ChapterExportVal]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 403 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /chapters/{chapter_id}/export [get]
func ExportChapter(st *state.AppState) iris.Handler {
	chapterPortApp := st.ChapterPortApp

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

		re := chapterPortApp.Export(newReqCx(cx), currUid, chapterId)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

// `ExportChapterLp` godoc
// @Summary Export Chapter LabelPlus
// @Description Export one chapter in LabelPlus text format
// @Description The caller must have any assignment on the target chapter
// @Description Auth: `authorization` cookie is preferred over `Authorization` header when both are present
// @Tags chapter
// @Security ApiKeyAuth
// @Produce text/plain
// @Param chapter_id path string true "chapter id"
// @Success 200 {string} string
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 403 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /chapters/{chapter_id}/export/lp [get]
func ExportChapterLp(st *state.AppState) iris.Handler {
	chapterPortApp := st.ChapterPortApp

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

		re := chapterPortApp.ExportLp(newReqCx(cx), currUid, chapterId)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		fileName := fmt.Sprintf("chapter-%s.lp.txt", chapterId)

		cx.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
		cx.ContentType("text/plain; charset=utf-8")
		cx.StatusCode(iris.StatusOK)

		if re.Data() != nil {
			_, _ = cx.WriteString(*re.Data())
		}
	}
}

// `ImportChapter` godoc
// @Summary Import Chapter
// @Description Import one chapter from Poprako JSON or LabelPlus text content
// @Description The caller must be translator or proofreader of the target chapter
// @Description Auth: `authorization` cookie is preferred over `Authorization` header when both are present
// @Tags chapter
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param chapter_id path string true "chapter id"
// @Param body body val.ImportChapterBody true "import chapter args"
// @Success 200 {object} res.HttpRes[val.ImportChapterRes]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 403 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /chapters/{chapter_id}/import [post]
func ImportChapter(st *state.AppState) iris.Handler {
	chapterPortApp := st.ChapterPortApp

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

		var args val.ImportChapterArgs
		if err := cx.ReadJSON(&args); err != nil {
			res.Reject(cx, iris.StatusBadRequest, "请求参数解析失败")
			return
		}

		args.ChapterId = chapterId

		re := chapterPortApp.Import(newReqCx(cx), currUid, &args)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}
