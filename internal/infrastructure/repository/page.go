package repository

import (
	"labelplus-next-web-be/internal/domain/model"
	intf "labelplus-next-web-be/internal/domain/repository"
	"labelplus-next-web-be/internal/infrastructure/repository/entity"

	"gorm.io/gorm/clause"
)

type pageRepository struct {
	executor intf.Executor
}

func NewPageRepository(executor intf.Executor) intf.PageRepository {
	return &pageRepository{executor: executor}
}

func (r *pageRepository) withTransaction(executor intf.Executor) intf.Executor {
	if executor != nil {
		return executor
	}

	return r.executor
}

func (r *pageRepository) BeginTransaction() intf.Executor {
	return r.executor.Begin()
}

func (r *pageRepository) LockByChapterID(executor intf.Executor, chapterID string) error {
	executor = r.withTransaction(executor)

	var lockedIDs []string

	return executor.
		Table(entity.PageTable).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("chapter_id = ?", chapterID).
		Pluck("id", &lockedIDs).Error
}

func (r *pageRepository) List(executor intf.Executor, options ...intf.QueryOption) ([]model.PageInfo, error) {
	executor = r.withTransaction(executor)
	executor = executor.Table(entity.PageTable)

	for _, opt := range options {
		executor = opt(executor)
	}

	var rows []entity.PageInfoRow
	if err := executor.Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]model.PageInfo, len(rows))
	for i, row := range rows {
		result[i] = entity.ToPageInfo(row)
	}

	return result, nil
}

func (r *pageRepository) Get(executor intf.Executor, options ...intf.QueryOption) (model.PageInfo, error) {
	executor = r.withTransaction(executor)
	executor = executor.Table(entity.PageTable)

	for _, opt := range options {
		executor = opt(executor)
	}

	var row entity.PageInfoRow
	if err := executor.First(&row).Error; err != nil {
		return model.PageInfo{}, err
	}

	return entity.ToPageInfo(row), nil
}

func (r *pageRepository) GetChapterByID(executor intf.Executor, chapterID string) (model.ChapterDetail, error) {
	executor = r.withTransaction(executor)
	executor = executor.Table(entity.ChapterTable).Where("id = ?", chapterID)

	var row entity.ChapterInfoRow
	if err := executor.First(&row).Error; err != nil {
		return model.ChapterDetail{}, err
	}

	return entity.ToChapterInfo(row), nil
}

func (r *pageRepository) CreateBatch(executor intf.Executor, pages []model.PageCreation) error {
	executor = r.withTransaction(executor)

	rows := make([]entity.PageInsertRow, len(pages))
	for i := range pages {
		rows[i] = entity.PageInsertRow{
			ID:        pages[i].ID,
			ChapterID: pages[i].ChapterID,
			Index:     pages[i].Index,
			OSSKey:    pages[i].OSSKey,
			CreatorID: pages[i].CreatorID,
		}
	}

	if len(rows) == 0 {
		return nil
	}

	return executor.Create(&rows).Error
}

func (r *pageRepository) Update(executor intf.Executor, update model.PageUpdate) error {
	executor = r.withTransaction(executor)

	return executor.
		Table(entity.PageTable).
		Where("id = ?", update.ID).
		Updates(map[string]any{
			"index":                 update.Index,
			"oss_key":               update.OSSKey,
			"uploaded":              update.IsUploaded,
			"total_unit_count":      update.TotalUnitCount,
			"translated_unit_count": update.TranslatedUnitCount,
			"proved_unit_count":     update.ProofreadUnitCount,
		}).Error
}

func (r *pageRepository) Delete(executor intf.Executor, id string) error {
	executor = r.withTransaction(executor)
	executor = executor.Table(entity.PageTable)

	return executor.Where("id = ?", id).Delete(&entity.PageInfoRow{}).Error
}

func (r *pageRepository) DeleteBatch(executor intf.Executor, ids []string) error {
	executor = r.withTransaction(executor)
	executor = executor.Table(entity.PageTable)

	if len(ids) == 0 {
		return nil
	}

	return executor.Where("id IN ?", ids).Delete(&entity.PageInfoRow{}).Error
}
