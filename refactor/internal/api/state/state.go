package state

import (
	app_iface "poprako-s/internal/app"
	"poprako-s/internal/cfg"
)

type AppState struct {
	Cfg *cfg.AppCfg

	UserApp app_iface.UserApp
	TeamApp app_iface.TeamApp
}

func NewAppState(
	appCfg *cfg.AppCfg,
	userApp app_iface.UserApp,
	teamApp app_iface.TeamApp,
) *AppState {
	return &AppState{
		Cfg:     appCfg,
		UserApp: userApp,
		TeamApp: teamApp,
	}
}
