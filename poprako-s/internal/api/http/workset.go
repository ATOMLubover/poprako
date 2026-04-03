package http

import (
	"poprako-s/internal/app/val"
	"poprako-s/internal/state"

	"github.com/kataras/iris/v12"
)

// ListWorksets godoc
// @Summary 	获取指定汉化组的工作集列表
// @Description 获取指定汉化组的工作集列表，支持分页，注意当列表为空，会返回 null 而不是空数组
//
// @Tags 		workset
// @Security 	ApiKeyAuth
// @Produce 	json
// @Param 		team_id query string true "汉化组 ID"
// @Param 		"includes[]" query []string false "include 关联信息，可选值：team"
// @Param 		offset query int true "偏移量"
// @Param 		limit query int true "每页数量"
//
// @Success 	200 {object} []val.WorksetInfo
//
// @Router 		/worksets [get]
func ListWorksets(appState *state.AppState) iris.Handler {
	worksetApp := appState.WorksetApp

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		var args val.ListWorksetArgs
		if err := ctx.ReadQuery(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "查询参数格式错误: "+err.Error())
			return
		}

		result, err := worksetApp.List(
			buildReqCx(ctx),
			currentUserID,
			&args,
		)
		if err != nil {
			reject(ctx, iris.StatusForbidden, err.Error())
			return
		}

		accept(ctx, "获取工作集列表成功", result)
	}
}

// CreateWorkset godoc
// @Summary 	创建工作集
// @Description 在指定汉化组中创建工作集
//
// @Tags 		workset
// @Security 	ApiKeyAuth
// @Accept 		json
// @Produce 	json
// @Param 		body body val.CreateWorksetArgs true "创建工作集参数"
//
// @Success 	201 {object} val.CreateWorksetRes
//
// @Router 		/worksets [post]
func CreateWorkset(appState *state.AppState) iris.Handler {
	worksetApp := appState.WorksetApp

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		var args val.CreateWorksetArgs
		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误: "+err.Error())
			return
		}

		result, err := worksetApp.Create(
			buildReqCx(ctx),
			currentUserID,
			&args,
		)
		if err != nil {
			reject(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		ctx.StatusCode(iris.StatusCreated)
		accept(ctx, "创建工作集成功", result)
	}
}

// UpdateWorkset godoc
// @Summary 	更新工作集
// @Description 更新指定工作集的信息
//
// @Tags 		workset
// @Security 	ApiKeyAuth
// @Accept 		json
// @Produce 	json
// @Param 		workset_id path string true "工作集 ID"
// @Param 		body body val.UpdateWorksetArgs true "更新工作集参数"
//
// @Success 	200
//
// @Router 		/worksets/{workset_id} [put]
func UpdateWorkset(appState *state.AppState) iris.Handler {
	worksetApp := appState.WorksetApp

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		worksetID := ctx.Params().Get("workset_id")
		if worksetID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 workset_id 路径参数")
			return
		}

		var args val.UpdateWorksetArgs
		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误: "+err.Error())
			return
		}

		if args.ID != worksetID {
			reject(ctx, iris.StatusBadRequest, "路径参数 workset_id 与请求体中的 ID 不匹配")
			return
		}

		if err := worksetApp.Update(
			buildReqCx(ctx),
			currentUserID,
			&args,
		); err != nil {
			reject(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		accept(ctx, "更新工作集成功", nil)
	}
}

// DeleteWorkset godoc
// @Summary 	删除工作集
// @Description 删除指定工作集
//
// @Tags 		workset
// @Security 	ApiKeyAuth
// @Produce 	json
// @Param 		workset_id path string true "工作集 ID"
//
// @Success 	200
//
// @Router 		/worksets/{workset_id} [delete]
func DeleteWorkset(appState *state.AppState) iris.Handler {
	worksetApp := appState.WorksetApp

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		worksetID := ctx.Params().Get("workset_id")
		if worksetID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 workset_id 路径参数")
			return
		}

		if err := worksetApp.Remove(
			buildReqCx(ctx),
			currentUserID,
			worksetID,
		); err != nil {
			reject(ctx, iris.StatusForbidden, err.Error())
			return
		}

		accept(ctx, "删除工作集成功", nil)
	}
}
