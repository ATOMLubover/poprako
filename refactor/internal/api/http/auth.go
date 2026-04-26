package http

import (
	"poprako-s/internal/api/http/res"
	"poprako-s/internal/api/state"
	"poprako-s/internal/app/val"

	"github.com/kataras/iris/v12"
)

func LoginUser(st *state.AppState) iris.Handler {
	userApp := st.UserApp

	return func(cx iris.Context) {
		var args val.UserLoginArgs

		if err := cx.ReadJSON(&args); err != nil {
			res.Reject(cx, iris.StatusBadRequest, "请求参数错误")
			return
		}

		re := userApp.Login(newReqCx(cx), &args)
		if re.IsReject() {
			res.Reject(cx, re.Code(), "登录失败")
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

func RegUser(st *state.AppState) iris.Handler {
	userApp := st.UserApp

	return func(cx iris.Context) {
		var args val.UserRegArgs

		if err := cx.ReadJSON(&args); err != nil {
			res.Reject(cx, iris.StatusBadRequest, "请求参数错误")
			return
		}

		re := userApp.Reg(newReqCx(cx), &args)
		if re.IsReject() {
			res.Reject(cx, re.Code(), "注册失败")
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}
