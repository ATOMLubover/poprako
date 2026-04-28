package app_iface

import (
	"context"

	"poprako-s/internal/app/res"
	"poprako-s/internal/app/val"
)

// `TeamApp` defines application use-cases for `Team`.
type TeamApp interface {
	// `GetInfo` returns team details by team id.
	GetInfo(cx context.Context, id string) res.AppRes[val.TeamVal]

	// NOTE: `ListByUser` can also be used as ListMyTeams.
	ListByUser(cx context.Context, userId string) res.AppRes[[]val.TeamVal]

	Update(cx context.Context, args *val.TeamUpdArgs) res.AppRes[res.None]

	ResvAvatar(cx context.Context, args *val.ResvTeamAvatarArgs) res.AppRes[val.ResvTeamAvatarRes]

	MarkAvatarUploaded(cx context.Context, teamId string) res.AppRes[res.None]
}
