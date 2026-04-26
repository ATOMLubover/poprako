package state

import (
	app_iface "poprako-s/internal/app"
	"poprako-s/internal/cfg"
)

type AppState struct {
	Cfg *cfg.AppCfg

	UserApp app_iface.UserApp
}
