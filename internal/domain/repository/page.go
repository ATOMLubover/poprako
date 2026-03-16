package repository

import "labelplus-next-web-be/internal/domain/model"

type PageRepository interface {
	Transactor
	LockByID(executor Executor, pageID string) error
	LockByChapterID(executor Executor, chapterID string) error
	List(executor Executor, options ...QueryOption) ([]model.PageInfo, error)
	Get(executor Executor, options ...QueryOption) (model.PageInfo, error)
	GetStatsByID(executor Executor, pageID string) (model.PageStats, error)
	CreateBatch(executor Executor, pages []model.PageCreation) error
	UpdateStats(executor Executor, pageStats model.PageStats) error
	Update(executor Executor, update model.PageUpdate) error
	Delete(executor Executor, id string) error
	DeleteBatch(executor Executor, ids []string) error
}
