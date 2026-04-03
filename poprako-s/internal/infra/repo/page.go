package repo_infra

import (
	"context"
	"errors"

	"poprako-s/internal/domain/model"
	iface "poprako-s/internal/domain/repo"

	"gorm.io/gorm"
)

type pageRepoImpl struct {
	gdb *gorm.DB
}

func NewPageRepo(gdb *gorm.DB) iface.PageRepo {
	return &pageRepoImpl{gdb: gdb}
}

func NewPageRepoFromCx(cx context.Context) (iface.PageRepo, error) {
	gdb, ok := cx.Value(txnKey).(*gorm.DB)
	if !ok {
		return nil, errors.New("[NewPageRepoFromCx]: 无法从上下文中获取事务数据库连接")
	}

	return &pageRepoImpl{gdb: gdb}, nil
}

func (r *pageRepoImpl) FromTxnCx(cx context.Context) (iface.PageRepo, error) {
	return NewPageRepoFromCx(cx)
}

func (r *pageRepoImpl) GetByID(id string) (*model.PageInfo, error) {
	return nil, errors.New("not implemented")
}

func (r *pageRepoImpl) List(opt model.PageQueryOpt) ([]model.PageInfo, error) {
	return nil, errors.New("not implemented")
}

func (r *pageRepoImpl) GetStatsByID(pageID string) (*model.PageStats, error) {
	return nil, errors.New("not implemented")
}

func (r *pageRepoImpl) CreateBatch(pages []*model.PageCreation) error {
	return errors.New("not implemented")
}

func (r *pageRepoImpl) Update(u *model.PageUpdate) error {
	return errors.New("not implemented")
}

func (r *pageRepoImpl) UpdateStats(stats *model.PageStats) error {
	return errors.New("not implemented")
}

func (r *pageRepoImpl) Delete(id string) error {
	return errors.New("not implemented")
}

func (r *pageRepoImpl) DeleteBatch(ids []string) error {
	return errors.New("not implemented")
}
