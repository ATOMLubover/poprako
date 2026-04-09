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

type chapterInvitationRepoImpl struct {
	gdb *gorm.DB
}

func NewChapterInvitationRepo(gdb *gorm.DB) iface.ChapterInvitationRepo {
	return &chapterInvitationRepoImpl{gdb: gdb}
}

func NewChapterInvitationRepoFromCx(cx context.Context) (iface.ChapterInvitationRepo, error) {
	gdb, ok := cx.Value(txnKey).(*gorm.DB)
	if !ok {
		return nil, errors.New("[NewChapterInvitationRepoFromCx]: 无法从上下文中获取事务数据库连接")
	}

	return &chapterInvitationRepoImpl{gdb: gdb}, nil
}

func (r *chapterInvitationRepoImpl) FromTxnCx(cx context.Context) (iface.ChapterInvitationRepo, error) {
	return NewChapterInvitationRepoFromCx(cx)
}

func (r *chapterInvitationRepoImpl) List(opt model.ChapterInvitationQueryOpt) ([]model.ChapterInvitationInfo, error) {
	db := r.gdb.Table(entity.ChapterInvitationTable)

	if opt.ChapterID != nil {
		db = db.Where("chapter_id = ?", *opt.ChapterID)
	}
	if opt.InvitationCode != nil {
		db = db.Where("invitation_code = ?", *opt.InvitationCode)
	}
	if opt.InviteeQQ != nil {
		db = db.Where("invitee_qq = ?", *opt.InviteeQQ)
	}
	if opt.OnlyPendingTrue {
		db = db.Where("pending = TRUE")
	}

	var rows []entity.ChapterInvitationInfoRow

	if err := db.Order("created_at DESC").Find(&rows).Error; err != nil {
		return nil, err
	}

	items := make([]model.ChapterInvitationInfo, 0, len(rows))
	for _, row := range rows {
		items = append(items, entity.ToChapterInvitationInfo(row))
	}

	return items, nil
}

func (r *chapterInvitationRepoImpl) Create(c *model.ChapterInvitationCreation) (*model.ChapterInvitationInfo, error) {
	now := time.Now()
	row := map[string]any{
		"id":                 c.ID,
		"chapter_id":         c.ChapterID,
		"inviter_id":         c.InviterID,
		"invitee_qq":         c.InviteeQQ,
		"invitation_code":    c.InvitationCode,
		"pending":            true,
		"to_be_raw_provider": c.ToBeRawProvider,
		"to_be_translator":   c.ToBeTranslator,
		"to_be_proofreader":  c.ToBeProofreader,
		"to_be_typesetter":   c.ToBeTypesetter,
		"to_be_redrawer":     c.ToBeRedrawer,
		"to_be_reviewer":     c.ToBeReviewer,
		"to_be_publisher":    c.ToBePublisher,
		"created_at":         now,
		"updated_at":         now,
	}

	if err := r.gdb.Table(entity.ChapterInvitationTable).Create(row).Error; err != nil {
		return nil, err
	}

	var created entity.ChapterInvitationInfoRow

	if err := r.gdb.Table(entity.ChapterInvitationTable).Where("id = ?", c.ID).First(&created).Error; err != nil {
		return nil, err
	}

	info := entity.ToChapterInvitationInfo(created)
	return &info, nil
}

func (r *chapterInvitationRepoImpl) Invalidate(id string) error {
	return r.gdb.Table(entity.ChapterInvitationTable).
		Where("id = ?", id).
		Updates(map[string]any{
			"pending":    false,
			"updated_at": time.Now(),
		}).Error
}
