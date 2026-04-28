package repo_infra

import (
	"poprako-s/internal/domain/model/aggr"
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

// `MarkRead` marks one system mail record as read.
func (r *sysMailRepoImpl) MarkRead(id string) repo_iface.RepoErr {
	upd := entity.NewSysMailMarkReadUpdRow()

	err := r.gdb.
		Table(entity.SYS_MAIL_TABLE).
		Where("id = ?", id).
		Select("read").
		Updates(upd).Error
	if err != nil {
		return err
	}

	return nil
}
