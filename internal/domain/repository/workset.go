package repository

import "labelplus-next-web-be/internal/domain/model"

type WorksetRepository interface {
	Transactor
	List(executor Executor, options ...QueryOption) ([]model.WorksetInfo, error)
	Get(executor Executor, options ...QueryOption) (model.WorksetInfo, error)
	Count(executor Executor, options ...QueryOption) (int64, error)
	LockByTeamID(executor Executor, teamID string) error
	Create(executor Executor, creation model.WorksetCreation) (string, error)
	Update(executor Executor, update model.WorksetUpdate) error
	Delete(executor Executor, id string) error
}
