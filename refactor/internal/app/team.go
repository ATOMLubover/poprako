package app_iface

import (
	"context"

	"poprako-s/internal/app/res"
	"poprako-s/internal/app/val"
)

// `TeamApp` defines application use-cases for `Team`.
type TeamApp interface {
	// `Create` creates one team and returns created id.
	Create(cx context.Context, currUid string, args *val.TeamCreArgs) app_res.AppRes[val.TeamCreRes]

	// `GetInfo` returns team details by team id.
	GetInfo(cx context.Context, id string) app_res.AppRes[val.TeamVal]

	// `List` lists all teams.
	List(cx context.Context, currUid string, args *val.ListTeamArgs) app_res.AppRes[[]val.TeamVal]

	// NOTE: `ListByUser` can also be used as ListMyTeams.
	ListByUser(cx context.Context, userId string, args *val.ListTeamArgs) app_res.AppRes[[]val.TeamVal]

	Update(cx context.Context, currUid string, args *val.TeamUpdArgs) app_res.AppRes[app_res.None]

	ResvAvatar(cx context.Context, currUid string, args *val.ResvTeamAvatarArgs) app_res.AppRes[val.ResvTeamAvatarRes]

	MarkAvatarUploaded(cx context.Context, currUid string, teamId string) app_res.AppRes[app_res.None]
}
