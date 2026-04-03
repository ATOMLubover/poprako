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

type invitationRepoImpl struct {
	gdb *gorm.DB
}

func NewInvitationRepo(
	gdb *gorm.DB,
) iface.InvitationRepo {
	return &invitationRepoImpl{
		gdb: gdb,
	}
}

func NewInvitationRepoFromCx(cx context.Context) (iface.InvitationRepo, error) {
	gdb, err := cx.Value(txnKey).(*gorm.DB)
	if !err {
		return nil, errors.New("[NewInvitationRepoFromCx]: 无法从上下文中获取事务数据库连接")
	}

	return &invitationRepoImpl{
		gdb: gdb,
	}, nil
}

func (r *invitationRepoImpl) FromTxnCx(cx context.Context) (iface.InvitationRepo, error) {
	return NewInvitationRepoFromCx(cx)
}

func (r *invitationRepoImpl) GetByInviteeQQ(id string) (*model.InvitationInfo, error) {

	var row entity.InvitationInfoRow

	err := r.gdb.Table(entity.InvitationTable).
		Where("invitee_qq = ?", id).
		Order("created_at DESC").
		First(&row).Error
	if err != nil {
		return nil, err
	}

	info := entity.ToInvitationInfo(row)
	return &info, nil
}

func (r *invitationRepoImpl) List(opt model.InvitationQueryOpt) ([]model.InvitationInfo, error) {
	db := r.gdb.Table(entity.InvitationTable)

	if opt.TeamID != nil {
		db = db.Where("team_id = ?", *opt.TeamID)
	}
	if opt.InvitationCode != nil {
		db = db.Where("invitation_code = ?", *opt.InvitationCode)
	}
	db = db.Where("pending = ?", opt.Pending)

	var rows []entity.InvitationInfoRow

	if err := db.Order("created_at DESC").Find(&rows).Error; err != nil {
		return nil, err
	}

	items := make([]model.InvitationInfo, 0, len(rows))
	for _, row := range rows {
		items = append(items, entity.ToInvitationInfo(row))
	}

	return items, nil
}

func (r *invitationRepoImpl) Create(c *model.InvitationCreation) (*model.InvitationInfo, error) {
	now := time.Now()
	row := map[string]any{
		"id":                 c.ID,
		"invitor_id":         c.InvitorID,
		"team_id":            c.TargetTeamID,
		"invitee_qq":         c.InviteeQQ,
		"invitation_code":    c.InvitationCode,
		"to_be_raw_provider": c.ToBeRawProvider,
		"to_be_translator":   c.ToBeTranslator,
		"to_be_proofreader":  c.ToBeProofreader,
		"to_be_typesetter":   c.ToBeTypesetter,
		"to_be_reviewer":     c.ToBeReviewer,
		"to_be_publisher":    c.ToBePublisher,
		"to_be_admin":        c.ToBeAdmin,
		"pending":            true,
		"created_at":         now,
		"updated_at":         now,
	}

	if err := r.gdb.Table(entity.InvitationTable).Create(row).Error; err != nil {
		return nil, err
	}

	var created entity.InvitationInfoRow

	err := r.gdb.Table(entity.InvitationTable).Where("id = ?", c.ID).First(&created).Error
	if err != nil {
		return nil, err
	}

	info := entity.ToInvitationInfo(created)
	return &info, nil
}

func (r *invitationRepoImpl) Update(u *model.InvitationUpdate) error {
	return r.gdb.Table(entity.InvitationTable).
		Where("id = ?", u.ID).
		Updates(map[string]any{
			"to_be_raw_provider": u.ToBeRawProvider,
			"to_be_translator":   u.ToBeTranslator,
			"to_be_proofreader":  u.ToBeProofreader,
			"to_be_typesetter":   u.ToBeTypesetter,
			"to_be_reviewer":     u.ToBeReviewer,
			"to_be_publisher":    u.ToBePublisher,
			"to_be_admin":        u.ToBeAdmin,
			"updated_at":         time.Now(),
		}).Error
}

func (r *invitationRepoImpl) Invalidate(id string) error {
	return r.gdb.Table(entity.InvitationTable).
		Where("id = ?", id).
		Updates(map[string]any{
			"pending":    false,
			"updated_at": time.Now(),
		}).Error
}

func (r *invitationRepoImpl) Delete(id string) error {
	return r.gdb.Table(entity.InvitationTable).Where("id = ?", id).Delete(nil).Error
}
