package query_option

import (
	intf "labelplus-next-web-be/internal/domain/repository"
)

type worksetQuery struct{}

func WorksetQuery() worksetQuery {
	return worksetQuery{}
}

// FilterByTeamID 精确匹配工作集所属汉化组 ID
func (worksetQuery) FilterByTeamID(teamID string) intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.Where("workset_table.team_id = ?", teamID)
	}
}

// OrderByIndexAsc 根据工作集索引升序排序
func (worksetQuery) OrderByIndexAsc() intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.Order("workset_table.index ASC")
	}
}

// IncludeTeamInfo 聚合查询 workset 所属的 team 字段。
func (worksetQuery) IncludeTeamInfo() intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.
			Joins("JOIN team_table ON team_table.id = workset_table.team_id AND team_table.deleted_at IS NULL").
			Select(
				"workset_table.id, workset_table.team_id, workset_table.index, workset_table.name, workset_table.description, workset_table.comic_count, workset_table.created_at, workset_table.updated_at, " +
					"team_table.name AS team_name, team_table.description AS team_description, team_table.avatar_oss_key AS team_avatar_oss_key, team_table.is_avatar_uploaded AS team_is_avatar_uploaded, team_table.created_at AS team_created_at, team_table.updated_at AS team_updated_at",
			)
	}
}
