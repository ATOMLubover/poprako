package http

import (
	"poprako-s/internal/app/val"
	"poprako-s/internal/state"

	"github.com/kataras/iris/v12"
)

// ListPageUnits godoc
// @Summary 	获取页面 unit 列表
// @Description 获取指定页面的所有翻校单元，按 index 升序排列，注意当列表为空时返回 null 而非空数组
//
// @Tags 		unit
// @Security 	ApiKeyAuth
// @Produce 	json
// @Param 		page_id query string true "页面 ID"
//
// @Success 	200 {object} []val.UnitInfo
//
// @Router 		/units [get]
func ListPageUnits(appState *state.AppState) iris.Handler {
	unitApp := appState.UnitApp

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		pageID := ctx.URLParam("page_id")
		if pageID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 page_id 查询参数")
			return
		}

		result, err := unitApp.List(
			buildReqCx(ctx),
			currentUserID,
			pageID,
		)
		if err != nil {
			reject(ctx, iris.StatusForbidden, err.Error())
			return
		}

		accept(ctx, "获取 unit 列表成功", result)
	}
}

// SavePageUnits godoc
// @Summary 	保存页面 unit diff
// @Description 以 diff 语义（insert / patch / delete）保存页面的翻校单元，并同步更新页面和章节的统计字段
//
// @Tags 		unit
// @Security 	ApiKeyAuth
// @Accept 		json
// @Produce 	json
// @Param 		body body val.SavePageUnitArgs true "unit diff 参数"
//
// @Success 	200
//
// @Router 		/units [put]
func SavePageUnits(appState *state.AppState) iris.Handler {
	unitApp := appState.UnitApp

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		var args val.SavePageUnitArgs
		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误: "+err.Error())
			return
		}

		if err := unitApp.Save(
			buildReqCx(ctx),
			currentUserID,
			&args,
		); err != nil {
			reject(ctx, iris.StatusForbidden, err.Error())
			return
		}

		accept(ctx, "保存 unit 成功", nil)
	}
}
