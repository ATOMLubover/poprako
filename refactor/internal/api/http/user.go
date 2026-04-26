package http

import (
	"poprako-s/internal/api/http/res"
	"poprako-s/internal/api/state"
	"poprako-s/internal/app/val"

	"github.com/kataras/iris/v12"
)

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
			res.Reject(cx, re.Code(), "获取用户信息失败")
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

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
			res.Reject(cx, re.Code(), "获取用户信息失败")
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

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
			res.Reject(cx, re.Code(), "预留头像失败")
			return
		}

		res.Accept(cx, iris.StatusOK, nil)
	}
}

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
			res.Reject(cx, re.Code(), "标记头像上传失败")
			return
		}

		res.Accept(cx, iris.StatusOK, nil)
	}
}
