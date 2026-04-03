package repo_infra

import (
	"context"
	"errors"

	"poprako-s/internal/domain/model"
	iface "poprako-s/internal/domain/repo"

	"gorm.io/gorm"
)

type memberRepoImpl struct {
	gdb *gorm.DB
}

func NewMemberRepo(gdb *gorm.DB) iface.MemberRepo {
	return &memberRepoImpl{gdb: gdb}
}

func NewMemberRepoFromCx(cx context.Context) (iface.MemberRepo, error) {
	gdb, ok := cx.Value(txnKey).(*gorm.DB)
	if !ok {
		return nil, errors.New("[NewMemberRepoFromCx]: 无法从上下文中获取事务数据库连接")
	}

	return &memberRepoImpl{gdb: gdb}, nil
}

func (r *memberRepoImpl) FromTxnCx(cx context.Context) (iface.MemberRepo, error) {
	return NewMemberRepoFromCx(cx)
}

func (r *memberRepoImpl) GetByID(id string) (*model.MemberInfo, error) {
	return nil, errors.New("not implemented")
}

func (r *memberRepoImpl) Get(opt model.MemberQueryOpt) (*model.MemberInfo, error) {
	return nil, errors.New("not implemented")
}

func (r *memberRepoImpl) List(opt model.MemberQueryOpt) ([]model.MemberInfo, error) {
	return nil, errors.New("not implemented")
}

func (r *memberRepoImpl) Exist(opt model.MemberQueryOpt) (bool, error) {
	return false, errors.New("not implemented")
}

func (r *memberRepoImpl) Create(c *model.MemberCreation) (*model.MemberInfo, error) {
	return nil, errors.New("not implemented")
}

func (r *memberRepoImpl) Update(u *model.MemberUpdate) error {
	return errors.New("not implemented")
}

func (r *memberRepoImpl) Delete(id string) error {
	return errors.New("not implemented")
}
