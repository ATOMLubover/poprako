package repo_infra

import (
	"context"
	"errors"
	"time"

	"poprako-s/internal/domain/model"
	iface "poprako-s/internal/domain/repo"
	entity "poprako-s/internal/infra/repo/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

	var row entity.UserCredsRow

	err := r.gdb.Table(entity.UserTable).
		Select("qq", "password_hash").
		Where("qq = ? AND deleted_at IS NULL", qq).
		First(&row).Error
	if err != nil {
		return nil, err
	}

	creds := entity.ToUserCreds(row)
	return &creds, nil
}

func (r *userRepoImpl) GetByID(id string) (*model.UserInfo, error) {

	var row entity.UserInfoRow

	err := r.gdb.Table(entity.UserTable).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&row).Error
	if err != nil {
		return nil, err
	}

	info := entity.ToUserInfo(row)
	return &info, nil
}

func (r *userRepoImpl) GetByQQ(qq string) (*model.UserInfo, error) {

	var row entity.UserInfoRow

	err := r.gdb.Table(entity.UserTable).
		Where("qq = ? AND deleted_at IS NULL", qq).
		First(&row).Error
	if err != nil {
		return nil, err
	}

	info := entity.ToUserInfo(row)
	return &info, nil
}

func (r *userRepoImpl) List(opt model.UserQueryOpt) ([]model.UserInfo, error) {
	db := r.gdb.Table(entity.UserTable).Where("deleted_at IS NULL")

	if opt.ID != nil {
		db = db.Where("id = ?", *opt.ID)
	}
	if opt.QQ != nil {
		db = db.Where("qq = ?", *opt.QQ)
	}

	var rows []entity.UserInfoRow

	if err := db.Order("created_at DESC").Find(&rows).Error; err != nil {
		return nil, err
	}

	items := make([]model.UserInfo, 0, len(rows))
	for _, row := range rows {
		items = append(items, entity.ToUserInfo(row))
	}

	return items, nil
}

func (r *userRepoImpl) Create(c *model.UserCreation) (*model.UserInfo, error) {
	now := time.Now()

	row := map[string]any{
		"id":                 c.ID,
		"name":               c.Name,
		"qq":                 c.QQ,
		"avatar_oss_key":     "",
		"is_avatar_uploaded": false,
		"password_hash":      c.PwdHash,
		"is_super_admin":     false,
		"last_login_at":      now,
		"created_at":         now,
		"updated_at":         now,
	}

	if err := r.gdb.Table(entity.UserTable).Create(row).Error; err != nil {
		return nil, err
	}

	return r.GetByID(c.ID)
}

func (r *userRepoImpl) Update(u *model.UserUpdate) error {
	return r.gdb.Table(entity.UserTable).
		Where("id = ? AND deleted_at IS NULL", u.ID).
		Updates(map[string]any{
			"name":       u.Name,
			"qq":         u.QQ,
			"updated_at": time.Now(),
		}).Error
}

func (r *userRepoImpl) Remove(id string) error {
	return r.gdb.Table(entity.UserTable).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{
			"deleted_at": time.Now(),
			"updated_at": time.Now(),
		}).Error
}

func (r *userRepoImpl) RefreshLastLogin(qq string, t time.Time) error {
	return r.gdb.Table(entity.UserTable).
		Where("qq = ? AND deleted_at IS NULL", qq).
		Updates(map[string]any{
			"last_login_at": t,
			"updated_at":    time.Now(),
		}).Error
}

func (r *userRepoImpl) PreFillAvatarOSSKey(id string, avatarOSSKey string) error {
	return r.gdb.Table(entity.UserTable).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{
			"avatar_oss_key":     avatarOSSKey,
			"is_avatar_uploaded": false,
			"updated_at":         time.Now(),
		}).Error
}

func (r *userRepoImpl) ConfirmAvatarUploaded(id string) error {
	return r.gdb.Table(entity.UserTable).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{
			"is_avatar_uploaded": true,
			"updated_at":         time.Now(),
		}).Error
}

func (r *userRepoImpl) GetOrCreateStats(userID string) (*model.UserStats, error) {

	var row entity.UserStatsRow

	err := r.gdb.Table(entity.UserStatsTable).
		Where("user_id = ?", userID).
		First(&row).Error
	if err == nil {
		stats := entity.ToUserStats(row)
		return &stats, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	createRow := map[string]any{
		"id":                        userID,
		"user_id":                   userID,
		"total_assignment_count":    0,
		"active_assignment_count":   0,
		"finished_assignment_count": 0,
		"created_at":                time.Now(),
		"updated_at":                time.Now(),
	}

	if err := r.gdb.Table(entity.UserStatsTable).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(createRow).Error; err != nil {
		return nil, err
	}

	err = r.gdb.Table(entity.UserStatsTable).
		Where("user_id = ?", userID).
		First(&row).Error
	if err != nil {
		return nil, err
	}

	stats := entity.ToUserStats(row)
	return &stats, nil
}

func (r *userRepoImpl) PatchStats(stats *model.UserStatsPatch) error {
	return r.gdb.Table(entity.UserStatsTable).
		Where("user_id = ?", stats.UserID).
		Updates(map[string]any{
			"total_assignment_count":    gorm.Expr("total_assignment_count + ?", stats.TotalAssignmentCountDelta),
			"active_assignment_count":   gorm.Expr("active_assignment_count + ?", stats.ActiveAssignmentCountDelta),
			"finished_assignment_count": gorm.Expr("finished_assignment_count + ?", stats.FinishedAssignmentCountDelta),
			"updated_at":                time.Now(),
		}).Error
}
