package repo_iface

// NOTE: DO NOT expose `UserStatsRepo` now, as relevant logic is not stably implemented.
// type UserStatsRepo interface {
// 	// `EnsureGet` get or create a user stats record.
// 	EnsureGet(id string) (*aggr.UserStats, RepoErr)
//
// 	// `Patch` update optimistically user stats.
// 	Patch(pat *aggr.UserStatsPatch) RepoErr
// }
