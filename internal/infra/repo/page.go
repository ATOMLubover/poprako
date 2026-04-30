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

type pageRepoImpl struct {
	gdb *gorm.DB
}

func NewPageRepo(gdb *gorm.DB) iface.PageRepo {
	return &pageRepoImpl{gdb: gdb}
}

func NewPageRepoFromCx(cx context.Context) (iface.PageRepo, error) {
	gdb, ok := cx.Value(txnKey).(*gorm.DB)
	if !ok {
		return nil, errors.New("[NewPageRepoFromCx]: 无法从上下文中获取事务数据库连接")
	}

	return &pageRepoImpl{gdb: gdb}, nil
}

func (r *pageRepoImpl) FromTxnCx(cx context.Context) (iface.PageRepo, error) {
	return NewPageRepoFromCx(cx)
}

func (r *pageRepoImpl) GetByID(id string) (*model.PageInfo, error) {

	var row entity.PageInfoRow

	err := r.gdb.Table(entity.PageTable).Where("id = ?", id).First(&row).Error
	if err != nil {
		return nil, err
	}

	info := entity.ToPageInfo(row)
	return &info, nil
}

func (r *pageRepoImpl) List(opt model.PageQueryOpt) ([]model.PageInfo, error) {
	db := r.gdb.Table(entity.PageTable)

	if opt.ID != nil {
		db = db.Where("id = ?", *opt.ID)
	}
	if opt.ChapterID != nil {
		db = db.Where("chapter_id = ?", *opt.ChapterID)
	}

	var rows []entity.PageInfoRow

	if err := db.Order("index ASC").Find(&rows).Error; err != nil {
		return nil, err
	}

	items := make([]model.PageInfo, 0, len(rows))
	for _, row := range rows {
		items = append(items, entity.ToPageInfo(row))
	}

	return items, nil
}

func (r *pageRepoImpl) LockByID(id string) error {
	return r.gdb.Table(entity.PageTable).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", id).
		Select("id").
		First(&entity.PageInfoRow{}).Error
}

func (r *pageRepoImpl) CreateBatch(pages []*model.PageCreation) error {
	if len(pages) == 0 {
		return nil
	}

	now := time.Now()
	rows := make([]map[string]any, 0, len(pages))
	for _, page := range pages {
		rows = append(rows, map[string]any{
			"id":                    page.ID,
			"chapter_id":            page.ChapterID,
			"index":                 page.Index,
			"oss_key":               page.OSSKey,
			"uploaded":              false,
			"creator_id":            page.CreatorID,
			"total_unit_count":      0,
			"translated_unit_count": 0,
			"proofread_unit_count":  0,
			"created_at":            now,
			"updated_at":            now,
		})
	}

	return r.gdb.Table(entity.PageTable).Create(rows).Error
}

func (r *pageRepoImpl) Update(u *model.PageUpdate) error {
	return r.gdb.Table(entity.PageTable).
		Where("id = ?", u.ID).
		Updates(map[string]any{
			"index":      u.Index,
			"oss_key":    u.OSSKey,
			"uploaded":   u.IsUploaded,
			"updated_at": time.Now(),
		}).Error
}

func (r *pageRepoImpl) UpdateStats(id string, totalDelta, translatedDelta, proofreadDelta int) error {
	return r.gdb.Table(entity.PageTable).
		Where("id = ?", id).
		Updates(map[string]any{
			"total_unit_count":      gorm.Expr("total_unit_count + ?", totalDelta),
			"translated_unit_count": gorm.Expr("translated_unit_count + ?", translatedDelta),
			"proofread_unit_count":  gorm.Expr("proofread_unit_count + ?", proofreadDelta),
			"updated_at":            time.Now(),
		}).Error
}

func (r *pageRepoImpl) Delete(id string) error {
	return r.gdb.Table(entity.PageTable).Where("id = ?", id).Delete(nil).Error
}

func (r *pageRepoImpl) DeleteBatch(ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	return r.gdb.Table(entity.PageTable).Where("id IN ?", ids).Delete(nil).Error
}
