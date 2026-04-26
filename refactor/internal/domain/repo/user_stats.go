package repo_iface

import "poprako-s/internal/domain/model/aggr"

type UserStatsRepo interface {
	// `EnsureGet` get or create a user stats record.
	EnsureGet(id string) (*aggr.UserStats, RepoErr)

	// `Patch` update optimistically user stats.
	Patch(pat *aggr.UserStatsPatch) RepoErr
}
