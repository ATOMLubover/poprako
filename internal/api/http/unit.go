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
// @Router /page/{page_id}/units [get]
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
// @Router /page/{page_id}/units [post]
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
