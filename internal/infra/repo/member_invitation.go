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

type memberInvitationRepoImpl struct {
	gdb *gorm.DB
}

func NewMemberInvitationRepo(
	gdb *gorm.DB,
) iface.MemberInvitationRepo {
	return &memberInvitationRepoImpl{
		gdb: gdb,
	}
}

func NewMemberInvitationRepoFromCx(cx context.Context) (iface.MemberInvitationRepo, error) {
	gdb, err := cx.Value(txnKey).(*gorm.DB)
	if !err {
		return nil, errors.New("[NewInvitationRepoFromCx]: 无法从上下文中获取事务数据库连接")
	}

	return &memberInvitationRepoImpl{
		gdb: gdb,
	}, nil
}

func (r *memberInvitationRepoImpl) FromTxnCx(cx context.Context) (iface.MemberInvitationRepo, error) {
	return NewMemberInvitationRepoFromCx(cx)
}

func (r *memberInvitationRepoImpl) GetByInviteeQQ(id string) (*model.MemberInvitationInfo, error) {
	var row entity.MemberInvitationInfoRow

	err := r.gdb.Table(entity.MemberInvitationTable).
		Where("invitee_qq = ?", id).
		Order("created_at DESC").
		First(&row).Error
	if err != nil {
		return nil, err
	}

	info := entity.ToMemberInvitationInfo(row)
	return &info, nil
}

func (r *memberInvitationRepoImpl) List(opt model.MemberInvitationQueryOpt) ([]model.MemberInvitationInfo, error) {
	db := r.gdb.Table(entity.MemberInvitationTable)

	if opt.TeamID != nil {
		db = db.Where("team_id = ?", *opt.TeamID)
	}
	if opt.InvitationCode != nil {
		db = db.Where("invitation_code = ?", *opt.InvitationCode)
	}
	if opt.Pending != nil {
		db = db.Where("pending = ?", *opt.Pending)
	}

	var rows []entity.MemberInvitationInfoRow

	if err := db.Order("created_at DESC").Find(&rows).Error; err != nil {
		return nil, err
	}

	items := make([]model.MemberInvitationInfo, 0, len(rows))
	for _, row := range rows {
		items = append(items, entity.ToMemberInvitationInfo(row))
	}

	return items, nil
}

func (r *memberInvitationRepoImpl) Create(c *model.MemberInvitationCreation) (*model.MemberInvitationInfo, error) {
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

	if err := r.gdb.Table(entity.MemberInvitationTable).Create(row).Error; err != nil {
		return nil, err
	}

	var created entity.MemberInvitationInfoRow

	err := r.gdb.Table(entity.MemberInvitationTable).Where("id = ?", c.ID).First(&created).Error
	if err != nil {
		return nil, err
	}

	info := entity.ToMemberInvitationInfo(created)
	return &info, nil
}

func (r *memberInvitationRepoImpl) Update(u *model.MemberInvitationUpdate) error {
	return r.gdb.Table(entity.MemberInvitationTable).
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

func (r *memberInvitationRepoImpl) Invalidate(id string) error {
	return r.gdb.Table(entity.MemberInvitationTable).
		Where("id = ?", id).
		Updates(map[string]any{
			"pending":    false,
			"updated_at": time.Now(),
		}).Error
}

func (r *memberInvitationRepoImpl) Delete(id string) error {
	return r.gdb.Table(entity.MemberInvitationTable).Where("id = ?", id).Delete(nil).Error
}
