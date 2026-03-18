package query_option

import intf "labelplus-next-web-be/internal/domain/repository"

type userStatsQuery struct{}

func UserStatsQuery() userStatsQuery {
	return userStatsQuery{}
}

func (userStatsQuery) FilterByUserID(userID string) intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.Where("user_stats_table.user_id = ?", userID)
	}
}
