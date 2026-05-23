package http

import (
	"poprako-s/internal/api/http/res"
	"poprako-s/internal/api/state"
	"poprako-s/internal/app/val"

	"github.com/kataras/iris/v12"
)

// `ListAnnouncements` godoc
// @Summary List Team Announcements
// @Description List announcements under one team
// @Description The caller must be member of the target team
// @Description Auth: `authorization` cookie is preferred over `Authorization` header when both are present
// @Tags announcement
// @Security ApiKeyAuth
// @Produce json
// @Param team_id query string true "team id"
// @Param offset query int false "pagination offset"
// @Param limit query int false "pagination limit"
// @Success 200 {object} res.HttpRes[[]val.AnnouncementVal]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /api/v1/announcements [get]
func ListAnnouncements(st *state.AppState) iris.Handler {
	announcementApp := st.AnnouncementApp

	return func(cx iris.Context) {
		currUid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		var args val.ListAnnouncementArgs
		if err := cx.ReadQuery(&args); err != nil {
			res.Reject(cx, iris.StatusBadRequest, "请求参数解析失败")
			return
		}

		if args.TeamId == "" {
			res.Reject(cx, iris.StatusBadRequest, "缺少 team_id 参数")
			return
		}

		re := announcementApp.List(newReqCx(cx), currUid, &args)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}

// `CreateAnnouncement` godoc
// @Summary Create Team Announcement
// @Description Create one announcement under one team
// @Description The caller must be team admin
// @Description Auth: `authorization` cookie is preferred over `Authorization` header when both are present
// @Tags announcement
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param body body val.CreateAnnouncementArgs true "create announcement args"
// @Success 201 {object} res.HttpRes[val.AnnouncementCreatedRes]
// @Failure 400 {object} res.HttpRes[any]
// @Failure 401 {object} res.HttpRes[any]
// @Failure 500 {object} res.HttpRes[any]
// @Router /api/v1/announcements [post]
func CreateAnnouncement(st *state.AppState) iris.Handler {
	announcementApp := st.AnnouncementApp

	return func(cx iris.Context) {
		currUid, ok := takeCurrUid(cx)
		if !ok {
			res.Reject(cx, iris.StatusUnauthorized, "未授权的访问")
			return
		}

		var args val.CreateAnnouncementArgs
		if err := cx.ReadJSON(&args); err != nil {
			res.Reject(cx, iris.StatusBadRequest, "请求参数解析失败")
			return
		}

		re := announcementApp.Create(newReqCx(cx), currUid, &args)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), re.Msg())
			return
		}

		res.Accept(cx, iris.StatusCreated, re.Data())
	}
}
