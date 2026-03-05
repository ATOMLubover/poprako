package query_option

import (
	intf "labelplus-next-web-be/internal/domain/repository"
)

type invitationQuery struct{}

func InvitationQuery() invitationQuery {
	return invitationQuery{}
}

func (invitationQuery) FilterByTeamID(teamID string) intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.Where("target_team_id = ?", teamID)
	}
}
