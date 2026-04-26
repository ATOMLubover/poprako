package repo_infra

import (
	"context"
	"errors"
	"time"

	"poprako-s/internal/domain/model/aggr"
	repo_iface "poprako-s/internal/domain/repo"
	"poprako-s/internal/infra/repo/entity"

	"gorm.io/gorm"
)

// NOTE:
// - any repo impl should have its corresponding `Txn*Repo` function
//   and `New*Repo` constructor. No literal struct constructor.
// - must execute query with strongly typed structs, and select necessary
// 	 fields only.

type userRepoImpl struct {
	gdb *gorm.DB
}

func TxnUserRepo(cx context.Context) (repo_iface.UserRepo, repo_iface.RepoErr) {
	gdb := takeTxnGdb(cx)
	if gdb == nil {
		return nil, errors.New("[TxnUserRepo] no transaction context found for UserRepo")
	}

	return &userRepoImpl{gdb: gdb}, nil
}

func NewUserRepo(gdb *gorm.DB) repo_iface.UserRepo {
	return &userRepoImpl{gdb: gdb}
}

func (r *userRepoImpl) GetById(id string) (*aggr.User, repo_iface.RepoErr) {
	var row entity.UserRow

	err := r.gdb.
		Table(row.TableName()).
		Where("id = ?", id).
		First(&row).Error
	if err != nil {
		// Treat "record not found" as an error.
		return nil, err
	}

	return row.ToUserAggr(), nil
}

func (r *userRepoImpl) GetByQid(qid string) (*aggr.User, repo_iface.RepoErr) {
	var row entity.UserRow

	err := r.gdb.
		Table(row.TableName()).
		Where("qid = ?", qid).
		First(&row).Error
	if err != nil {
		// Treat "record not found" as an error.
		return nil, err
	}

	return row.ToUserAggr(), nil
}

func (r *userRepoImpl) GetCredsByQid(qid string) (*aggr.UserCreds, repo_iface.RepoErr) {
	var row entity.UserCredsRow

	err := r.gdb.
		Table(row.TableName()).
		Where("qid = ?", qid).
		First(&row).Error
	if err != nil {
		// Treat "record not found" as an error.
		return nil, err
	}

	return row.ToUserCredsAggr(), nil
}

func (r *userRepoImpl) Reg(reg *aggr.UserReg) (*aggr.User, repo_iface.RepoErr) {
	regRow := entity.NewUserRegRowFromAggr(reg)

	err := r.gdb.
		Table(regRow.TableName()).
		Create(regRow).Error
	if err != nil {
		return nil, err
	}

	return r.GetById(regRow.Id)
}

func (r *userRepoImpl) Refresh(id string, activeAt time.Time) repo_iface.RepoErr {
	return r.gdb.
		Table(entity.USER_TABLE).
		Where("id = ?", id).
		Update("last_active_at", activeAt).Error
}

func (r *userRepoImpl) PrefillAvatarKey(id string, key string) repo_iface.RepoErr {
	return r.gdb.
		Table(entity.USER_TABLE).
		Where("id = ?", id).
		Update("avatar_key", key).Error
}

func (r *userRepoImpl) MarkAvatarUploaded(id string) repo_iface.RepoErr {
	return r.gdb.
		Table(entity.USER_TABLE).
		Where("id = ?", id).
		Update("avatar_uploaded", true).Error
}
