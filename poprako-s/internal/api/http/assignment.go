package http

import (
	"poprako-s/internal/app/val"
	"poprako-s/internal/state"

	"github.com/kataras/iris/v12"
)

// ListChapterAssignments godoc
// @Summary 	获取章节分配列表
// @Description 获取指定章节的所有分配记录，仅汉化组成员可访问
//
// @Tags 		assignment
// @Security 	ApiKeyAuth
// @Produce 	json
// @Param 		chapter_id query string true "章节 ID"
// @Param 		"includes[]" query []string false "include 关联信息，可选值：user, chapter, chapter.comic, chapter.creator"
// @Param 		offset query int true "偏移量"
// @Param 		limit query int true "每页数量"
//
// @Success 	200 {object} []val.AssignmentInfo
//
// @Router 		/assignments [get]
func ListChapterAssignments(appState *state.AppState) iris.Handler {
	assignmentApp := appState.AssignmentApp

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		var args val.ListChapterAssignmentArgs

		if err := ctx.ReadQuery(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "查询参数格式错误: "+err.Error())
			return
		}

		result, err := assignmentApp.ListByChapter(
			buildReqCx(ctx),
			currentUserID,
			&args,
		)
		if err != nil {
			reject(ctx, iris.StatusForbidden, err.Error())
			return
		}

		accept(ctx, "获取分配列表成功", result)
	}
}

// ListMyAssignments godoc
// @Summary 	获取我的分配列表
// @Description 获取当前登录用户的所有分配记录，支持分页
//
// @Tags 		assignment
// @Security 	ApiKeyAuth
// @Produce 	json
// @Param 		"includes[]" query []string false "include 关联信息，可选值：chapter, chapter.comic, chapter.creator"
// @Param 		offset query int true "偏移量"
// @Param 		limit query int true "每页数量"
//
// @Success 	200 {object} []val.AssignmentInfo
//
// @Router 		/assignments/mine [get]
func ListMyAssignments(appState *state.AppState) iris.Handler {
	assignmentApp := appState.AssignmentApp

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		var args val.ListMyAssignmentArgs

		if err := ctx.ReadQuery(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "查询参数格式错误: "+err.Error())
			return
		}

		result, err := assignmentApp.ListMy(
			buildReqCx(ctx),
			currentUserID,
			&args,
		)
		if err != nil {
			reject(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		accept(ctx, "获取我的分配列表成功", result)
	}
}

// CreateChapterAssignment godoc
// @Summary 	创建章节分配
// @Description 为指定章节创建一条分配记录，需要当前用户在该章节中拥有 reviewer 角色
//
// @Tags 		assignment
// @Security 	ApiKeyAuth
// @Accept 		json
// @Produce 	json
// @Param 		body body val.CreateAssignmentArgs true "创建分配参数"
//
// @Success 	201 {object} val.CreateAssignmentRes
//
// @Router 		/assignments [post]
func CreateChapterAssignment(appState *state.AppState) iris.Handler {
	assignmentApp := appState.AssignmentApp

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		var args val.CreateAssignmentArgs

		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误: "+err.Error())
			return
		}

		result, err := assignmentApp.Create(
			buildReqCx(ctx),
			currentUserID,
			&args,
		)
		if err != nil {
			reject(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		ctx.StatusCode(iris.StatusCreated)
		accept(ctx, "创建分配成功", result)
	}
}

// UpdateAssignment godoc
// @Summary 	更新分配角色
// @Description 全量替换指定分配的角色（PUT 语义），需要当前用户在该章节中拥有 reviewer 角色
//
// @Tags 		assignment
// @Security 	ApiKeyAuth
// @Accept 		json
// @Produce 	json
// @Param 		assignment_id path string true "分配 ID"
// @Param 		body body val.UpdateAssignmentArgs true "更新分配参数"
//
// @Success 	200
//
// @Router 		/assignments/{assignment_id} [put]
func UpdateAssignment(appState *state.AppState) iris.Handler {
	assignmentApp := appState.AssignmentApp

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		assignmentID := ctx.Params().Get("assignment_id")
		if assignmentID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 assignment_id 路径参数")
			return
		}

		var args val.UpdateAssignmentArgs

		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误: "+err.Error())
			return
		}
		args.ID = assignmentID

		if err := assignmentApp.Update(
			buildReqCx(ctx),
			currentUserID,
			&args,
		); err != nil {
			reject(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		accept(ctx, "更新分配成功", nil)
	}
}

// RemoveAssignment godoc
// @Summary 	删除分配
// @Description 删除指定分配记录，需要当前用户在该章节中拥有 reviewer 角色
//
// @Tags 		assignment
// @Security 	ApiKeyAuth
// @Produce 	json
// @Param 		assignment_id path string true "分配 ID"
//
// @Success 	200
//
// @Router 		/assignments/{assignment_id} [delete]
func RemoveAssignment(appState *state.AppState) iris.Handler {
	assignmentApp := appState.AssignmentApp

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		assignmentID := ctx.Params().Get("assignment_id")
		if assignmentID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 assignment_id 路径参数")
			return
		}

		if err := assignmentApp.Remove(
			buildReqCx(ctx),
			currentUserID,
			assignmentID,
		); err != nil {
			reject(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		accept(ctx, "删除分配成功", nil)
	}
}
