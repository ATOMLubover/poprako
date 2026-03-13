package repository

import (
	"labelplus-next-web-be/internal/domain/model"
	intf "labelplus-next-web-be/internal/domain/repository"
	"labelplus-next-web-be/internal/infrastructure/repository/entity"
	"labelplus-next-web-be/internal/util"

	"gorm.io/gorm/clause"
)

type worksetRepository struct {
	executor intf.Executor
}

func NewWorksetRepository(executor intf.Executor) intf.WorksetRepository {
	return &worksetRepository{executor: executor}
}

func (r *worksetRepository) withTransaction(executor intf.Executor) intf.Executor {
	if executor != nil {
		return executor
	}

	return r.executor
}

func (r *worksetRepository) BeginTransaction() intf.Executor {
	return r.executor.Begin()
}

func (r *worksetRepository) List(executor intf.Executor, options ...intf.QueryOption) ([]model.WorksetInfo, error) {
	executor = r.withTransaction(executor)
	executor = executor.Table(entity.WorksetTable)
	for _, opt := range options {
		executor = opt(executor)
	}

	var rows []entity.WorksetInfoRow
	if err := executor.Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]model.WorksetInfo, len(rows))
	for i, row := range rows {
		result[i] = entity.ToWorksetInfo(row)
	}

	return result, nil
}

func (r *worksetRepository) Get(executor intf.Executor, options ...intf.QueryOption) (model.WorksetInfo, error) {
	executor = r.withTransaction(executor)
	executor = executor.Table(entity.WorksetTable)
	for _, opt := range options {
		executor = opt(executor)
	}

	var row entity.WorksetInfoRow
	if err := executor.First(&row).Error; err != nil {
		return model.WorksetInfo{}, err
	}

	return entity.ToWorksetInfo(row), nil
}

func (r *worksetRepository) LockByTeamID(executor intf.Executor, teamID string) error {
	executor = r.withTransaction(executor)

	var lockedIDs []string

	return executor.
		Table(entity.WorksetTable).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("team_id = ?", teamID).
		Pluck("id", &lockedIDs).Error
}

func (r *worksetRepository) Count(executor intf.Executor, options ...intf.QueryOption) (int64, error) {
	executor = r.withTransaction(executor)
	executor = executor.Table(entity.WorksetTable)
	for _, opt := range options {
		executor = opt(executor)
	}

	var count int64
	if err := executor.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *worksetRepository) Create(executor intf.Executor, creation model.WorksetCreation) (string, error) {
	executor = r.withTransaction(executor)

	row := entity.WorksetInsertRow{
		ID:          util.GenerateUUID(),
		TeamID:      creation.TeamID,
		Index:       creation.Index,
		Name:        creation.Name,
		Description: creation.Description,
	}

	if err := executor.Create(&row).Error; err != nil {
		return "", err
	}

	return row.ID, nil
}

func (r *worksetRepository) Update(executor intf.Executor, update model.WorksetUpdate) error {
	executor = r.withTransaction(executor)

	updates := map[string]any{
		"name": update.Name,
	}

	if update.Description != nil {
		updates["description"] = *update.Description
	} else {
		updates["description"] = nil
	}

	return executor.
		Table(entity.WorksetTable).
		Where("id = ?", update.ID).
		Updates(updates).Error
}

func (r *worksetRepository) Delete(executor intf.Executor, id string) error {
	executor = r.withTransaction(executor)

	return executor.
		Table(entity.WorksetTable).
		Where("id = ?", id).
		Delete(nil).Error
}
