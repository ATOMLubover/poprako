package query_option

import intf "labelplus-next-web-be/internal/domain/repository"

type chapterQuery struct{}

func ChapterQuery() chapterQuery {
	return chapterQuery{}
}

// FilterByComicID 精确匹配章节所属漫画 ID
func (chapterQuery) FilterByComicID(comicID string) intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.Where("chapter_table.comic_id = ?", comicID)
	}
}

// OrderByIndexAsc 按章节顺序降序排列
func (chapterQuery) OrderByIndexDesc() intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.Order("chapter_table.index DESC")
	}
}

// IncludeCreatorInfo 聚合查询章节创建者的用户字段。
func (chapterQuery) IncludeCreatorInfo() intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.
			Joins("JOIN user_table AS creator_table ON creator_table.id = chapter_table.creator_id AND creator_table.deleted_at IS NULL").
			Select(
				"chapter_table.*, " +
					"creator_table.name AS creator_name, creator_table.qq AS creator_qq, " +
					"creator_table.avatar_oss_key AS creator_avatar_oss_key, " +
					"creator_table.is_avatar_uploaded AS creator_is_avatar_uploaded, " +
					"creator_table.is_super_admin AS creator_is_super_admin, " +
					"creator_table.created_at AS creator_created_at, " +
					"creator_table.updated_at AS creator_updated_at",
			)
	}
}
