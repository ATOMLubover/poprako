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

// IncludeInvitorInfo 聚合查询邀请发起人的用户字段。
func (invitationQuery) IncludeInvitorInfo() intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.
			Joins("LEFT JOIN user_table AS invitor_table ON invitor_table.id = invitation_table.invitor_id").
			Select(
				"invitation_table.*, " +
					"invitor_table.name AS invitor_name, " +
					"invitor_table.qq AS invitor_qq, " +
					"invitor_table.avatar_oss_key AS invitor_avatar_oss_key, " +
					"invitor_table.is_avatar_uploaded AS invitor_is_avatar_uploaded, " +
					"invitor_table.is_super_admin AS invitor_is_super_admin, " +
					"invitor_table.created_at AS invitor_created_at, " +
					"invitor_table.updated_at AS invitor_updated_at",
			)
	}
}
