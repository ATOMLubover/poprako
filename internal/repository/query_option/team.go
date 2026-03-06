package query_option

import intf "labelplus-next-web-be/internal/domain/repository"

type teamQuery struct{}

func TeamQuery() teamQuery {
	return teamQuery{}
}

func (teamQuery) FilterByTeamID(teamID string) func(executor intf.Executor) intf.Executor {
	return func(executor intf.Executor) intf.Executor {
		return executor.Where("team_id = ?", teamID)
	}
}
