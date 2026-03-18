package query_option

import (
	intf "labelplus-next-web-be/internal/domain/repository"
)

type memberQuery struct{}

func MemberQuery() memberQuery {
	return memberQuery{}
}

// FilterByUserID 精确匹配成员的用户 ID
func (memberQuery) FilterByUserID(userID string) intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.Where("member_table.user_id = ?", userID)
	}
}

// FilterByUserQQ 精确匹配成员的 QQ
func (memberQuery) FilterOnUserQQ(qq string) intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.Where("user_table.qq = ?", qq)
	}
}

// FilterByTeamID 精准匹配成员所在的汉化组 ID
func (memberQuery) FilterByTeamID(teamID string) intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.Where("member_table.team_id = ?", teamID)
	}
}

// JoinUser 为查询添加 user_table 的 LEFT JOIN（与 user 字段一起使用）
func (memberQuery) JoinUser() intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.Joins("LEFT JOIN user_table ON user_table.id = member_table.user_id AND user_table.deleted_at IS NULL")
	}
}

// IncludeUserInfo 聚合查询 member 所属用户字段。
func (memberQuery) IncludeUserInfo() intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.
			Joins("JOIN user_table ON user_table.id = member_table.user_id AND user_table.deleted_at IS NULL").
			Select(
				"member_table.id, member_table.user_id, member_table.team_id, " +
					"member_table.assigned_raw_provider_at, member_table.assigned_translator_at, " +
					"member_table.assigned_proofreader_at, member_table.assigned_typesetter_at, " +
					"member_table.assigned_reviewer_at, member_table.assigned_publisher_at, " +
					"member_table.assigned_admin_at, member_table.created_at, member_table.updated_at, " +
					"user_table.name AS user_name, user_table.qq AS user_qq, " +
					"user_table.avatar_oss_key AS user_avatar_oss_key, " +
					"user_table.is_avatar_uploaded AS user_is_avatar_uploaded, " +
					"user_table.is_super_admin AS user_is_super_admin, " +
					"user_table.created_at AS user_created_at, " +
					"user_table.updated_at AS user_updated_at",
			)
	}
}

// IncludeTeamInfo 聚合查询 member 所属汉化组字段。
func (memberQuery) IncludeTeamInfo() intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.
			Joins("JOIN team_table ON team_table.id = member_table.team_id AND team_table.deleted_at IS NULL").
			Select(
				"member_table.id, member_table.user_id, member_table.team_id, " +
					"member_table.assigned_raw_provider_at, member_table.assigned_translator_at, " +
					"member_table.assigned_proofreader_at, member_table.assigned_typesetter_at, " +
					"member_table.assigned_reviewer_at, member_table.assigned_publisher_at, " +
					"member_table.assigned_admin_at, member_table.created_at, member_table.updated_at, " +
					"team_table.name AS team_name, team_table.description AS team_description, " +
					"team_table.avatar_oss_key AS team_avatar_oss_key, " +
					"team_table.is_avatar_uploaded AS team_is_avatar_uploaded, " +
					"team_table.created_at AS team_created_at, " +
					"team_table.updated_at AS team_updated_at",
			)
	}
}
