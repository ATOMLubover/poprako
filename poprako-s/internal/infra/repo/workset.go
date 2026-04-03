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

type worksetRepoImpl struct {
	gdb *gorm.DB
}

func NewWorksetRepo(gdb *gorm.DB) iface.WorksetRepo {
	return &worksetRepoImpl{gdb: gdb}
}

func NewWorksetRepoFromCx(cx context.Context) (iface.WorksetRepo, error) {
	gdb, ok := cx.Value(txnKey).(*gorm.DB)
	if !ok {
		return nil, errors.New("[NewWorksetRepoFromCx]: 无法从上下文中获取事务数据库连接")
	}

	return &worksetRepoImpl{gdb: gdb}, nil
}

func (r *worksetRepoImpl) FromTxnCx(cx context.Context) (iface.WorksetRepo, error) {
	return NewWorksetRepoFromCx(cx)
}

func (r *worksetRepoImpl) GetByID(id string) (*model.WorksetInfo, error) {

	var row entity.WorksetInfoRow

	err := r.gdb.Table(entity.WorksetTable).Where("id = ?", id).First(&row).Error
	if err != nil {
		return nil, err
	}

	info := entity.ToWorksetInfo(row)
	return &info, nil
}

func (r *worksetRepoImpl) List(opt model.WorksetQueryOpt) ([]model.WorksetInfo, error) {
	db := r.gdb.Table(entity.WorksetTable)

	if opt.ID != nil {
		db = db.Where("id = ?", *opt.ID)
	}
	if opt.TeamID != nil {
		db = db.Where("team_id = ?", *opt.TeamID)
	}

	var rows []entity.WorksetInfoRow

	if err := db.Order("index ASC").Find(&rows).Error; err != nil {
		return nil, err
	}

	items := make([]model.WorksetInfo, 0, len(rows))
	for _, row := range rows {
		items = append(items, entity.ToWorksetInfo(row))
	}

	return items, nil
}

func (r *worksetRepoImpl) Count(opt model.WorksetQueryOpt) (int64, error) {
	db := r.gdb.Table(entity.WorksetTable)

	if opt.ID != nil {
		db = db.Where("id = ?", *opt.ID)
	}
	if opt.TeamID != nil {
		db = db.Where("team_id = ?", *opt.TeamID)
	}

	var n int64

	err := db.Count(&n).Error
	return n, err
}

func (r *worksetRepoImpl) Create(c *model.WorksetCreation) (*model.WorksetInfo, error) {
	now := time.Now()
	row := map[string]any{
		"id":          c.ID,
		"team_id":     c.TeamID,
		"index":       c.Index,
		"name":        c.Name,
		"description": c.Description,
		"comic_count": 0,
		"created_at":  now,
		"updated_at":  now,
	}

	if err := r.gdb.Table(entity.WorksetTable).Create(row).Error; err != nil {
		return nil, err
	}

	return r.GetByID(c.ID)
}

func (r *worksetRepoImpl) Update(u *model.WorksetUpdate) error {
	updates := map[string]any{
		"name":       u.Name,
		"updated_at": time.Now(),
	}
	if u.Description != nil {
		updates["description"] = *u.Description
	}

	return r.gdb.Table(entity.WorksetTable).
		Where("id = ?", u.ID).
		Updates(updates).Error
}

func (r *worksetRepoImpl) UpdateComicCount(id string, delta int) error {
	return r.gdb.Table(entity.WorksetTable).
		Where("id = ?", id).
		Updates(map[string]any{
			"comic_count": gorm.Expr("comic_count + ?", delta),
			"updated_at":  time.Now(),
		}).Error
}

func (r *worksetRepoImpl) Delete(id string) error {
	return r.gdb.Table(entity.WorksetTable).Where("id = ?", id).Delete(nil).Error
}
