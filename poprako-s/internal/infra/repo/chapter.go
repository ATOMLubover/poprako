package repo_infra

import (
	"context"
	"errors"

	"poprako-s/internal/domain/model"
	iface "poprako-s/internal/domain/repo"

	"gorm.io/gorm"
)

type chapterRepoImpl struct {
	gdb *gorm.DB
}

func NewChapterRepo(gdb *gorm.DB) iface.ChapterRepo {
	return &chapterRepoImpl{gdb: gdb}
}

func NewChapterRepoFromCx(cx context.Context) (iface.ChapterRepo, error) {
	gdb, ok := cx.Value(txnKey).(*gorm.DB)
	if !ok {
		return nil, errors.New("[NewChapterRepoFromCx]: 无法从上下文中获取事务数据库连接")
	}

	return &chapterRepoImpl{gdb: gdb}, nil
}

func (r *chapterRepoImpl) FromTxnCx(cx context.Context) (iface.ChapterRepo, error) {
	return NewChapterRepoFromCx(cx)
}

func (r *chapterRepoImpl) GetByID(id string) (*model.ChapterInfo, error) {
	return nil, errors.New("not implemented")
}

func (r *chapterRepoImpl) FindPinnedByComicID(comicID string) (*model.ChapterInfo, error) {
	return nil, errors.New("not implemented")
}

func (r *chapterRepoImpl) List(opt model.ChapterQueryOpt) ([]model.ChapterInfo, error) {
	return nil, errors.New("not implemented")
}

func (r *chapterRepoImpl) Count(opt model.ChapterQueryOpt) (int64, error) {
	return 0, errors.New("not implemented")
}

func (r *chapterRepoImpl) Create(c *model.ChapterCreation) (*model.ChapterInfo, error) {
	return nil, errors.New("not implemented")
}

func (r *chapterRepoImpl) Update(u *model.ChapterUpdate) error {
	return errors.New("not implemented")
}

func (r *chapterRepoImpl) Remove(id string) error {
	return errors.New("not implemented")
}

func (r *chapterRepoImpl) UpdateStats(stats *model.ChapterStats) error {
	return errors.New("not implemented")
}
