package http

import (
	"poprako-s/internal/api/http/res"
	"poprako-s/internal/api/state"
	"poprako-s/internal/app/val"

	"github.com/kataras/iris/v12"
)

// `ListSysMail` godoc
// @Summary List Unread System Mails
// @Description List unread system mails for current authorized user with pagination
// @Description Auth: `authorization` cookie is preferred over `Authorization` header when both are present
// @Tags sys-mail
// @Security ApiKeyAuth
// @Produce json
// @Param offset query int false "pagination offset"
// @Param limit query int false "pagination limit"
// @Success 200 {object} res.HttpRes[[]val.SysMailVal]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /api/v1/sys-mails [get]
func ListSysMail(st *state.AppState) iris.Handler {
	sysMailApp := st.SysMailApp

	return func(cx iris.Context) {
		currUid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		offset, err := cx.URLParamInt("offset")
		if err != nil {
			res.Reject(cx, iris.StatusBadRequest, "offset 参数格式错误")
			return
		}

		limit, err := cx.URLParamInt("limit")
		if err != nil {
			res.Reject(cx, iris.StatusBadRequest, "limit 参数格式错误")
			return
		}

		args := &val.ListSysMailArgs{
			Offset: offset,
			Limit:  limit,
		}

		re := sysMailApp.List(newReqCx(cx), currUid, args)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

// `MarkSysMailRead` godoc
// @Summary Mark System Mail Read
// @Description Mark one system mail as read for current authorized user
// @Description Auth: `authorization` cookie is preferred over `Authorization` header when both are present
// @Tags sys-mail
// @Security ApiKeyAuth
// @Produce json
// @Param sys_mail_id path string true "system mail id"
// @Success 200 {object} res.HttpRes[any]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /api/v1/sys-mails/{sys_mail_id}/read [post]
func MarkSysMailRead(st *state.AppState) iris.Handler {
	sysMailApp := st.SysMailApp

	return func(cx iris.Context) {
		currUid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		sysMailId := cx.Params().Get("sys_mail_id")
		if sysMailId == "" {
			res.Reject(cx, iris.StatusBadRequest, "缺少 sys_mail_id 参数")
			return
		}

		re := sysMailApp.MarkRead(newReqCx(cx), currUid, sysMailId)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}
