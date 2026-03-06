package query_option

import (
	intf "labelplus-next-web-be/internal/domain/repository"
)

type comicQuery struct{}

func ComicQuery() comicQuery {
	return comicQuery{}
}

// FilterByCreatorID 精确匹配漫画的创建者 ID
func (comicQuery) FilterByCreatorID(creatorID string) intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.Where("creator_id = ?", creatorID)
	}
}

// FilterByTeamID 精确匹配漫画所属汉化组 ID
func (comicQuery) FilterByTeamID(teamID string) intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.Where("team_id = ?", teamID)
	}
}
