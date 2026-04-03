package repo_infra

import (
	"context"
	"errors"

	"poprako-s/internal/domain/model"
	iface "poprako-s/internal/domain/repo"

	"gorm.io/gorm"
)

type unitRepoImpl struct {
	gdb *gorm.DB
}

func NewUnitRepo(gdb *gorm.DB) iface.UnitRepo {
	return &unitRepoImpl{gdb: gdb}
}

func NewUnitRepoFromCx(cx context.Context) (iface.UnitRepo, error) {
	gdb, ok := cx.Value(txnKey).(*gorm.DB)
	if !ok {
		return nil, errors.New("[NewUnitRepoFromCx]: 无法从上下文中获取事务数据库连接")
	}

	return &unitRepoImpl{gdb: gdb}, nil
}

func (r *unitRepoImpl) FromTxnCx(cx context.Context) (iface.UnitRepo, error) {
	return NewUnitRepoFromCx(cx)
}

func (r *unitRepoImpl) List(opt model.UnitQueryOpt) ([]model.UnitInfo, error) {
	return nil, errors.New("not implemented")
}

func (r *unitRepoImpl) CreateBatch(units []*model.UnitCreation) error {
	return errors.New("not implemented")
}

func (r *unitRepoImpl) PatchBatch(patches []*model.UnitPatch) error {
	return errors.New("not implemented")
}

func (r *unitRepoImpl) DeleteBatch(unitIDs []string) error {
	return errors.New("not implemented")
}
