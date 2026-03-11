package repository

import (
	"labelplus-next-web-be/internal/domain/model"
	intf "labelplus-next-web-be/internal/domain/repository"
	"labelplus-next-web-be/internal/infrastructure/repository/entity"
	"labelplus-next-web-be/internal/util"
)

type assignmentRepository struct {
	executor intf.Executor
}

func NewAssignmentRepository(executor intf.Executor) intf.AssignmentRepository {
	return &assignmentRepository{executor: executor}
}

func (r *assignmentRepository) withTransaction(executor intf.Executor) intf.Executor {
	if executor != nil {
		return executor
	}
	return r.executor
}

func (r *assignmentRepository) BeginTransaction() intf.Executor {
	return r.executor.Begin()
}

func (r *assignmentRepository) Get(executor intf.Executor, options ...intf.QueryOption) (model.AssignmentInfo, error) {
	executor = r.withTransaction(executor)
	executor = executor.Table(entity.AssignmentTable)

	for _, opt := range options {
		executor = opt(executor)
	}

	var row entity.AssignmentInfoRow
	if err := executor.First(&row).Error; err != nil {
		return model.AssignmentInfo{}, err
	}

	return entity.ToAssignmentInfo(row), nil
}

func (r *assignmentRepository) Exist(executor intf.Executor, options ...intf.QueryOption) (bool, error) {
	executor = r.withTransaction(executor)
	executor = executor.Table(entity.AssignmentTable)

	for _, opt := range options {
		executor = opt(executor)
	}

	var count int64
	if err := executor.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *assignmentRepository) ListWithUserInfo(executor intf.Executor, options ...intf.QueryOption) ([]model.AssignmentWithUserInfo, error) {
	executor = r.withTransaction(executor)

	executor = executor.Table(entity.AssignmentTable).
		Select(`assignment_table.*,
			user_table.name               AS user_name,
			user_table.qq                 AS user_qq,
			user_table.avatar_oss_key     AS user_avatar_oss_key,
			user_table.is_avatar_uploaded AS user_is_avatar_uploaded,
			user_table.is_super_admin     AS user_is_super_admin,
			user_table.created_at         AS user_created_at,
			user_table.updated_at         AS user_updated_at`).
		Joins("LEFT JOIN user_table ON user_table.id = assignment_table.user_id AND user_table.deleted_at IS NULL")

	for _, opt := range options {
		executor = opt(executor)
	}

	var rows []entity.AssignmentWithUserRow
	if err := executor.Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]model.AssignmentWithUserInfo, len(rows))
	for i, row := range rows {
		result[i] = entity.ToAssignmentWithUserInfo(row)
	}

	return result, nil
}

func (r *assignmentRepository) ListWithChapterInfo(executor intf.Executor, options ...intf.QueryOption) ([]model.AssignmentWithChapterInfo, error) {
	executor = r.withTransaction(executor)

	executor = executor.Table(entity.AssignmentTable).
		Select(`assignment_table.*,
			chapter_table.comic_id      AS chapter_comic_id,
			chapter_table.index         AS chapter_index,
			chapter_table.subtitle      AS chapter_subtitle,
			chapter_table.page_count    AS chapter_page_count,
			chapter_table.cover_url     AS chapter_cover_url,
			chapter_table.created_at    AS chapter_created_at,
			chapter_table.updated_at    AS chapter_updated_at,
			comic_table.team_id         AS comic_team_id,
			comic_table.index           AS comic_index,
			comic_table.title           AS comic_title,
			comic_table.author          AS comic_author,
			comic_table.description     AS comic_description,
			comic_table.cover_url       AS comic_cover_url,
			comic_table.chapter_count   AS comic_chapter_count,
			comic_table.creator_id      AS comic_creator_id,
			comic_table.last_active_at  AS comic_last_active_at,
			comic_table.created_at      AS comic_created_at,
			comic_table.updated_at      AS comic_updated_at`).
		Joins("LEFT JOIN chapter_table ON chapter_table.id = assignment_table.chapter_id AND chapter_table.deleted_at IS NULL").
		Joins("LEFT JOIN comic_table ON comic_table.id = chapter_table.comic_id AND comic_table.deleted_at IS NULL")

	for _, opt := range options {
		executor = opt(executor)
	}

	var rows []entity.AssignmentWithChapterAndComicRow
	if err := executor.Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]model.AssignmentWithChapterInfo, len(rows))
	for i, row := range rows {
		result[i] = entity.ToAssignmentWithChapterInfo(row)
	}

	return result, nil
}

func (r *assignmentRepository) Create(executor intf.Executor, creation model.AssignmentCreation) (string, error) {
	executor = r.withTransaction(executor)

	row := entity.AssignmentInsertRow{
		ID:                    util.GenerateUUID(),
		ChapterID:             creation.ChapterID,
		UserID:                creation.UserID,
		AssignedRawProviderAt: creation.AssignedRawProviderAt,
		AssignedTranslatorAt:  creation.AssignedTranslatorAt,
		AssignedProofreaderAt: creation.AssignedProofreaderAt,
		AssignedTypesetterAt:  creation.AssignedTypesetterAt,
		AssignedReviewerAt:    creation.AssignedReviewerAt,
		AssignedPublisherAt:   creation.AssignedPublisherAt,
	}

	if err := executor.Create(&row).Error; err != nil {
		return "", err
	}

	return row.ID, nil
}

func (r *assignmentRepository) Update(executor intf.Executor, update model.AssignmentUpdate) error {
	executor = r.withTransaction(executor)

	updates := map[string]any{
		"assigned_raw_provider_at": update.AssignedRawProviderAt,
		"assigned_translator_at":   update.AssignedTranslatorAt,
		"assigned_proofreader_at":  update.AssignedProofreaderAt,
		"assigned_typesetter_at":   update.AssignedTypesetterAt,
		"assigned_reviewer_at":     update.AssignedReviewerAt,
		"assigned_publisher_at":    update.AssignedPublisherAt,
	}

	return executor.
		Table(entity.AssignmentTable).
		Where("id = ?", update.ID).
		Updates(updates).Error
}

func (r *assignmentRepository) Delete(executor intf.Executor, assignmentID string) error {
	executor = r.withTransaction(executor)

	return executor.
		Table(entity.AssignmentTable).
		Where("id = ?", assignmentID).
		Delete(nil).Error
}

func (r *assignmentRepository) LockByComicID(executor intf.Executor, comicID string) error {
	executor = r.withTransaction(executor)

	var lockedIDs []string

	return executor.
		Table(entity.AssignmentTable).
		Where("chapter_id IN (SELECT id FROM chapter_table WHERE comic_id = ?)", comicID).
		Pluck("id", &lockedIDs).Error
}
