package query_option

import (
	"strings"

	intf "labelplus-next-web-be/internal/domain/repository"
)

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

// IncludeRelationInfo 聚合查询 assignment 的关联字段。
// 支持 user, chapter, chapter.comic, chapter.creator 的任意组合。
func (assignmentQuery) IncludeRelationInfo(
	needUser bool,
	needChapter bool,
	needChapterComic bool,
	needChapterCreator bool,
) intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		selectFields := []string{"assignment_table.*"}

		if needUser {
			executor = executor.Joins("JOIN user_table AS assignment_user_table ON assignment_user_table.id = assignment_table.user_id AND assignment_user_table.deleted_at IS NULL")
			selectFields = append(
				selectFields,
				"assignment_user_table.name               AS user_name",
				"assignment_user_table.qq                 AS user_qq",
				"assignment_user_table.avatar_oss_key     AS user_avatar_oss_key",
				"assignment_user_table.is_avatar_uploaded AS user_is_avatar_uploaded",
				"assignment_user_table.is_super_admin     AS user_is_super_admin",
				"assignment_user_table.created_at         AS user_created_at",
				"assignment_user_table.updated_at         AS user_updated_at",
			)
		}

		if needChapter {
			executor = executor.Joins("LEFT JOIN chapter_table ON chapter_table.id = assignment_table.chapter_id AND chapter_table.deleted_at IS NULL")
			selectFields = append(
				selectFields,
				"chapter_table.comic_id      AS chapter_comic_id",
				"chapter_table.index         AS chapter_index",
				"chapter_table.subtitle      AS chapter_subtitle",
				"chapter_table.page_count    AS chapter_page_count",
				"chapter_table.creator_id    AS chapter_creator_id",
				"chapter_table.created_at    AS chapter_created_at",
				"chapter_table.updated_at    AS chapter_updated_at",
			)
		}

		if needChapterComic {
			executor = executor.
				Joins("LEFT JOIN comic_table ON comic_table.id = chapter_table.comic_id AND comic_table.deleted_at IS NULL").
				Joins("LEFT JOIN workset_table ON workset_table.id = comic_table.workset_id")
			selectFields = append(
				selectFields,
				"comic_table.workset_id      AS comic_workset_id",
				"comic_table.index           AS comic_index",
				"comic_table.title           AS comic_title",
				"comic_table.author          AS comic_author",
				"comic_table.description     AS comic_description",
				"comic_table.chapter_count   AS comic_chapter_count",
				"comic_table.creator_id      AS comic_creator_id",
				"comic_table.last_active_at  AS comic_last_active_at",
				"comic_table.created_at      AS comic_created_at",
				"comic_table.updated_at      AS comic_updated_at",
			)
		}

		if needChapterCreator {
			executor = executor.Joins("LEFT JOIN user_table AS chapter_creator_table ON chapter_creator_table.id = chapter_table.creator_id AND chapter_creator_table.deleted_at IS NULL")
			selectFields = append(
				selectFields,
				"chapter_creator_table.name               AS chapter_creator_name",
				"chapter_creator_table.qq                 AS chapter_creator_qq",
				"chapter_creator_table.avatar_oss_key     AS chapter_creator_avatar_oss_key",
				"chapter_creator_table.is_avatar_uploaded AS chapter_creator_is_avatar_uploaded",
				"chapter_creator_table.is_super_admin     AS chapter_creator_is_super_admin",
				"chapter_creator_table.created_at         AS chapter_creator_created_at",
				"chapter_creator_table.updated_at         AS chapter_creator_updated_at",
			)
		}

		return executor.Select(strings.Join(selectFields, ", "))
	}
}

// IncludeUserInfo 聚合查询 assignment 所属用户字段。
func (assignmentQuery) IncludeUserInfo() intf.QueryOption {
	return AssignmentQuery().IncludeRelationInfo(true, false, false, false)
}

// IncludeChapterInfo 聚合查询 assignment 所属章节及漫画字段（兼容旧语义）。
func (assignmentQuery) IncludeChapterInfo() intf.QueryOption {
	return AssignmentQuery().IncludeRelationInfo(false, true, true, false)
}
