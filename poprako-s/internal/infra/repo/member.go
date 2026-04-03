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

type memberRepoImpl struct {
	gdb *gorm.DB
}

func NewMemberRepo(gdb *gorm.DB) iface.MemberRepo {
	return &memberRepoImpl{gdb: gdb}
}

func NewMemberRepoFromCx(cx context.Context) (iface.MemberRepo, error) {
	gdb, ok := cx.Value(txnKey).(*gorm.DB)
	if !ok {
		return nil, errors.New("[NewMemberRepoFromCx]: 无法从上下文中获取事务数据库连接")
	}

	return &memberRepoImpl{gdb: gdb}, nil
}

func (r *memberRepoImpl) FromTxnCx(cx context.Context) (iface.MemberRepo, error) {
	return NewMemberRepoFromCx(cx)
}

func (r *memberRepoImpl) GetByID(id string) (*model.MemberInfo, error) {

	var row entity.MemberInfoRow

	err := r.gdb.Table(entity.MemberTable).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&row).Error
	if err != nil {
		return nil, err
	}

	info := entity.ToMemberInfo(row)
	return &info, nil
}

func (r *memberRepoImpl) Get(opt model.MemberQueryOpt) (*model.MemberInfo, error) {
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

func (r *memberRepoImpl) List(opt model.MemberQueryOpt) ([]model.MemberInfo, error) {
	db := r.gdb.Table(entity.MemberTable).Where("deleted_at IS NULL")

	if opt.ID != nil {
		db = db.Where("id = ?", *opt.ID)
	}
	if opt.UserID != nil {
		db = db.Where("user_id = ?", *opt.UserID)
	}
	if opt.TeamID != nil {
		db = db.Where("team_id = ?", *opt.TeamID)
	}

	var rows []entity.MemberInfoRow

	if err := db.Order("created_at ASC").Find(&rows).Error; err != nil {
		return nil, err
	}

	items := make([]model.MemberInfo, 0, len(rows))
	for _, row := range rows {
		items = append(items, entity.ToMemberInfo(row))
	}

	return items, nil
}

func (r *memberRepoImpl) Exist(opt model.MemberQueryOpt) (bool, error) {
	db := r.gdb.Table(entity.MemberTable).Where("deleted_at IS NULL")

	if opt.ID != nil {
		db = db.Where("id = ?", *opt.ID)
	}
	if opt.UserID != nil {
		db = db.Where("user_id = ?", *opt.UserID)
	}
	if opt.TeamID != nil {
		db = db.Where("team_id = ?", *opt.TeamID)
	}

	var n int64

	err := db.Count(&n).Error
	return n > 0, err
}

func (r *memberRepoImpl) Create(c *model.MemberCreation) (*model.MemberInfo, error) {
	now := time.Now()
	row := map[string]any{
		"id":         c.ID,
		"user_id":    c.UserID,
		"team_id":    c.TeamID,
		"created_at": now,
		"updated_at": now,
	}

	if c.ToBeRawProvider {
		row["assigned_raw_provider_at"] = now
	}
	if c.ToBeTranslator {
		row["assigned_translator_at"] = now
	}
	if c.ToBeProofreader {
		row["assigned_proofreader_at"] = now
	}
	if c.ToBeTypesetter {
		row["assigned_typesetter_at"] = now
	}
	if c.ToBeReviewer {
		row["assigned_reviewer_at"] = now
	}
	if c.ToBePublisher {
		row["assigned_publisher_at"] = now
	}
	if c.ToBeAdmin {
		row["assigned_admin_at"] = now
	}

	if err := r.gdb.Table(entity.MemberTable).Create(row).Error; err != nil {
		return nil, err
	}

	return r.GetByID(c.ID)
}

func (r *memberRepoImpl) Update(u *model.MemberUpdate) error {
	return r.gdb.Table(entity.MemberTable).
		Where("id = ? AND deleted_at IS NULL", u.ID).
		Updates(map[string]any{
			"assigned_raw_provider_at": u.AssignedRawProviderAt,
			"assigned_translator_at":   u.AssignedTranslatorAt,
			"assigned_proofreader_at":  u.AssignedProofreaderAt,
			"assigned_typesetter_at":   u.AssignedTypesetterAt,
			"assigned_reviewer_at":     u.AssignedReviewerAt,
			"assigned_publisher_at":    u.AssignedPublisherAt,
			"assigned_admin_at":        u.AssignedAdminAt,
			"updated_at":               time.Now(),
		}).Error
}

func (r *memberRepoImpl) Delete(id string) error {
	return r.gdb.Table(entity.MemberTable).Where("id = ?", id).Delete(nil).Error
}
