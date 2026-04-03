package repo_infra

import (
	"context"
	"errors"

	"poprako-s/internal/domain/model"
	iface "poprako-s/internal/domain/repo"

	"gorm.io/gorm"
)

type assignmentRepoImpl struct {
	gdb *gorm.DB
}

func NewAssignmentRepo(gdb *gorm.DB) iface.AssignmentRepo {
	return &assignmentRepoImpl{gdb: gdb}
}

func NewAssignmentRepoFromCx(cx context.Context) (iface.AssignmentRepo, error) {
	gdb, ok := cx.Value(txnKey).(*gorm.DB)
	if !ok {
		return nil, errors.New("[NewAssignmentRepoFromCx]: 无法从上下文中获取事务数据库连接")
	}

	return &assignmentRepoImpl{gdb: gdb}, nil
}

func (r *assignmentRepoImpl) FromTxnCx(cx context.Context) (iface.AssignmentRepo, error) {
	return NewAssignmentRepoFromCx(cx)
}

func (r *assignmentRepoImpl) GetByID(id string) (*model.AssignmentInfo, error) {
	return nil, errors.New("not implemented")
}

func (r *assignmentRepoImpl) Get(opt model.AssignmentQueryOpt) (*model.AssignmentInfo, error) {
	return nil, errors.New("not implemented")
}

func (r *assignmentRepoImpl) List(opt model.AssignmentQueryOpt) ([]model.AssignmentInfo, error) {
	return nil, errors.New("not implemented")
}

func (r *assignmentRepoImpl) Exist(opt model.AssignmentQueryOpt) (bool, error) {
	return false, errors.New("not implemented")
}

func (r *assignmentRepoImpl) Create(c *model.AssignmentCreation) (*model.AssignmentInfo, error) {
	return nil, errors.New("not implemented")
}

func (r *assignmentRepoImpl) Update(u *model.AssignmentUpdate) error {
	return errors.New("not implemented")
}

func (r *assignmentRepoImpl) Delete(id string) error {
	return errors.New("not implemented")
}
