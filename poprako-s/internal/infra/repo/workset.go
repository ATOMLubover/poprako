package repo_infra

import (
	"context"
	"errors"

	"poprako-s/internal/domain/model"
	iface "poprako-s/internal/domain/repo"

	"gorm.io/gorm"
)

type worksetRepoImpl struct {
	gdb *gorm.DB
}

func NewWorksetRepo(gdb *gorm.DB) iface.WorksetRepo {
	return &worksetRepoImpl{gdb: gdb}
}

func NewWorksetRepoFromCx(cx context.Context) (iface.WorksetRepo, error) {
	gdb, ok := cx.Value(txnKey).(*gorm.DB)
	if !ok {
		return nil, errors.New("[NewWorksetRepoFromCx]: 无法从上下文中获取事务数据库连接")
	}

	return &worksetRepoImpl{gdb: gdb}, nil
}

func (r *worksetRepoImpl) FromTxnCx(cx context.Context) (iface.WorksetRepo, error) {
	return NewWorksetRepoFromCx(cx)
}

func (r *worksetRepoImpl) GetByID(id string) (*model.WorksetInfo, error) {
	return nil, errors.New("not implemented")
}

func (r *worksetRepoImpl) List(opt model.WorksetQueryOpt) ([]model.WorksetInfo, error) {
	return nil, errors.New("not implemented")
}

func (r *worksetRepoImpl) Count(opt model.WorksetQueryOpt) (int64, error) {
	return 0, errors.New("not implemented")
}

func (r *worksetRepoImpl) Create(c *model.WorksetCreation) (*model.WorksetInfo, error) {
	return nil, errors.New("not implemented")
}

func (r *worksetRepoImpl) Update(u *model.WorksetUpdate) error {
	return errors.New("not implemented")
}

func (r *worksetRepoImpl) UpdateComicCount(id string, delta int) error {
	return errors.New("not implemented")
}

func (r *worksetRepoImpl) Delete(id string) error {
	return errors.New("not implemented")
}
