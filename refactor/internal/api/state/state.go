package state

import (
	app_iface "poprako-s/internal/app"
	"poprako-s/internal/cfg"
)

// `AppState` holds all application-wide dependencies shared across HTTP handlers.
type AppState struct {
	// `Cfg` is the application configuration.
	Cfg *cfg.AppCfg

	// `UserApp` handles user-related use-cases.
	UserApp app_iface.UserApp

	// `TeamApp` handles team-related use-cases.
	TeamApp app_iface.TeamApp

	// `WorksetApp` handles workset-related use-cases.
	WorksetApp app_iface.WorksetApp
}

// `NewAppState` constructs an `AppState` from its dependencies.
func NewAppState(
	appCfg *cfg.AppCfg,
	userApp app_iface.UserApp,
	teamApp app_iface.TeamApp,
	worksetApp app_iface.WorksetApp,
) *AppState {
	return &AppState{
		Cfg:        appCfg,
		UserApp:    userApp,
		TeamApp:    teamApp,
		WorksetApp: worksetApp,
	}
}
