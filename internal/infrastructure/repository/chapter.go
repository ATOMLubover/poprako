package repository

import (
	"time"

	"labelplus-next-web-be/internal/domain/model"
	intf "labelplus-next-web-be/internal/domain/repository"
	"labelplus-next-web-be/internal/infrastructure/repository/entity"
	"labelplus-next-web-be/internal/util"

	"gorm.io/gorm/clause"
)

type chapterRepository struct {
	executor intf.Executor
}

func NewChapterRepository(executor intf.Executor) intf.ChapterRepository {
	return &chapterRepository{executor: executor}
}

func (r *chapterRepository) withTransaction(executor intf.Executor) intf.Executor {
	if executor != nil {
		return executor
	}

	return r.executor
}

func (r *chapterRepository) BeginTransaction() intf.Executor {
	return r.executor.Begin()
}

func (r *chapterRepository) LockByID(executor intf.Executor, chapterID string) error {
	executor = r.withTransaction(executor)

	var lockedID string

	return executor.
		Table(entity.ChapterTable).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ? AND deleted_at IS NULL", chapterID).
		Pluck("id", &lockedID).Error
}

func (r *chapterRepository) LockByComicID(executor intf.Executor, comicID string) error {
	executor = r.withTransaction(executor)

	var lockedIDs []string

	return executor.
		Table(entity.ChapterTable).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("comic_id = ? AND deleted_at IS NULL", comicID).
		Pluck("id", &lockedIDs).Error
}

func (r *chapterRepository) List(executor intf.Executor, options ...intf.QueryOption) ([]model.ChapterInfo, error) {
	executor = r.withTransaction(executor)
	executor = executor.Table(entity.ChapterTable).Where("deleted_at IS NULL")

	for _, opt := range options {
		executor = opt(executor)
	}

	var rows []entity.ChapterWithInfoRow
	if err := executor.Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]model.ChapterInfo, len(rows))
	for i, row := range rows {
		result[i] = entity.ToChapterWithInfo(row)
	}

	return result, nil
}

func (r *chapterRepository) Get(executor intf.Executor, options ...intf.QueryOption) (model.ChapterInfo, error) {
	executor = r.withTransaction(executor)
	executor = executor.Table(entity.ChapterTable).Where("deleted_at IS NULL")

	for _, opt := range options {
		executor = opt(executor)
	}

	var row entity.ChapterWithInfoRow
	if err := executor.First(&row).Error; err != nil {
		return model.ChapterInfo{}, err
	}

	return entity.ToChapterWithInfo(row), nil
}

func (r *chapterRepository) GetStatsByID(executor intf.Executor, chapterID string) (model.ChapterStats, error) {
	executor = r.withTransaction(executor)

	var row entity.ChapterInfoRow
	if err := executor.
		Table(entity.ChapterTable).
		Select("id, total_unit_count, translated_unit_count, proofread_unit_count").
		Where("id = ? AND deleted_at IS NULL", chapterID).
		First(&row).Error; err != nil {
		return model.ChapterStats{}, err
	}

	return model.NewChapterStats(
		row.ID,
		row.TotalUnitCount,
		row.TranslatedUnitCount,
		row.ProofreadUnitCount,
	), nil
}

func (r *chapterRepository) Count(executor intf.Executor, options ...intf.QueryOption) (int64, error) {
	executor = r.withTransaction(executor)
	executor = executor.Table(entity.ChapterTable).Where("deleted_at IS NULL")

	for _, opt := range options {
		executor = opt(executor)
	}

	var count int64
	if err := executor.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *chapterRepository) Create(executor intf.Executor, creation model.ChapterCreation) (string, error) {
	executor = r.withTransaction(executor)

	row := entity.ChapterInsertRow{
		ID:        util.GenerateUUID(),
		ComicID:   creation.ComicID,
		Index:     creation.Index,
		ChapterNo: creation.ChapterNo,
		CreatorID: creation.CreatorID,
	}

	if err := executor.Create(&row).Error; err != nil {
		return "", err
	}

	return row.ID, nil
}

func (r *chapterRepository) Update(executor intf.Executor, update model.ChapterUpdate) error {
	executor = r.withTransaction(executor)

	return executor.
		Table(entity.ChapterTable).
		Where("id = ? AND deleted_at IS NULL", update.ID).
		Updates(map[string]any{
			"subtitle":        update.ChapterNo,
			"uploaded_at":     update.UploadedAt,
			"transalating_at": update.TransalatingAt,
			"translated_at":   update.TranslatedAt,
			"proofreading_at": update.ProofreadingAt,
			"proofread_at":    update.ProofreadAt,
			"typesetting_at":  update.TypesettingAt,
			"typeset_at":      update.TypesetAt,
			"reviewed_at":     update.ReviewedAt,
			"published_at":    update.PublishedAt,
		}).Error
}

func (r *chapterRepository) UpdateStats(executor intf.Executor, chapterStats model.ChapterStats) error {
	executor = r.withTransaction(executor)

	return executor.
		Table(entity.ChapterTable).
		Where("id = ? AND deleted_at IS NULL", chapterStats.ChapterID).
		Updates(map[string]any{
			"total_unit_count":      chapterStats.TotalUnitCount,
			"translated_unit_count": chapterStats.TranslatedUnitCount,
			"proofread_unit_count":  chapterStats.ProofreadUnitCount,
		}).Error
}

func (r *chapterRepository) Delete(executor intf.Executor, id string) error {
	executor = r.withTransaction(executor)

	return executor.
		Table(entity.ChapterTable).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).Error
}
