package repo_infra

import (
	"context"
	"errors"

	"poprako-s/internal/domain/model"
	iface "poprako-s/internal/domain/repo"

	"gorm.io/gorm"
)

type invitationRepoImpl struct {
	gdb *gorm.DB
}

func NewInvitationRepo(
	gdb *gorm.DB,
) iface.InvitationRepo {
	return &invitationRepoImpl{
		gdb: gdb,
	}
}

func NewInvitationRepoFromCx(cx context.Context) (iface.InvitationRepo, error) {
	gdb, err := cx.Value(txnKey).(*gorm.DB)
	if !err {
		return nil, errors.New("[NewInvitationRepoFromCx]: 无法从上下文中获取事务数据库连接")
	}

	return &invitationRepoImpl{
		gdb: gdb,
	}, nil
}

func (r *invitationRepoImpl) GetByInviteeQQ(id string) (*model.InvitationInfo, error) {
	return nil, errors.New("not implemented")
}

func (r *invitationRepoImpl) List(opt model.InvitationQueryOpt) ([]model.InvitationInfo, error) {
	return nil, errors.New("not implemented")
}

func (r *invitationRepoImpl) Create(c *model.InvitationCreation) (*model.InvitationInfo, error) {
	return nil, errors.New("not implemented")
}

func (r *invitationRepoImpl) Update(u *model.InvitationUpdate) error {
	return errors.New("not implemented")
}

func (r *invitationRepoImpl) Invalidate(id string) error {
	return errors.New("not implemented")
}

func (r *invitationRepoImpl) Delete(id string) error {
	return errors.New("not implemented")
}
