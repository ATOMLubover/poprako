package query_option

import intf "labelplus-next-web-be/internal/domain/repository"

type pageQuery struct{}

func PageQuery() pageQuery {
	return pageQuery{}
}

func (pageQuery) FilterByChapterID(chapterID string) intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.Where("page_table.chapter_id = ?", chapterID)
	}
}

func (pageQuery) OrderByIndexAsc() intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.Order("page_table.index ASC")
	}
}

// IncludeCreatorInfo 聚合查询页面创建者字段。
func (pageQuery) IncludeCreatorInfo() intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.
			Joins("LEFT JOIN user_table AS creator_table ON creator_table.id = page_table.creator_id").
			Select(
				"page_table.*, " +
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
