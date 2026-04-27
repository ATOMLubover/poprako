package http

import (
	"poprako-s/internal/api/http/res"
	"poprako-s/internal/api/state"
	"poprako-s/internal/app/val"

	"github.com/kataras/iris/v12"
)

// `LoginUser` godoc
// @Summary User Login
//
//	Login by qid and password and return a `res.HttpRes` wrapper with `val.UserLoginRes`
//
// @Tags auth
// @Accept json
// @Produce json
// @Param body body val.UserLoginArgs true "login args"
// @Success 200 {object} res.HttpRes
// @Failure 400 {object} res.HttpRes
// @Router /auth/login [post]
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
			res.Reject(cx, int(re.Code()), "登录失败")
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

// `RegUser` godoc
// @Summary User Registration
//
//	Register by invitation code and return a `res.HttpRes` wrapper with `val.UserRegRes`
//
// @Tags auth
// @Accept json
// @Produce json
// @Param body body val.UserRegArgs true "registration args"
// @Success 200 {object} res.HttpRes
// @Failure 400 {object} res.HttpRes
// @Router /auth/reg [post]
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
			res.Reject(cx, int(re.Code()), "注册失败")
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}
