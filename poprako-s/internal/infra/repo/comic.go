package repo_infra

import (
	"context"
	"errors"

	"poprako-s/internal/domain/model"
	iface "poprako-s/internal/domain/repo"

	"gorm.io/gorm"
)

type comicRepoImpl struct {
	gdb *gorm.DB
}

func NewComicRepo(gdb *gorm.DB) iface.ComicRepo {
	return &comicRepoImpl{gdb: gdb}
}

func NewComicRepoFromCx(cx context.Context) (iface.ComicRepo, error) {
	gdb, ok := cx.Value(txnKey).(*gorm.DB)
	if !ok {
		return nil, errors.New("[NewComicRepoFromCx]: 无法从上下文中获取事务数据库连接")
	}

	return &comicRepoImpl{gdb: gdb}, nil
}

func (r *comicRepoImpl) FromTxnCx(cx context.Context) (iface.ComicRepo, error) {
	return NewComicRepoFromCx(cx)
}

func (r *comicRepoImpl) GetByID(id string) (*model.ComicInfo, error) {
	return nil, errors.New("not implemented")
}

func (r *comicRepoImpl) List(opt model.ComicQueryOpt) ([]model.ComicInfo, error) {
	return nil, errors.New("not implemented")
}

func (r *comicRepoImpl) Count(opt model.ComicQueryOpt) (int64, error) {
	return 0, errors.New("not implemented")
}

func (r *comicRepoImpl) Create(c *model.ComicCreation) (*model.ComicInfo, error) {
	return nil, errors.New("not implemented")
}

func (r *comicRepoImpl) Update(u *model.ComicUpdate) error {
	return errors.New("not implemented")
}

func (r *comicRepoImpl) UpdateChapterCount(id string, delta int) error {
	return errors.New("not implemented")
}

func (r *comicRepoImpl) Delete(id string) error {
	return errors.New("not implemented")
}
