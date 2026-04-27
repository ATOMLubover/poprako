package http

import (
	"poprako-s/internal/api/http/res"
	"poprako-s/internal/api/state"

	"github.com/kataras/iris/v12"
)

// `GetTeamInfo` godoc
// @Summary Get Team Info
//
//	Get team info by team id and return a `res.HttpRes` wrapper with `val.TeamVal`.
//	Auth: `authorization` cookie is preferred over `Authorization` header when both are present
//
// @Tags team
// @Security ApiKeyAuth
// @Produce json
// @Param team_id path string true "team id"
// @Success 200 {object} res.HttpRes
// @Failure 400 {object} res.HttpRes
// @Router /team/{team_id} [get]
func GetTeamInfo(st *state.AppState) iris.Handler {
	teamApp := st.TeamApp

	return func(cx iris.Context) {
		teamId := cx.Params().Get("team_id")
		if teamId == "" {
			res.Reject(cx, iris.StatusBadRequest, "缺少 team_id 参数")
			return
		}

		re := teamApp.GetInfo(newReqCx(cx), teamId)
		if re.IsReject() {
			res.Reject(cx, int(re.Code()), "获取团队信息失败")
			return
		}

		res.Accept(cx, iris.StatusOK, re.Data())
	}
}
