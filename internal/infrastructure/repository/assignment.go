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

func (r *assignmentRepository) List(executor intf.Executor, options ...intf.QueryOption) ([]model.AssignmentInfo, error) {
	executor = r.withTransaction(executor)
	executor = executor.Table(entity.AssignmentTable)

	for _, opt := range options {
		executor = opt(executor)
	}

	var rows []entity.AssignmentIncludeRow
	if err := executor.Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]model.AssignmentInfo, len(rows))
	for i, row := range rows {
		result[i] = entity.ToAssignmentInfoFromIncludeRow(row)
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
		AssignedRedrawerAt:    creation.AssignedRedrawerAt,
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
		"assigned_redrawer_at":     update.AssignedRedrawerAt,
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
