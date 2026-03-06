package repository

import "labelplus-next-web-be/internal/domain/model"

type ComicRepository interface {
	Transactor
	List(executor Executor, options ...QueryOption) ([]*model.ComicInfo, error)
	GetByID(executor Executor, comicID string) (*model.ComicInfo, error)
	LockByTeamID(executor Executor, teamID string) error
	CountByTeamID(executor Executor, teamID string) (int64, error)
	Create(executor Executor, creation *model.ComicCreation) (string, error)
	Update(executor Executor, update *model.ComicUpdate) error
	Delete(executor Executor, comicID string) error
}
