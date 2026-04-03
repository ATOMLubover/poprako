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

type assignmentRepoImpl struct {
	gdb *gorm.DB
}

func NewAssignmentRepo(gdb *gorm.DB) iface.AssignmentRepo {
	return &assignmentRepoImpl{gdb: gdb}
}

func NewAssignmentRepoFromCx(cx context.Context) (iface.AssignmentRepo, error) {
	gdb, ok := cx.Value(txnKey).(*gorm.DB)
	if !ok {
		return nil, errors.New("[NewAssignmentRepoFromCx]: 无法从上下文中获取事务数据库连接")
	}

	return &assignmentRepoImpl{gdb: gdb}, nil
}

func (r *assignmentRepoImpl) FromTxnCx(cx context.Context) (iface.AssignmentRepo, error) {
	return NewAssignmentRepoFromCx(cx)
}

func (r *assignmentRepoImpl) GetByID(id string) (*model.AssignmentInfo, error) {
	var row entity.AssignmentInfoRow
	err := r.gdb.Table(entity.AssignmentTable).Where("id = ?", id).First(&row).Error
	if err != nil {
		return nil, err
	}

	info := entity.ToAssignmentInfo(row)
	return &info, nil
}

func (r *assignmentRepoImpl) Get(opt model.AssignmentQueryOpt) (*model.AssignmentInfo, error) {
	items, err := r.List(opt)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	item := items[0]
	return &item, nil
}

func (r *assignmentRepoImpl) List(opt model.AssignmentQueryOpt) ([]model.AssignmentInfo, error) {
	db := r.gdb.Table(entity.AssignmentTable)

	if opt.ChapterID != nil {
		db = db.Where("chapter_id = ?", *opt.ChapterID)
	}
	if opt.UserID != nil {
		db = db.Where("user_id = ?", *opt.UserID)
	}

	var rows []entity.AssignmentInfoRow
	if err := db.Order("created_at ASC").Find(&rows).Error; err != nil {
		return nil, err
	}

	items := make([]model.AssignmentInfo, 0, len(rows))
	for _, row := range rows {
		items = append(items, entity.ToAssignmentInfo(row))
	}

	return items, nil
}

func (r *assignmentRepoImpl) Exist(opt model.AssignmentQueryOpt) (bool, error) {
	db := r.gdb.Table(entity.AssignmentTable)

	if opt.ChapterID != nil {
		db = db.Where("chapter_id = ?", *opt.ChapterID)
	}
	if opt.UserID != nil {
		db = db.Where("user_id = ?", *opt.UserID)
	}

	var n int64
	err := db.Count(&n).Error
	return n > 0, err
}

func (r *assignmentRepoImpl) Create(c *model.AssignmentCreation) (*model.AssignmentInfo, error) {
	now := time.Now()
	row := map[string]any{
		"id":                       c.ID,
		"chapter_id":               c.ChapterID,
		"user_id":                  c.UserID,
		"assigned_raw_provider_at": c.AssignedRawProviderAt,
		"assigned_translator_at":   c.AssignedTranslatorAt,
		"assigned_proofreader_at":  c.AssignedProofreaderAt,
		"assigned_typesetter_at":   c.AssignedTypesetterAt,
		"assigned_redrawer_at":     c.AssignedRedrawerAt,
		"assigned_reviewer_at":     c.AssignedReviewerAt,
		"assigned_publisher_at":    c.AssignedPublisherAt,
		"created_at":               now,
		"updated_at":               now,
	}

	if err := r.gdb.Table(entity.AssignmentTable).Create(row).Error; err != nil {
		return nil, err
	}

	return r.GetByID(c.ID)
}

func (r *assignmentRepoImpl) Update(u *model.AssignmentUpdate) error {
	return r.gdb.Table(entity.AssignmentTable).
		Where("id = ?", u.ID).
		Updates(map[string]any{
			"assigned_raw_provider_at": u.AssignedRawProviderAt,
			"assigned_translator_at":   u.AssignedTranslatorAt,
			"assigned_proofreader_at":  u.AssignedProofreaderAt,
			"assigned_typesetter_at":   u.AssignedTypesetterAt,
			"assigned_redrawer_at":     u.AssignedRedrawerAt,
			"assigned_reviewer_at":     u.AssignedReviewerAt,
			"assigned_publisher_at":    u.AssignedPublisherAt,
			"updated_at":               time.Now(),
		}).Error
}

func (r *assignmentRepoImpl) Delete(id string) error {
	return r.gdb.Table(entity.AssignmentTable).Where("id = ?", id).Delete(nil).Error
}
