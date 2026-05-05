package http

import (
	"poprako-s/internal/api/http/res"
	"poprako-s/internal/api/state"
	"poprako-s/internal/app/val"

	"github.com/kataras/iris/v12"
)

// `ListPageUnits` godoc
// @Summary List Page Units
// @Description List units for one page
// @Description The caller must have any assignment on the target chapter
// @Description Auth: `authorization` cookie is preferred over `Authorization` header when both are present
// @Tags unit
// @Security ApiKeyAuth
// @Produce json
// @Param page_id path string true "page id"
// @Success 200 {object} res.HttpRes[val.ListPageUnitsRes]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 403 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /pages/{page_id}/units [get]
func ListPageUnits(st *state.AppState) iris.Handler {
	unitApp := st.UnitApp

	return func(cx iris.Context) {
		currUid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		pageId := cx.Params().Get("page_id")
		if pageId == "" {
			res.Reject(cx, iris.StatusBadRequest, "缺少 page_id 参数")
			return
		}

		re := unitApp.ListByPage(newReqCx(cx), currUid, &val.ListPageUnitsArgs{PageId: pageId})
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

// `SavePageUnits` godoc
// @Summary Save Page Units
// @Description Apply one page unit diff and synchronize page and chapter counters
// @Description
// @Description ## Operations
// @Description Each op in `diff.ops` is classified automatically by the server:
// @Description - **CREATE**: set `local_id`, leave `id` empty  Must include `is_bubble`, `is_proofread`, `x_coord`, `y_coord`
// @Description - **SAVE (update/upsert)**: set `id` plus >=1 mutable field  Must include all geometry fields (`is_bubble`, `is_proofread`, `x_coord`, `y_coord`)
// @Description - **DELETE**: set `id` only, no mutable fields
// @Description Mutable fields: `is_bubble`, `is_proofread`, `x_coord`, `y_coord`, `translated_text`, `translator_comment`, `last_translator_id`, `proofread_text`, `proofreader_comment`, `last_proofreader_id`
// @Description
// @Description ## `cand_order` (client-suggested reindex order)
// @Description An ordered list of unit identifiers to control the final display order after mutations are applied
// @Description Use `local_id` for CREATE ops or real `id` for SAVE ops
// @Description Must include every `local_id` and `id` from CREATE/SAVE ops  Must exclude all deleted ids  No duplicates  No empty strings
// @Description Units not listed in `cand_order` are kept near their original neighbours
// @Description
// @Description ## Validation rules (return 400 on failure)
// @Description - `diff.page_id` must match the path `page_id`
// @Description - Each op must have `local_id` xor `id` (not both, not neither)
// @Description - CREATE ops require `is_bubble`, `is_proofread`, `x_coord`, `y_coord`
// @Description - SAVE ops require all geometry fields plus at least one mutable field
// @Description - The same unit cannot appear in both a CREATE/SAVE op and a DELETE op
// @Description - `cand_order` must not contain ids being deleted
// @Description
// @Description ## Permission
// @Description The caller must be translator or proofreader on the target chapter
// @Description Auth: `authorization` cookie is preferred over `Authorization` header when both are present
// @Tags unit
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param page_id path string true "page id"
// @Param body body val.SavePageUnitsArgs true "save page units args"
// @Success 200 {object} res.HttpRes[val.SavePageUnitsRes]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 403 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /pages/{page_id}/units [post]
func SavePageUnits(st *state.AppState) iris.Handler {
	unitApp := st.UnitApp

	return func(cx iris.Context) {
		currUid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		pageId := cx.Params().Get("page_id")
		if pageId == "" {
			res.Reject(cx, iris.StatusBadRequest, "缺少 page_id 参数")
			return
		}

		var args val.SavePageUnitsArgs
		if err := cx.ReadJSON(&args); err != nil {
			res.Reject(cx, iris.StatusBadRequest, "请求参数解析失败")
			return
		}

		args.PageId = pageId

		re := unitApp.SaveByPage(newReqCx(cx), currUid, &args)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}
