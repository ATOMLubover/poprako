package repo_infra

import (
	"context"
	"errors"

	iface "poprako-s/internal/domain/repo"

	"gorm.io/gorm"
)

type invitationRepoImpl struct {
	gdb *gorm.DB
}

func NewInvitationRepo(
	gdb *gorm.DB,
) iface.InvitationRepo {
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
