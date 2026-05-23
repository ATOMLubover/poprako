package http

import (
	"poprako-s/internal/api/http/res"
	"poprako-s/internal/api/state"
	"poprako-s/internal/app/val"

	"github.com/kataras/iris/v12"
)

// `ListComments` godoc
// @Summary List Team Comments
// @Description List comments under one team
// @Description The caller must be member of the target team
// @Description Auth: `authorization` cookie is preferred over `Authorization` header when both are present
// @Tags comment
// @Security ApiKeyAuth
// @Produce json
// @Param team_id query string true "team id"
// @Param offset query int false "pagination offset"
// @Param limit query int false "pagination limit"
// @Success 200 {object} res.HttpRes[[]val.CommentVal]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /api/v1/comments [get]
func ListComments(st *state.AppState) iris.Handler {
	commentApp := st.CommentApp

	return func(cx iris.Context) {
		currUid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		var args val.ListCommentArgs
		if err := cx.ReadQuery(&args); err != nil {
			res.Reject(cx, iris.StatusBadRequest, "请求参数解析失败")
			return
		}

		if args.TeamId == "" {
			res.Reject(cx, iris.StatusBadRequest, "缺少 team_id 参数")
			return
		}

		re := commentApp.List(newReqCx(cx), currUid, &args)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

// `CreateComment` godoc
// @Summary Create Team Comment
// @Description Create one team board comment
// @Description The caller must be team member
// @Description Auth: `authorization` cookie is preferred over `Authorization` header when both are present
// @Tags comment
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param body body val.CreateCommentArgs true "create comment args"
// @Success 201 {object} res.HttpRes[val.CommentCreatedRes]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /api/v1/comments [post]
func CreateComment(st *state.AppState) iris.Handler {
	commentApp := st.CommentApp

	return func(cx iris.Context) {
		currUid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		var args val.CreateCommentArgs
		if err := cx.ReadJSON(&args); err != nil {
			res.Reject(cx, iris.StatusBadRequest, "请求参数解析失败")
			return
		}

		re := commentApp.Create(newReqCx(cx), currUid, &args)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusCreated, re.Data())
	}
}
