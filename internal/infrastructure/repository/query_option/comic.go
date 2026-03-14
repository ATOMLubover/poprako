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
		return executor.Where("comic_table.creator_id = ?", creatorID)
	}
}

// FilterByWorksetID 精确匹配漫画所属工作集 ID
func (comicQuery) FilterByWorksetID(worksetID string) intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.Where("comic_table.workset_id = ?", worksetID)
	}
}

// OrderByLastActiveAtDesc 根据漫画的最后活跃时间降序排序
func (comicQuery) OrderByLastActiveAtDesc() intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.Order("comic_table.last_active_at DESC")
	}
}

// IncludeWorksetInfo 聚合查询 comic 所属工作集字段。
func (comicQuery) IncludeWorksetInfo() intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.Select(
			"comic_table.*, workset_table.team_id AS team_id, " +
				"workset_table.name AS workset_name, " +
				"workset_table.index AS workset_index, " +
				"workset_table.comic_count AS workset_comic_count, " +
				"workset_table.created_at AS workset_created_at, " +
				"workset_table.updated_at AS workset_updated_at",
		)
	}
}
