package http

import (
	"poprako-s/internal/api/http/res"
	"poprako-s/internal/api/state"
	"poprako-s/internal/app/val"

	"github.com/kataras/iris/v12"
)

// `GetUserInfo` godoc
// @Summary Get User Info
//
//	Get user info by user id and return a `res.HttpRes` wrapper with `val.UserVal`.
//	Auth: `authorization` cookie is preferred over `Authorization` header when both are present
//
// @Tags user
// @Security ApiKeyAuth
// @Produce json
// @Param user_id path string true "user id"
// @Success 200 {object} res.HttpRes
// @Failure 400 {object} res.HttpRes
// @Router /user/{user_id} [get]
func GetUserInfo(st *state.AppState) iris.Handler {
	userApp := st.UserApp

	return func(cx iris.Context) {
		uid := cx.Params().Get("user_id")
		if uid == "" {
			res.Reject(cx, iris.StatusBadRequest, "缺少 user_id 参数")
			return
		}

		re := userApp.GetInfo(newReqCx(cx), uid)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), "获取用户信息失败")
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

// `GetMyUserInfo` godoc
// @Summary Get My User Info
//
//	Get current authorized user info and return a `res.HttpRes` wrapper with `val.UserVal`.
//	Auth: `authorization` cookie is preferred over `Authorization` header when both are present
//
// @Tags user
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} res.HttpRes
// @Failure 401 {object} res.HttpRes
// @Router /user/me [get]
func GetMyUserInfo(st *state.AppState) iris.Handler {
	userApp := st.UserApp

	return func(cx iris.Context) {
		uid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		re := userApp.GetInfo(newReqCx(cx), uid)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), "获取用户信息失败")
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

// `ResvUserAvatar` godoc
// @Summary Reserve User Avatar Upload
//
//	Reserve a signed upload url for user avatar and return a `res.HttpRes` wrapper with `val.ResvUserAvatarRes`.
//	Auth: `authorization` cookie is preferred over `Authorization` header when both are present
//
// @Tags user
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param body body val.ResvUserAvatarArgs true "reserve avatar args"
// @Success 200 {object} res.HttpRes
// @Failure 400 {object} res.HttpRes
// @Router /user/avatar [post]
func ResvUserAvatar(st *state.AppState) iris.Handler {
	userApp := st.UserApp

	return func(cx iris.Context) {
		var args val.ResvUserAvatarArgs

		if err := cx.ReadJSON(&args); err != nil {
			res.Reject(cx, iris.StatusBadRequest, "请求参数错误")
			return
		}

		re := userApp.ResvAvatar(newReqCx(cx), &args)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), "预留头像失败")
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

// `MarkUserAvatarUploaded` godoc
// @Summary Confirm User Avatar Uploaded
//
//	Confirm avatar uploaded after client upload completed and return no JSON body.
//	Auth: `authorization` cookie is preferred over `Authorization` header when both are present
//
// @Tags user
// @Security ApiKeyAuth
// @Produce json
// @Success 200
// @Failure 401 {object} res.HttpRes
// @Router /user/avatar/confirm [post]
func MarkUserAvatarUploaded(st *state.AppState) iris.Handler {
	userApp := st.UserApp

	return func(cx iris.Context) {
		uid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		re := userApp.MarkAvatarUploaded(newReqCx(cx), uid)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), "标记头像上传失败")
			return
		}

		res.Accept(cx, iris.StatusOK, nil)
	}
}
