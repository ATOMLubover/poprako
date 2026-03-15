package query_option

import intf "labelplus-next-web-be/internal/domain/repository"

type assignmentQuery struct{}

func AssignmentQuery() assignmentQuery {
	return assignmentQuery{}
}

func (assignmentQuery) FilterByChapterID(chapterID string) intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.Where("assignment_table.chapter_id = ?", chapterID)
	}
}

func (assignmentQuery) FilterByUserID(userID string) intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.Where("assignment_table.user_id = ?", userID)
	}
}

// IncludeUserInfo 聚合查询 assignment 所属用户字段。
func (assignmentQuery) IncludeUserInfo() intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.
			Joins("JOIN user_table ON user_table.id = assignment_table.user_id AND user_table.deleted_at IS NULL").
			Select(
				"assignment_table.*, " +
					"user_table.name               AS user_name, " +
					"user_table.qq                 AS user_qq, " +
					"user_table.avatar_oss_key     AS user_avatar_oss_key, " +
					"user_table.is_avatar_uploaded AS user_is_avatar_uploaded, " +
					"user_table.is_super_admin     AS user_is_super_admin, " +
					"user_table.created_at         AS user_created_at, " +
					"user_table.updated_at         AS user_updated_at",
			)
	}
}

// IncludeChapterInfo 聚合查询 assignment 所属章节及漫画字段（chapter 隐含 chapter.comic）。
func (assignmentQuery) IncludeChapterInfo() intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.
			Joins("LEFT JOIN chapter_table ON chapter_table.id = assignment_table.chapter_id AND chapter_table.deleted_at IS NULL").
			Joins("LEFT JOIN comic_table ON comic_table.id = chapter_table.comic_id AND comic_table.deleted_at IS NULL").
			Joins("LEFT JOIN workset_table ON workset_table.id = comic_table.workset_id").
			Select(
				"assignment_table.*, " +
					"chapter_table.comic_id      AS chapter_comic_id, " +
					"chapter_table.index         AS chapter_index, " +
					"chapter_table.subtitle      AS chapter_subtitle, " +
					"chapter_table.page_count    AS chapter_page_count, " +
					"chapter_table.cover_url     AS chapter_cover_url, " +
					"chapter_table.created_at    AS chapter_created_at, " +
					"chapter_table.updated_at    AS chapter_updated_at, " +
					"comic_table.workset_id      AS comic_workset_id, " +
					"workset_table.team_id       AS comic_team_id, " +
					"comic_table.index           AS comic_index, " +
					"comic_table.title           AS comic_title, " +
					"comic_table.author          AS comic_author, " +
					"comic_table.description     AS comic_description, " +
					"comic_table.cover_url       AS comic_cover_url, " +
					"comic_table.chapter_count   AS comic_chapter_count, " +
					"comic_table.creator_id      AS comic_creator_id, " +
					"comic_table.last_active_at  AS comic_last_active_at, " +
					"comic_table.created_at      AS comic_created_at, " +
					"comic_table.updated_at      AS comic_updated_at",
			)
	}
}
