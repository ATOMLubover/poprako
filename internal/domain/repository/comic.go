package repository

import "labelplus-next-web-be/internal/domain/model"

type ComicRepository interface {
	List(executor Executor, options ...QueryOption) ([]*model.ComicInfo, error)
	Create(executor Executor, creation *model.ComicCreation) (string, error)
	Update(executor Executor, update *model.ComicUpdate) error
	Delete(executor Executor, comicID string) error
}
