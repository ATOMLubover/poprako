package repo_infra

import (
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/query"
	repo_iface "poprako-s/internal/domain/repo"
	"poprako-s/internal/infra/repo/entity"

	"gorm.io/gorm"
)

// `sysMailRepoImpl` is the gorm implementation of `SysMailRepo`.
type sysMailRepoImpl struct {
	gdb *gorm.DB
}

// `NewSysMailRepo` creates a non transaction-scoped `SysMailRepo`.
func NewSysMailRepo(gdb *gorm.DB) repo_iface.SysMailRepo {
	return &sysMailRepoImpl{gdb: gdb}
}

// `Send` creates a new system mail record for the given creation aggregate.
func (r *sysMailRepoImpl) Send(cre *aggr.SysMailCre) repo_iface.RepoErr {
	creRow := entity.NewSysMailCreRowFromAggr(cre)

	err := r.gdb.
		Table(creRow.TableName()).
		Create(creRow).Error
	if err != nil {
		return err
	}

	return nil
}

// `ListUnreadByRcvId` lists unread system mails by receiver with pagination.
func (r *sysMailRepoImpl) ListUnreadByRcvId(rcvId string, pagi query.PagiOpt) ([]aggr.SysMail, repo_iface.RepoErr) {
	rows := []entity.SysMailRow{}

	qry := r.gdb.
		Table(entity.SYS_MAIL_TABLE).
		Where("receiver_id = ? AND read = FALSE", rcvId)

	if pagi.Offset > 0 {
		qry = qry.Offset(pagi.Offset)
	}

	if pagi.Limit > 0 {
		qry = qry.Limit(pagi.Limit)
	}

	err := qry.
		Order("created_at DESC").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	items := make([]aggr.SysMail, len(rows))

	for i := range rows {
		item := rows[i].ToSysMailAggr()
		if item != nil {
			items[i] = *item
		}
	}

	return items, nil
}

// `MarkReadByRcvId` marks one system mail record as read by id and receiver id.
func (r *sysMailRepoImpl) MarkReadByRcvId(id string, rcvId string) repo_iface.RepoErr {
	upd := entity.NewSysMailMarkReadUpdRow()

	res := r.gdb.
		Table(entity.SYS_MAIL_TABLE).
		Where("id = ? AND receiver_id = ?", id, rcvId).
		Select("read").
		Updates(upd)
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
