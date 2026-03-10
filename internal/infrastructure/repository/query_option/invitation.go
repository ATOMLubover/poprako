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
		return executor.Where("invitation_table.target_team_id = ?", teamID)
	}
}

func (invitationQuery) FilterByInviteeQQ(inviteeQQ string) intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.Where("invitation_table.invitee_qq = ?", inviteeQQ)
	}
}

func (invitationQuery) FilterByCode(code string) intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.Where("invitation_table.invitation_code = ?", code)
	}
}

func (invitationQuery) FilterPending(pending bool) intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.Where("invitation_table.pending = ?", pending)
	}
}
