package repository

import "labelplus-next-web-be/internal/domain/model"

type PageRepository interface {
	Transactor
	LockByChapterID(executor Executor, chapterID string) error
	List(executor Executor, options ...QueryOption) ([]model.PageInfo, error)
	Get(executor Executor, options ...QueryOption) (model.PageInfo, error)
	GetChapterByID(executor Executor, chapterID string) (model.ChapterInfo, error)
	CreateBatch(executor Executor, pages []model.PageCreation) error
	Update(executor Executor, update model.PageUpdate) error
	Delete(executor Executor, id string) error
	DeleteBatch(executor Executor, ids []string) error
}
