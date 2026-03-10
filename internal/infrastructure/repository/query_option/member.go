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
