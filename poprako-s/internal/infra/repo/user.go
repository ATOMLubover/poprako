package repo_infra

import (
	"context"
	"errors"
	"time"

	"poprako-s/internal/domain/model"
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
func NewUserRepoFromCx(cx context.Context) (iface.UserRepo, error) {
	gdb, err := cx.Value(txnKey).(*gorm.DB)
	if !err {
		return nil, errors.New("[NewUserRepoFromCx]: 无法从上下文中获取事务数据库连接")
	}

	return &userRepoImpl{
		gdb: gdb,
	}, nil
}

func (r *userRepoImpl) FromTxnCx(cx context.Context) (iface.UserRepo, error) {
	return NewUserRepoFromCx(cx)
}

func (r *userRepoImpl) GetCredsByQQ(qq string) (*model.UserCreds, error) {
	return nil, errors.New("not implemented")
}

func (r *userRepoImpl) GetByID(id string) (*model.UserInfo, error) {
	return nil, errors.New("not implemented")
}

func (r *userRepoImpl) GetByQQ(qq string) (*model.UserInfo, error) {
	return nil, errors.New("not implemented")
}

func (r *userRepoImpl) List(opt model.UserQueryOpt) ([]model.UserInfo, error) {
	return nil, errors.New("not implemented")
}

func (r *userRepoImpl) Create(c *model.UserCreation) (*model.UserInfo, error) {
	return nil, errors.New("not implemented")
}

func (r *userRepoImpl) Update(u *model.UserUpdate) error {
	return errors.New("not implemented")
}

func (r *userRepoImpl) Remove(id string) error {
	return errors.New("not implemented")
}

func (r *userRepoImpl) RefreshLastLogin(qq string, t time.Time) error {
	return errors.New("not implemented")
}

func (r *userRepoImpl) PreFillAvatarOSSKey(id string, avatarOSSKey string) error {
	return errors.New("not implemented")
}

func (r *userRepoImpl) ConfirmAvatarUploaded(id string) error {
	return errors.New("not implemented")
}

func (r *userRepoImpl) GetOrCreateStats(userID string) (*model.UserStats, error) {
	return nil, errors.New("not implemented")
}

func (r *userRepoImpl) PatchStats(stats *model.UserStatsPatch) error {
	return errors.New("not implemented")
}
