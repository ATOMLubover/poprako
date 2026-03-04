package state

import (
	"labelplus-next-web-be/internal/application"
	"labelplus-next-web-be/internal/config"
)

type AppState struct {
    Config *config.AppConfig

	UserApplication       application.UserApplication
	InvitationApplication application.InvitationApplication
	TeamApplication       application.TeamApplication
	MemberApplication     application.MemberApplication
}
