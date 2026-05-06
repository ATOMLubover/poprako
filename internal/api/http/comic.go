package http

import (
	"poprako-s/internal/api/http/res"
	"poprako-s/internal/api/state"
	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/model/enum"

	"github.com/kataras/iris/v12"
)

// `ListComics` godoc
// @Summary List Comics
// @Description List active comics for one workset
// @Description The caller must be a member of the target workset team
// @Description Auth: `authorization` cookie is preferred over `Authorization` header when both are present
// @Tags comic
// @Security ApiKeyAuth
// @Produce json
// @Param workset_id query string true "workset id"
// @Param fuzzy_title query string false "fuzzy title"
// @Param upload_phase query int false "upload phase 0 pending 1 ongoing 2 completed"
// @Param translate_phase query int false "translate phase 0 pending 1 ongoing 2 completed"
// @Param proofread_phase query int false "proofread phase 0 pending 1 ongoing 2 completed"
// @Param typeset_phase query int false "typeset phase 0 pending 1 ongoing 2 completed"
// @Param review_phase query int false "review phase 0 pending 1 ongoing 2 completed"
// @Param publish_phase query int false "publish phase 0 pending 1 ongoing 2 completed"
// @Param includes query []string false "include related fields, optional: workset, workset.team, creator"
// @Param offset query int false "pagination offset"
// @Param limit query int false "pagination limit"
// @Success 200 {object} res.HttpRes[[]val.ComicVal]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /api/v1/comics [get]
func ListComics(st *state.AppState) iris.Handler {
	comicApp := st.ComicApp

	return func(cx iris.Context) {
		currUid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		var args val.ListComicArgs
		if err := cx.ReadQuery(&args); err != nil {
			res.Reject(cx, iris.StatusBadRequest, "请求参数解析失败")
			return
		}

		if args.WorksetId == "" {
			res.Reject(cx, iris.StatusBadRequest, "缺少 workset_id 参数")
			return
		}

		if !parseComicIncludes(args.Includes) {
			res.Reject(cx, iris.StatusBadRequest, "includes 参数格式错误")
			return
		}

		if !validateWorkflowPhase(args.UploadPhase) {
			res.Reject(cx, iris.StatusBadRequest, "upload_phase 参数格式错误")
			return
		}

		if !validateWorkflowPhase(args.TranslatePhase) {
			res.Reject(cx, iris.StatusBadRequest, "translate_phase 参数格式错误")
			return
		}

		if !validateWorkflowPhase(args.ProofreadPhase) {
			res.Reject(cx, iris.StatusBadRequest, "proofread_phase 参数格式错误")
			return
		}

		if !validateWorkflowPhase(args.TypesetPhase) {
			res.Reject(cx, iris.StatusBadRequest, "typeset_phase 参数格式错误")
			return
		}

		if !validateWorkflowPhase(args.ReviewPhase) {
			res.Reject(cx, iris.StatusBadRequest, "review_phase 参数格式错误")
			return
		}

		if !validateWorkflowPhase(args.PublishPhase) {
			res.Reject(cx, iris.StatusBadRequest, "publish_phase 参数格式错误")
			return
		}

		re := comicApp.List(newReqCx(cx), currUid, &args)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

// `validateWorkflowPhase` validates workflow phase filter values.
func validateWorkflowPhase(phase *enum.WorkflowPhase) bool {
	if phase == nil {
		return true
	}

	if *phase < enum.WorkflowPending || *phase > enum.WorkflowCompleted {
		return false
	}

	return true
}

// `GetComicById` godoc
// @Summary Get Comic By Id
// @Description Get one comic by id
// @Description The caller must be a member of the owning team
// @Description Auth: `authorization` cookie is preferred over `Authorization` header when both are present
// @Tags comic
// @Security ApiKeyAuth
// @Produce json
// @Param comic_id path string true "comic id"
// @Param includes query []string false "include related fields, optional: workset, workset.team, creator"
// @Success 200 {object} res.HttpRes[val.ComicVal]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 404 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /api/v1/comics/{comic_id} [get]
func GetComicById(st *state.AppState) iris.Handler {
	comicApp := st.ComicApp

	return func(cx iris.Context) {
		currUid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		comicId := cx.Params().Get("comic_id")
		if comicId == "" {
			res.Reject(cx, iris.StatusBadRequest, "缺少 comic_id 参数")
			return
		}

		var args val.GetComicByIdArgs
		if err := cx.ReadQuery(&args); err != nil {
			res.Reject(cx, iris.StatusBadRequest, "请求参数解析失败")
			return
		}

		args.ComicId = comicId

		if !parseComicIncludes(args.Includes) {
			res.Reject(cx, iris.StatusBadRequest, "includes 参数格式错误")
			return
		}

		re := comicApp.GetById(newReqCx(cx), currUid, &args)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

// `parseComicIncludes` validates include query values.
func parseComicIncludes(includes []enum.ComicIncl) bool {
	if len(includes) == 0 {
		return true
	}

	for i := range includes {
		switch includes[i] {
		case enum.ComicInclWorkset,
			enum.ComicInclWorksetTeam,
			enum.ComicInclCreator:
			continue

		default:
			return false
		}
	}

	return true
}

// `CreateComic` godoc
// @Summary Create Comic
// @Description Create one comic under a workset
// @Description The caller must be an admin of the target workset team
// @Description Auth: `authorization` cookie is preferred over `Authorization` header when both are present
// @Tags comic
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param body body val.CreateComicArgs true "create comic args"
// @Success 201 {object} res.HttpRes[val.ComicCreatedRes]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /api/v1/comics [post]
func CreateComic(st *state.AppState) iris.Handler {
	comicApp := st.ComicApp

	return func(cx iris.Context) {
		currUid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		var args val.CreateComicArgs

		if err := cx.ReadJSON(&args); err != nil {
			res.Reject(cx, iris.StatusBadRequest, "请求参数解析失败")
			return
		}

		re := comicApp.Create(newReqCx(cx), currUid, &args)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusCreated, re.Data())
	}
}

// `UpdateComic` godoc
// @Summary Update Comic
// @Description Update comic fields with PUT semantics
// @Description The caller must be an admin of the owning team
// @Description Auth: `authorization` cookie is preferred over `Authorization` header when both are present
// @Tags comic
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param comic_id path string true "comic id"
// @Param body body val.ComicUpdArgs true "update comic args"
// @Success 200 {object} res.HttpRes[any]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /api/v1/comics/{comic_id} [put]
func UpdateComic(st *state.AppState) iris.Handler {
	comicApp := st.ComicApp

	return func(cx iris.Context) {
		currUid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		comicId := cx.Params().Get("comic_id")
		if comicId == "" {
			res.Reject(cx, iris.StatusBadRequest, "缺少 comic_id 参数")
			return
		}

		var args val.ComicUpdArgs

		if err := cx.ReadJSON(&args); err != nil {
			res.Reject(cx, iris.StatusBadRequest, "请求参数解析失败")
			return
		}

		args.Id = comicId

		re := comicApp.Update(newReqCx(cx), currUid, &args)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

// `ResvComicCover` godoc
// @Summary Reserve Comic Cover Upload
// @Description Reserve a signed upload URL for one comic cover
// @Description The caller must be an admin of the owning team
// @Description Auth: `authorization` cookie is preferred over `Authorization` header when both are present
// @Tags comic
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param comic_id path string true "comic id"
// @Param body body val.ResvComicCoverBody true "reserve comic cover args"
// @Success 200 {object} res.HttpRes[val.ResvComicCoverRes]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /api/v1/comics/{comic_id}/cover [post]
func ResvComicCover(st *state.AppState) iris.Handler {
	comicApp := st.ComicApp

	return func(cx iris.Context) {
		currUid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		comicId := cx.Params().Get("comic_id")
		if comicId == "" {
			res.Reject(cx, iris.StatusBadRequest, "缺少 comic_id 参数")
			return
		}

		var args val.ResvComicCoverArgs
		if err := cx.ReadJSON(&args); err != nil {
			res.Reject(cx, iris.StatusBadRequest, "请求参数解析失败")
			return
		}

		args.ComicId = comicId

		re := comicApp.ResvCover(newReqCx(cx), currUid, &args)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

// `MarkComicCoverUploaded` godoc
// @Summary Confirm Comic Cover Uploaded
// @Description Confirm one comic cover upload after client upload completed
// @Description The caller must be an admin of the owning team
// @Description Auth: `authorization` cookie is preferred over `Authorization` header when both are present
// @Tags comic
// @Security ApiKeyAuth
// @Produce json
// @Param comic_id path string true "comic id"
// @Success 200 {object} res.HttpRes[any]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /api/v1/comics/{comic_id}/cover/confirm [post]
func MarkComicCoverUploaded(st *state.AppState) iris.Handler {
	comicApp := st.ComicApp

	return func(cx iris.Context) {
		currUid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		comicId := cx.Params().Get("comic_id")
		if comicId == "" {
			res.Reject(cx, iris.StatusBadRequest, "缺少 comic_id 参数")
			return
		}

		re := comicApp.MarkCoverUploaded(newReqCx(cx), currUid, comicId)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

// `DeleteComic` godoc
// @Summary Delete Comic
// @Description Hard-delete one comic by id
// @Description The caller must be an admin of the owning team
// @Description Auth: `authorization` cookie is preferred over `Authorization` header when both are present
// @Tags comic
// @Security ApiKeyAuth
// @Produce json
// @Param comic_id path string true "comic id"
// @Success 200 {object} res.HttpRes[any]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /api/v1/comics/{comic_id} [delete]
func DeleteComic(st *state.AppState) iris.Handler {
	comicApp := st.ComicApp

	return func(cx iris.Context) {
		currUid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		comicId := cx.Params().Get("comic_id")
		if comicId == "" {
			res.Reject(cx, iris.StatusBadRequest, "缺少 comic_id 参数")
			return
		}

		re := comicApp.Delete(newReqCx(cx), currUid, comicId)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}
