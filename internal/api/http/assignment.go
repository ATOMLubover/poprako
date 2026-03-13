package http

import (
	"labelplus-next-web-be/internal/state"
	"labelplus-next-web-be/internal/value"

	"github.com/kataras/iris/v12"
)

// ListChapterAssignments godoc
// @Summary 	获取章节分配列表
// @Description 获取指定章节的所有分配记录（含被分配用户信息），仅汉化组成员可访问
//
// @Tags 		assignment
// @Security 	ApiKeyAuth
// @Produce 	json
// @Param 		chapter_id query string true "章节 ID"
//
// @Success 	200 {object} []value.AssignmentWithUserInfo
//
// @Router 		/assignments [get]
func ListChapterAssignments(appState *state.AppState) iris.Handler {
	assignmentApplication := appState.AssignmentApplication

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrentUserID(ctx)
		if !ok {
			return
		}

		var args value.ListChapterAssignmentArgs
		if err := ctx.ReadQuery(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "查询参数格式错误: "+err.Error())
			return
		}

		result, err := assignmentApplication.ListChapterAssignments(
			buildTraceScope(ctx),
			currentUserID,
			args,
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
// @Description 获取当前登录用户的所有分配记录（含章节和漫画信息），支持分页
//
// @Tags 		assignment
// @Security 	ApiKeyAuth
// @Produce 	json
// @Param 		offset query int true "偏移量"
// @Param 		limit query int true "每页数量"
//
// @Success 	200 {object} []value.AssignmentWithChapterInfo
//
// @Router 		/assignments/mine [get]
func ListMyAssignments(appState *state.AppState) iris.Handler {
	assignmentApplication := appState.AssignmentApplication

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrentUserID(ctx)
		if !ok {
			return
		}

		var args value.ListAssignmentArgs
		if err := ctx.ReadQuery(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "查询参数格式错误: "+err.Error())
			return
		}

		result, err := assignmentApplication.ListMyAssignments(
			buildTraceScope(ctx),
			currentUserID,
			args,
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
// @Param 		body body value.CreateChapterAssignmentArgs true "创建分配参数"
//
// @Success 	201 {object} value.CreateChapterAssignmentResult
//
// @Router 		/assignments [post]
func CreateChapterAssignment(appState *state.AppState) iris.Handler {
	assignmentApplication := appState.AssignmentApplication

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrentUserID(ctx)
		if !ok {
			return
		}

		var args value.CreateChapterAssignmentArgs
		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误: "+err.Error())
			return
		}

		result, err := assignmentApplication.CreateChapterAssignment(
			buildTraceScope(ctx),
			currentUserID,
			args,
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
// @Param 		body body value.UpdateAssignmentArgs true "更新分配参数"
//
// @Success 	200
//
// @Router 		/assignments/{assignment_id} [put]
func UpdateAssignment(appState *state.AppState) iris.Handler {
	assignmentApplication := appState.AssignmentApplication

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrentUserID(ctx)
		if !ok {
			return
		}

		assignmentID := ctx.Params().Get("assignment_id")
		if assignmentID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 assignment_id 路径参数")
			return
		}

		var args value.UpdateAssignmentArgs
		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误: "+err.Error())
			return
		}
		args.ID = assignmentID

		if err := assignmentApplication.UpdateAssignment(
			buildTraceScope(ctx),
			currentUserID,
			args,
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
	assignmentApplication := appState.AssignmentApplication

	return func(ctx iris.Context) {
		currentUserID, ok := extractCurrentUserID(ctx)
		if !ok {
			return
		}

		assignmentID := ctx.Params().Get("assignment_id")
		if assignmentID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 assignment_id 路径参数")
			return
		}

		if err := assignmentApplication.RemoveAssignment(
			buildTraceScope(ctx),
			currentUserID,
			assignmentID,
		); err != nil {
			reject(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		accept(ctx, "删除分配成功", nil)
	}
}
