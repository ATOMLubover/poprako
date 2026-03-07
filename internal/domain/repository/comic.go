package repository

import "labelplus-next-web-be/internal/domain/model"

type ComicRepository interface {
	Transactor
	List(executor Executor, options ...QueryOption) ([]model.ComicInfo, error)
	Get(executor Executor, options ...QueryOption) (model.ComicInfo, error)
	Count(executor Executor, options ...QueryOption) (int64, error)
	LockByTeamID(executor Executor, teamID string) error
	Create(executor Executor, creation model.ComicCreation) (string, error)
	Update(executor Executor, update model.ComicUpdate) error
	Delete(executor Executor, id string) error
}
