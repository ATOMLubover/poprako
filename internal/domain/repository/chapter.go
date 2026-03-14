package repository

import "labelplus-next-web-be/internal/domain/model"

type ChapterRepository interface {
	Transactor
	LockByComicID(executor Executor, comicID string) error
	List(executor Executor, options ...QueryOption) ([]model.ChapterInfo, error)
	Get(executor Executor, options ...QueryOption) (model.ChapterInfo, error)
	Count(executor Executor, options ...QueryOption) (int64, error)
	Create(executor Executor, creation model.ChapterCreation) (string, error)
	Update(executor Executor, update model.ChapterUpdate) error
	Delete(executor Executor, id string) error
}
