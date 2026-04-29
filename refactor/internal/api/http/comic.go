package http

import (
	"strconv"

	"poprako-s/internal/api/http/res"
	"poprako-s/internal/api/state"
	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/model/enum"

	"github.com/kataras/iris/v12"
)

// `ListComics` godoc
// @Summary List Comics
//
//	List active comics for one workset
//	The caller must be a member of the target workset team
//	Auth: `authorization` cookie is preferred over `Authorization` header when both are present
//
// @Tags comic
// @Security ApiKeyAuth
// @Produce json
// @Param workset_id path string true "workset id"
// @Param fuzzy_title query string false "fuzzy title"
// @Param upload_phase query int false "upload phase 0 pending 1 ongoing 2 completed"
// @Param translate_phase query int false "translate phase 0 pending 1 ongoing 2 completed"
// @Param proofread_phase query int false "proofread phase 0 pending 1 ongoing 2 completed"
// @Param typeset_phase query int false "typeset phase 0 pending 1 ongoing 2 completed"
// @Param review_phase query int false "review phase 0 pending 1 ongoing 2 completed"
// @Param publish_phase query int false "publish phase 0 pending 1 ongoing 2 completed"
// @Param offset query int false "pagination offset"
// @Param limit query int false "pagination limit"
// @Success 200 {object} res.HttpRes "res.HttpRes{data=[]val.ComicVal}"
// @Failure 400 {object} res.HttpRes
// @Failure 401 {object} res.HttpRes
// @Failure 500 {object} res.HttpRes
// @Router /comic/workset/{workset_id} [get]
func ListComics(st *state.AppState) iris.Handler {
	comicApp := st.ComicApp

	return func(cx iris.Context) {
		currUid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		worksetId := cx.Params().Get("workset_id")
		if worksetId == "" {
			res.Reject(cx, iris.StatusBadRequest, "缺少 workset_id 参数")
			return
		}

		offset, err := cx.URLParamInt("offset")
		if err != nil {
			res.Reject(cx, iris.StatusBadRequest, "offset 参数格式错误")
			return
		}

		limit, err := cx.URLParamInt("limit")
		if err != nil {
			res.Reject(cx, iris.StatusBadRequest, "limit 参数格式错误")
			return
		}

		uploadPhase, err := parseWorkflowPhaseParam(cx, "upload_phase")
		if err != nil {
			res.Reject(cx, iris.StatusBadRequest, "upload_phase 参数格式错误")
			return
		}

		translatePhase, err := parseWorkflowPhaseParam(cx, "translate_phase")
		if err != nil {
			res.Reject(cx, iris.StatusBadRequest, "translate_phase 参数格式错误")
			return
		}

		proofreadPhase, err := parseWorkflowPhaseParam(cx, "proofread_phase")
		if err != nil {
			res.Reject(cx, iris.StatusBadRequest, "proofread_phase 参数格式错误")
			return
		}

		typesetPhase, err := parseWorkflowPhaseParam(cx, "typeset_phase")
		if err != nil {
			res.Reject(cx, iris.StatusBadRequest, "typeset_phase 参数格式错误")
			return
		}

		reviewPhase, err := parseWorkflowPhaseParam(cx, "review_phase")
		if err != nil {
			res.Reject(cx, iris.StatusBadRequest, "review_phase 参数格式错误")
			return
		}

		publishPhase, err := parseWorkflowPhaseParam(cx, "publish_phase")
		if err != nil {
			res.Reject(cx, iris.StatusBadRequest, "publish_phase 参数格式错误")
			return
		}

		args := &val.ListComicArgs{
			WorksetId:      worksetId,
			FuzzyTitle:     cx.URLParamDefault("fuzzy_title", ""),
			UploadPhase:    uploadPhase,
			TranslatePhase: translatePhase,
			ProofreadPhase: proofreadPhase,
			TypesetPhase:   typesetPhase,
			ReviewPhase:    reviewPhase,
			PublishPhase:   publishPhase,
			Offset:         offset,
			Limit:          limit,
		}

		re := comicApp.List(newReqCx(cx), currUid, args)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

// `parseWorkflowPhaseParam` parses one workflow phase query parameter.
func parseWorkflowPhaseParam(cx iris.Context, key string) (*enum.WorkflowPhase, error) {
	phaseVal := cx.URLParamDefault(key, "")
	if phaseVal == "" {
		return nil, nil
	}

	v, err := strconv.Atoi(phaseVal)
	if err != nil {
		return nil, err
	}

	phase := enum.WorkflowPhase(v)
	if phase < enum.WorkflowPending || phase > enum.WorkflowCompleted {
		return nil, strconv.ErrSyntax
	}

	return &phase, nil
}

// `GetComicById` godoc
// @Summary Get Comic By Id
//
//	Get one comic by id.
//	The caller must be a member of the owning team.
//	Auth: `authorization` cookie is preferred over `Authorization` header when both are present.
//
// @Tags comic
// @Security ApiKeyAuth
// @Produce json
// @Param comic_id path string true "comic id"
// @Success 200 {object} res.HttpRes "res.HttpRes{data=val.ComicVal}"
// @Failure 400 {object} res.HttpRes
// @Failure 404 {object} res.HttpRes
// @Failure 401 {object} res.HttpRes
// @Failure 500 {object} res.HttpRes
// @Router /comic/{comic_id} [get]
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

		re := comicApp.GetById(newReqCx(cx), currUid, comicId)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

// `CreateComic` godoc
// @Summary Create Comic
//
//	Create one comic under a workset
//	The caller must be an admin of the target workset team
//	Auth: `authorization` cookie is preferred over `Authorization` header when both are present
//
// @Tags comic
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param body body val.CreateComicArgs true "create comic args"
// @Success 201 {object} res.HttpRes "res.HttpRes{data=val.ComicCreatedRes}"
// @Failure 400 {object} res.HttpRes
// @Failure 401 {object} res.HttpRes
// @Failure 500 {object} res.HttpRes
// @Router /comic [post]
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
//
//	Update comic fields with put semantics
//	The caller must be an admin of the owning team
//	Auth: `authorization` cookie is preferred over `Authorization` header when both are present
//
// @Tags comic
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param comic_id path string true "comic id"
// @Param body body val.ComicUpdArgs true "update comic args"
// @Success 200 {object} res.HttpRes
// @Failure 400 {object} res.HttpRes
// @Failure 401 {object} res.HttpRes
// @Failure 500 {object} res.HttpRes
// @Router /comic/{comic_id} [put]
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

// `RemoveComic` godoc
// @Summary Remove Comic
//
//	Soft-delete one comic by id
//	The caller must be an admin of the owning team
//	Auth: `authorization` cookie is preferred over `Authorization` header when both are present
//
// @Tags comic
// @Security ApiKeyAuth
// @Produce json
// @Param comic_id path string true "comic id"
// @Success 200 {object} res.HttpRes
// @Failure 400 {object} res.HttpRes
// @Failure 401 {object} res.HttpRes
// @Failure 500 {object} res.HttpRes
// @Router /comic/{comic_id} [delete]
func RemoveComic(st *state.AppState) iris.Handler {
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

		re := comicApp.Remove(newReqCx(cx), currUid, comicId)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}
