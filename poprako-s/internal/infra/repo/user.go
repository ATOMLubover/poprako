package repo

import (
	"context"
	"errors"

	iface "poprako-s/internal/domain/repo"

	"gorm.io/gorm"
)

type userRepoImpl struct {
	gdb *gorm.DB
}

func NewUserRepo(
	gdb *gorm.DB,
) iface.UserRepo {
	return &userRepoImpl{
		gdb: gdb,
	}
}

// 用于在事务上下文中获取已开启事务的 gdb
func NewUserRepoFromCtx(ctx context.Context) (iface.UserRepo, error) {
	gdb, err := ctx.Value(TxnKey).(*gorm.DB)
	if !err {
		return nil, errors.New("NewUserRepoFromCtx: 无法从上下文中获取事务数据库连接")
	}

	return &userRepoImpl{
		gdb: gdb,
	}, nil
}
