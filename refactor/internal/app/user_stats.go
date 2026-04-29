package app_iface

import (
	"context"

	"poprako-s/internal/app/res"
	"poprako-s/internal/app/val"
)

type UserStatsApp interface {
	// `GetStats` is used to get user stats by user id.
	// NOTE: it can be used as GetMyStats.
	GetStats(cx context.Context, currUid string) app_res.AppRes[val.UserStats]
}
