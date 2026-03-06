package repository

import (
	"time"

	"labelplus-next-web-be/internal/domain/model"
	intf "labelplus-next-web-be/internal/domain/repository"
	"labelplus-next-web-be/internal/repository/entity"
	"labelplus-next-web-be/internal/util"
)

type teamRepository struct {
	executor intf.Executor
}

func NewTeamRepository(executor intf.Executor) intf.TeamRepository {
	return &teamRepository{executor: executor}
}

func (r *teamRepository) withTransaction(executor intf.Executor) intf.Executor {
	if executor != nil {
		return executor
	}
	return r.executor
}

func (r *teamRepository) BeginTransaction() intf.Executor {
	return r.executor.Begin()
}

func (r *teamRepository) List(executor intf.Executor, options ...intf.QueryOption) ([]*model.TeamInfo, error) {
	executor = r.withTransaction(executor)

	executor = executor.Table(entity.TeamTable).Where("deleted_at IS NULL")
	for _, opt := range options {
		executor = opt(executor)
	}

	var rows []entity.TeamInfoRow

	if err := executor.Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]*model.TeamInfo, len(rows))
	for i, row := range rows {
		result[i] = entity.ToTeamInfo(row)
	}

	return result, nil
}

func (r *teamRepository) Create(executor intf.Executor, creation *model.TeamCreation) (string, error) {
	executor = r.withTransaction(executor)

	row := entity.TeamInsertRow{
		ID:          util.GenerateUUID(),
		Name:        creation.Name,
		Description: creation.Description,
	}
	if err := executor.Create(&row).Error; err != nil {
		return "", err
	}
	return row.ID, nil
}

func (r *teamRepository) Update(executor intf.Executor, update *model.TeamUpdate) error {
	executor = r.withTransaction(executor)

	updates := map[string]any{}

	if update.Name.State() == util.OptionSome {
		updates["name"] = update.Name.Unwrap()
	}
	if update.Description.State() == util.OptionSome {
		updates["description"] = update.Description.Unwrap()
	}

	if len(updates) == 0 {
		return nil
	}

	return executor.
		Table(entity.TeamTable).
		Where("id = ?", update.ID).
		Updates(updates).Error
}

func (r *teamRepository) DeleteByID(executor intf.Executor, teamID string) error {
	executor = r.withTransaction(executor)

	return executor.
		Table(entity.TeamTable).
		Where("id = ?", teamID).
		Update("deleted_at", time.Now()).Error
}
