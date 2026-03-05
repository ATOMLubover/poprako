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
		return executor.Where("user_id = ?", userID)
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
		return executor.Where("team_id = ?", teamID)
	}
}
