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

// IncludeCreatorInfo 聚合查询 comic 创建者字段。
func (comicQuery) IncludeCreatorInfo() intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.
			Joins("LEFT JOIN user_table AS creator_table ON creator_table.id = comic_table.creator_id").
			Select(
				"comic_table.*, workset_table.team_id AS team_id, " +
					"creator_table.name AS creator_name, " +
					"creator_table.qq AS creator_qq, " +
					"creator_table.avatar_oss_key AS creator_avatar_oss_key, " +
					"creator_table.is_avatar_uploaded AS creator_is_avatar_uploaded, " +
					"creator_table.is_super_admin AS creator_is_super_admin, " +
					"creator_table.created_at AS creator_created_at, " +
					"creator_table.updated_at AS creator_updated_at",
			)
	}
}

// IncludeWorksetAndCreatorInfo 聚合查询 comic 所属工作集与创建者字段。
func (comicQuery) IncludeWorksetAndCreatorInfo() intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.
			Joins("LEFT JOIN user_table AS creator_table ON creator_table.id = comic_table.creator_id").
			Select(
				"comic_table.*, workset_table.team_id AS team_id, " +
					"workset_table.name AS workset_name, " +
					"workset_table.index AS workset_index, " +
					"workset_table.comic_count AS workset_comic_count, " +
					"workset_table.created_at AS workset_created_at, " +
					"workset_table.updated_at AS workset_updated_at, " +
					"creator_table.name AS creator_name, " +
					"creator_table.qq AS creator_qq, " +
					"creator_table.avatar_oss_key AS creator_avatar_oss_key, " +
					"creator_table.is_avatar_uploaded AS creator_is_avatar_uploaded, " +
					"creator_table.is_super_admin AS creator_is_super_admin, " +
					"creator_table.created_at AS creator_created_at, " +
					"creator_table.updated_at AS creator_updated_at",
			)
	}
}
