package repository

import "labelplus-next-web-be/internal/domain/model"

type ChapterRepository interface {
	Transactor
	LockByID(executor Executor, chapterID string) error
	LockByComicID(executor Executor, comicID string) error
	List(executor Executor, options ...QueryOption) ([]model.ChapterInfo, error)
	Get(executor Executor, options ...QueryOption) (model.ChapterInfo, error)
	GetStatsByID(executor Executor, chapterID string) (model.ChapterStats, error)
	Count(executor Executor, options ...QueryOption) (int64, error)
	Create(executor Executor, creation model.ChapterCreation) (string, error)
	Update(executor Executor, update model.ChapterUpdate) error
	// UpdateStats 仅更新章节的三个统计字段，供 unit save 事务使用。
	UpdateStats(executor Executor, chapterStats model.ChapterStats) error
	Delete(executor Executor, id string) error
}
