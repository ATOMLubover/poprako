package query_option

import (
	"labelplus-next-web-be/internal/domain/model"
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

// FuzzyFilterByTitle 对漫画组合标题进行模糊匹配（不区分大小写）。
// 组合标题格式为：【index】[author]title。
func (comicQuery) FuzzyFilterByTitle(title string) intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.Where("comic_table.composed_title ILIKE ?", "%"+title+"%")
	}
}

// FilterByLatestChapterStatus 基于 comic 表中的最新章节副本时间戳字段筛选。
func (comicQuery) FilterByLatestChapterStatus(
	uploadStatus *model.WorkflowStatus,
	translateStatus *model.WorkflowStatus,
	proofreadStatus *model.WorkflowStatus,
	typesetStatus *model.WorkflowStatus,
	reviewStatus *model.WorkflowStatus,
	publishStatus *model.WorkflowStatus,
) intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		if uploadStatus != nil {
			switch *uploadStatus {
			case model.WorkflowPending:
				executor = executor.Where("comic_table.latest_uploaded_at IS NULL")
			case model.WorkflowCompleted:
				executor = executor.Where("comic_table.latest_uploaded_at IS NOT NULL")
			}
		}

		if translateStatus != nil {
			switch *translateStatus {
			case model.WorkflowPending:
				executor = executor.Where("comic_table.latest_transalating_at IS NULL")
			case model.WorkflowInProgress:
				executor = executor.Where("comic_table.latest_transalating_at IS NOT NULL").
					Where("comic_table.latest_translated_at IS NULL")
			case model.WorkflowCompleted:
				executor = executor.Where("comic_table.latest_translated_at IS NOT NULL")
			}
		}

		if proofreadStatus != nil {
			switch *proofreadStatus {
			case model.WorkflowPending:
				executor = executor.Where("comic_table.latest_proofreading_at IS NULL")
			case model.WorkflowInProgress:
				executor = executor.Where("comic_table.latest_proofreading_at IS NOT NULL").
					Where("comic_table.latest_proofread_at IS NULL")
			case model.WorkflowCompleted:
				executor = executor.Where("comic_table.latest_proofread_at IS NOT NULL")
			}
		}

		if typesetStatus != nil {
			switch *typesetStatus {
			case model.WorkflowPending:
				executor = executor.Where("comic_table.latest_typesetting_at IS NULL")
			case model.WorkflowInProgress:
				executor = executor.Where("comic_table.latest_typesetting_at IS NOT NULL").
					Where("comic_table.latest_typeset_at IS NULL")
			case model.WorkflowCompleted:
				executor = executor.Where("comic_table.latest_typeset_at IS NOT NULL")
			}
		}

		if reviewStatus != nil {
			switch *reviewStatus {
			case model.WorkflowPending:
				executor = executor.Where("comic_table.latest_reviewed_at IS NULL")
			case model.WorkflowCompleted:
				executor = executor.Where("comic_table.latest_reviewed_at IS NOT NULL")
			}
		}

		if publishStatus != nil {
			switch *publishStatus {
			case model.WorkflowPending:
				executor = executor.Where("comic_table.latest_published_at IS NULL")
			case model.WorkflowCompleted:
				executor = executor.Where("comic_table.latest_published_at IS NOT NULL")
			}
		}

		return executor
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
		return executor.
			Joins("LEFT JOIN workset_table ON workset_table.id = comic_table.workset_id").
			Select(
				"comic_table.*, workset_table.team_id AS workset_team_id, " +
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
				"comic_table.*, " +
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
			Joins("LEFT JOIN workset_table ON workset_table.id = comic_table.workset_id").
			Joins("LEFT JOIN user_table AS creator_table ON creator_table.id = comic_table.creator_id").
			Select(
				"comic_table.*, workset_table.team_id AS workset_team_id, " +
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
