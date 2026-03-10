package repository

import (
	"time"

	"labelplus-next-web-be/internal/domain/model"
	intf "labelplus-next-web-be/internal/domain/repository"
	"labelplus-next-web-be/internal/infrastructure/repository/entity"
	"labelplus-next-web-be/internal/util"

	"gorm.io/gorm/clause"
)

type comicRepository struct {
	executor intf.Executor
}

func NewComicRepository(executor intf.Executor) intf.ComicRepository {
	return &comicRepository{executor: executor}
}

func (r *comicRepository) withTransaction(executor intf.Executor) intf.Executor {
	if executor != nil {
		return executor
	}

	return r.executor
}

func (r *comicRepository) BeginTransaction() intf.Executor {
	return r.executor.Begin()
}

func (r *comicRepository) List(executor intf.Executor, options ...intf.QueryOption) ([]model.ComicInfo, error) {
	executor = r.withTransaction(executor)

	executor = executor.Table(entity.ComicTable).Where("deleted_at IS NULL")
	for _, opt := range options {
		executor = opt(executor)
	}

	var rows []entity.ComicInfoRow

	if err := executor.Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]model.ComicInfo, len(rows))
	for i, row := range rows {
		result[i] = entity.ToComicInfo(row)
	}

	return result, nil
}

func (r *comicRepository) Get(executor intf.Executor, options ...intf.QueryOption) (model.ComicInfo, error) {
	executor = r.withTransaction(executor)
	executor = executor.Table(entity.ComicTable).Where("deleted_at IS NULL")
	for _, opt := range options {
		executor = opt(executor)
	}

	var row entity.ComicInfoRow
	if err := executor.First(&row).Error; err != nil {
		return model.ComicInfo{}, err
	}

	return entity.ToComicInfo(row), nil
}

func (r *comicRepository) LockByTeamID(executor intf.Executor, teamID string) error {
	executor = r.withTransaction(executor)

	var lockedIDs []string

	return executor.
		Table(entity.ComicTable).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("team_id = ? AND deleted_at IS NULL", teamID).
		Pluck("id", &lockedIDs).Error
}

func (r *comicRepository) Count(executor intf.Executor, options ...intf.QueryOption) (int64, error) {
	executor = r.withTransaction(executor)
	executor = executor.Table(entity.ComicTable).Where("deleted_at IS NULL")
	for _, opt := range options {
		executor = opt(executor)
	}

	var count int64
	if err := executor.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *comicRepository) Create(executor intf.Executor, creation model.ComicCreation) (string, error) {
	executor = r.withTransaction(executor)

	row := entity.ComicInsertRow{
		ID:          util.GenerateUUID(),
		TeamID:      creation.TeamID,
		Index:       creation.Index,
		Title:       creation.Title,
		Author:      creation.Author,
		Description: creation.Description,
		CreatorID:   creation.CreatorID,
		CoverURL:    "",
	}

	if err := executor.Create(&row).Error; err != nil {
		return "", err
	}

	return row.ID, nil
}

func (r *comicRepository) Update(executor intf.Executor, update model.ComicUpdate) error {
	executor = r.withTransaction(executor)

	updates := map[string]any{}
	if update.Title != nil {
		updates["title"] = *update.Title
	}
	if update.Author != nil {
		updates["author"] = *update.Author
	}
	if update.Description != nil {
		updates["description"] = *update.Description
	}

	if len(updates) == 0 {
		return nil
	}

	return executor.
		Table(entity.ComicTable).
		Where("id = ? AND deleted_at IS NULL", update.ID).
		Updates(updates).Error
}

func (r *comicRepository) Delete(executor intf.Executor, comicID string) error {
	executor = r.withTransaction(executor)

	return executor.
		Table(entity.ComicTable).
		Where("id = ?", comicID).
		Update("deleted_at", time.Now()).Error
}
