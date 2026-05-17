package repo_infra

import (
	"time"

	"poprako-s/internal/domain/model/aggr"
	repo_iface "poprako-s/internal/domain/repo"
	"poprako-s/internal/infra/repo/entity"

	"gorm.io/gorm"
)

// NOTE:
// - any repo impl should have its `New*Repo` constructor. No literal struct constructor.
// - must execute query with strongly typed structs, and select necessary
// 	 fields only.

// `userRepoImpl` is the GORM-backed implementation of `repo_iface.UserRepo`
type userRepoImpl struct {
	gdb *gorm.DB
}

// `NewUserRepo` creates a non-transaction-scoped `UserRepo`
func NewUserRepo(gdb *gorm.DB) repo_iface.UserRepo {
	return &userRepoImpl{gdb: gdb}
}

// `GetById` retrieves one `User` aggregate by id
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

// `GetByQid` retrieves one `User` aggregate by QID
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

// `GetCredsByQid` retrieves the credential row for the user with the given QID
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

// `Register` inserts a new user row from `UserReg` and returns the created aggregate
func (r *userRepoImpl) Register(reg *aggr.UserReg) (*aggr.User, repo_iface.RepoErr) {
	regRow := entity.NewUserRegRowFromAggr(reg)

	err := r.gdb.
		Table(regRow.TableName()).
		Create(regRow).Error
	if err != nil {
		return nil, err
	}

	return r.GetById(regRow.Id)
}

// `Update` applies put-style mutable fields to one user row.
func (r *userRepoImpl) Update(upd *aggr.UserUpd) repo_iface.RepoErr {
	updRe := r.gdb.
		Table(entity.USER_TABLE).
		Where("id = ?", upd.Id).
		Updates(map[string]any{
			"nickname":   upd.Name,
			"qid":        upd.Qid,
			"updated_at": time.Now(),
		})
	if updRe.Error != nil {
		return updRe.Error
	}

	if updRe.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// `Refresh` updates the `last_active_at` timestamp for the given user id
func (r *userRepoImpl) Refresh(id string, activeAt time.Time) repo_iface.RepoErr {
	updRe := r.gdb.
		Table(entity.USER_TABLE).
		Where("id = ?", id).
		Updates(map[string]any{
			"last_active_at": activeAt,
			"updated_at":     time.Now(),
		})
	if updRe.Error != nil {
		return updRe.Error
	}
	if updRe.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// `PrefillAvatarKey` writes the OSS object key for the user avatar before the upload begins
func (r *userRepoImpl) PrefillAvatarKey(id string, key string) repo_iface.RepoErr {
	return r.gdb.
		Table(entity.USER_TABLE).
		Where("id = ?", id).
		Update("avatar_key", key).Error
}

// `MarkAvatarUploaded` sets `avatar_uploaded` to true for the given user id
func (r *userRepoImpl) MarkAvatarUploaded(id string) repo_iface.RepoErr {
	return r.gdb.
		Table(entity.USER_TABLE).
		Where("id = ?", id).
		Update("avatar_uploaded", true).Error
}
