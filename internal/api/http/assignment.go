package http

import (
	"poprako-s/internal/app/val"
	"poprako-s/internal/state"

	"github.com/kataras/iris/v12"
)

// GetAssignmentByID godoc
// @Summary 	根据 ID 获取分配详情
// @Description 根据分配 ID 获取单个分配记录的详细信息
//
// @Tags 		assignment
// @Security 	ApiKeyAuth
// @Produce 	json
// @Param 		assignment_id path string true "分配 ID"
//
// @Success 	200 {object} val.AssignmentInfo
//
// @Router 		/assignments/{assignment_id} [get]
func GetAssignmentByID(appState *state.AppState) iris.Handler {
	assignmentApp := appState.AssignmentApp

	return func(ctx iris.Context) {
		currUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		assignmentID := ctx.Params().Get("assignment_id")
		if assignmentID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 assignment_id 路径参数")
			return
		}

		result, err := assignmentApp.Get(
			buildReqCx(ctx),
			currUserID,
			assignmentID,
		)
		if err != nil {
			reject(ctx, iris.StatusForbidden, err.Error())
			return
		}

		accept(ctx, "获取分配详情成功", result)
	}
}

// ListChapterAssignments godoc
// @Summary 	获取章节分配列表
// @Description 获取指定章节的所有分配记录，仅汉化组成员可访问
//
// @Tags 		assignment
// @Security 	ApiKeyAuth
// @Produce 	json
// @Param 		chapter_id query string true "章节 ID"
// @Param 		"includes" query []string false "include 关联信息，可选值：user, chapter, chapter.comic, chapter.creator"
// @Param 		offset query int true "偏移量"
// @Param 		limit query int true "每页数量"
//
// @Success 	200 {object} []val.AssignmentInfo
//
// @Router 		/assignments [get]
func ListChapterAssignments(appState *state.AppState) iris.Handler {
	assignmentApp := appState.AssignmentApp

	return func(ctx iris.Context) {
		currUserID, ok := extractCurrUserID(ctx)
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
			currUserID,
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
// @Param 		"includes" query []string false "include 关联信息，可选值：chapter, chapter.comic, chapter.creator"
// @Param 		offset query int true "偏移量"
// @Param 		limit query int true "每页数量"
//
// @Success 	200 {object} []val.AssignmentInfo
//
// @Router 		/assignments/mine [get]
func ListMyAssignments(appState *state.AppState) iris.Handler {
	assignmentApp := appState.AssignmentApp

	return func(ctx iris.Context) {
		currUserID, ok := extractCurrUserID(ctx)
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
			currUserID,
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
		currUserID, ok := extractCurrUserID(ctx)
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
			currUserID,
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
		currUserID, ok := extractCurrUserID(ctx)
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
			currUserID,
			assignmentID,
		); err != nil {
			reject(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		accept(ctx, "删除分配成功", nil)
	}
}

// JoinInvitorChapter godoc
// @Summary 	通过章节邀请加入协作
// @Description 使用章节邀请码加入对应章节协作并授予邀请中的分工
//
// @Tags 		assignment
// @Security 	ApiKeyAuth
// @Accept 		json
// @Produce 	json
// @Param 		body body val.JoinInvitorChapterArgs true "加入章节协作参数"
//
// @Success 	200
//
// @Router 		/assignments/join [post]
func JoinInvitorChapter(appState *state.AppState) iris.Handler {
	assignmentApp := appState.AssignmentApp

	return func(ctx iris.Context) {
		currUserID, ok := extractCurrUserID(ctx)
		if !ok {
			return
		}

		var args val.JoinInvitorChapterArgs

		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误: "+err.Error())
			return
		}

		if err := assignmentApp.JoinInvitorChapter(
			buildReqCx(ctx),
			currUserID,
			&args,
		); err != nil {
			reject(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		accept(ctx, "加入章节协作成功", nil)
	}
}
