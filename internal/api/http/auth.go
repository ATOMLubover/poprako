package http

import (
	"net/http"

	"poprako-s/internal/api/http/middleware"
	"poprako-s/internal/api/http/res"
	"poprako-s/internal/api/state"
	"poprako-s/internal/app/val"

	"github.com/kataras/iris/v12"
)

// `LoginUser` godoc
// @Summary User Login
// @Description Login by qid and password and return a `res.HttpRes` wrapper with `val.UserLoginRes`
// @Tags auth
// @Accept json
// @Produce json
// @Param body body val.UserLoginArgs true "login args"
// @Success 200 {object} res.HttpRes[val.UserLoginRes]
// @Failure 400 {object} res.HttpRes[any]
// @Router /api/v1/auth/login [post]
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

		cx.SetCookie(&http.Cookie{
			Name:     middleware.AuthCookieName,
			Value:    string(re.Data().Token),
			Path:     "/",
			HttpOnly: true,
			MaxAge:   7 * 24 * 3600, // 7 days
		})

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

// `LogoutUser` godoc
// @Summary User Logout
// @Description Invalidate the auth cookie to log out the current user and return 200 OK
// @Tags auth
// @Success 200
// @Router /api/v1/auth/logout [post]
func LogoutUser() iris.Handler {
	return func(cx iris.Context) {
		cx.SetCookie(&http.Cookie{
			Name:   middleware.AuthCookieName,
			Value:  "",
			Path:   "/",
			MaxAge: -1, // delete the cookie
		})

		res.Accept[string](cx, iris.StatusOK, "退出登录成功")
	}
}

// `RegUser` godoc
// @Summary User Registration
// @Description Register by invitation code and return a `res.HttpRes` wrapper with `val.UserRegRes`
// @Tags auth
// @Accept json
// @Produce json
// @Param body body val.UserRegArgs true "registration args"
// @Success 200 {object} res.HttpRes[val.UserRegRes]
// @Failure 400 {object} res.HttpRes[any]
// @Router /api/v1/auth/register [post]
func RegUser(st *state.AppState) iris.Handler {
	userApp := st.UserApp

	return func(cx iris.Context) {
		var args val.UserRegArgs

		if err := cx.ReadJSON(&args); err != nil {
			res.Reject(cx, iris.StatusBadRequest, "请求参数错误")
			return
		}

		re := userApp.Register(newReqCx(cx), &args)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), "注册失败")
			return
		}

		cx.SetCookie(&http.Cookie{
			Name:     middleware.AuthCookieName,
			Value:    string(re.Data().Token),
			Path:     "/",
			HttpOnly: true,
			MaxAge:   7 * 24 * 3600, // 7 days
		})

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}
