package repo_infra

import (
	"context"
	"errors"
	"time"

	"poprako-s/internal/domain/model"
	iface "poprako-s/internal/domain/repo"
	entity "poprako-s/internal/infra/repo/entity"

	"gorm.io/gorm"
)

type teamRepoImpl struct {
	gdb *gorm.DB
}

func NewTeamRepo(gdb *gorm.DB) iface.TeamRepo {
	return &teamRepoImpl{gdb: gdb}
}

func NewTeamRepoFromCx(cx context.Context) (iface.TeamRepo, error) {
	gdb, ok := cx.Value(txnKey).(*gorm.DB)
	if !ok {
		return nil, errors.New("[NewTeamRepoFromCx]: 无法从上下文中获取事务数据库连接")
	}

	return &teamRepoImpl{gdb: gdb}, nil
}

func (r *teamRepoImpl) FromTxnCx(cx context.Context) (iface.TeamRepo, error) {
	return NewTeamRepoFromCx(cx)
}

func (r *teamRepoImpl) GetByID(id string) (*model.TeamInfo, error) {

	var row entity.TeamInfoRow

	err := r.gdb.Table(entity.TeamTable).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&row).Error
	if err != nil {
		return nil, err
	}

	info := entity.ToTeamInfo(row)
	return &info, nil
}

func (r *teamRepoImpl) List(opt model.TeamQueryOpt) ([]model.TeamInfo, error) {
	db := r.gdb.Table(entity.TeamTable).Where("deleted_at IS NULL")
	if opt.ID != nil {
		db = db.Where("id = ?", *opt.ID)
	}

	var rows []entity.TeamInfoRow

	if err := db.Order("created_at DESC").Find(&rows).Error; err != nil {
		return nil, err
	}

	items := make([]model.TeamInfo, 0, len(rows))
	for _, row := range rows {
		items = append(items, entity.ToTeamInfo(row))
	}
	return items, nil
}

func (r *teamRepoImpl) Create(c *model.TeamCreation) (*model.TeamInfo, error) {
	now := time.Now()
	row := map[string]any{
		"id":                 c.ID,
		"name":               c.Name,
		"description":        c.Description,
		"avatar_oss_key":     "",
		"is_avatar_uploaded": false,
		"created_at":         now,
		"updated_at":         now,
	}

	if err := r.gdb.Table(entity.TeamTable).Create(row).Error; err != nil {
		return nil, err
	}

	return r.GetByID(c.ID)
}

func (r *teamRepoImpl) Update(u *model.TeamUpdate) error {
	return r.gdb.Table(entity.TeamTable).
		Where("id = ? AND deleted_at IS NULL", u.ID).
		Updates(map[string]any{
			"name":        u.Name,
			"description": u.Description,
			"updated_at":  time.Now(),
		}).Error
}

func (r *teamRepoImpl) Delete(id string) error {
	return r.gdb.Table(entity.TeamTable).Where("id = ?", id).Delete(nil).Error
}

func (r *teamRepoImpl) PreFillAvatarOSSKey(id string, avatarOSSKey string) error {
	return r.gdb.Table(entity.TeamTable).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{
			"avatar_oss_key":     avatarOSSKey,
			"is_avatar_uploaded": false,
			"updated_at":         time.Now(),
		}).Error
}

func (r *teamRepoImpl) ConfirmAvatarUploaded(id string) error {
	return r.gdb.Table(entity.TeamTable).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{
			"is_avatar_uploaded": true,
			"updated_at":         time.Now(),
		}).Error
}
