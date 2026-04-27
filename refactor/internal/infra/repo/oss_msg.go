package repo_infra

import (
	"context"
	"errors"
	"time"

	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	repo_iface "poprako-s/internal/domain/repo"
	"poprako-s/internal/infra/repo/entity"

	"gorm.io/gorm"
)

// `ossMsgRepoImpl` is the gorm implementation of `OssMsgRepo`
type ossMsgRepoImpl struct {
	gdb *gorm.DB
}

// `TxnOssMsgRepo` creates a transaction-scoped `OssMsgRepo` from `cx`.
func TxnOssMsgRepo(cx context.Context) (repo_iface.OssMsgRepo, error) {
	// Ensure `cx` carries a transaction-scoped `gdb`.
	gdb := takeTxnGdb(cx)
	if gdb == nil {
		return nil, errors.New("[TxnOssMsgRepo] no transaction context found for OssMsgRepo")
	}

	return &ossMsgRepoImpl{gdb: gdb}, nil
}

// `NewOssMsgRepo` creates a non transaction-scoped `OssMsgRepo`
func NewOssMsgRepo(gdb *gorm.DB) repo_iface.OssMsgRepo {
	return &ossMsgRepoImpl{gdb: gdb}
}

// `SavePendingCre` saves a pending create message
func (r *ossMsgRepoImpl) SavePendingCre(msg *aggr.OssCreMsg) repo_iface.RepoErr {
	creRow := entity.NewOssCreMsgCreRowFromAggr(msg)

	upd := &entity.OssMsgRow{
		Id:         creRow.Id,
		Op:         creRow.Op,
		Status:     string(enum.OssMsgStatePend),
		ObjKeys:    creRow.ObjKeys,
		VisibleAt:  creRow.VisibleAt,
		ExpireAt:   creRow.ExpireAt,
		ProcAt:     nil,
		AttemptCnt: 0,
		LastErr:    "",
		UpdatedAt:  time.Now(),
	}

	res := r.gdb.
		Table(entity.OSS_MSG_TABLE).
		Where(
			"resource_type = ? AND resource_id = ? AND operation = ? AND status <> ?",
			creRow.ResTyp,
			creRow.ResId,
			creRow.Op,
			string(enum.OssMsgStateCmpl),
		).
		Select("id", "operation", "status", "object_keys", "visible_at", "expire_at", "processing_at", "attempt_count", "last_error", "updated_at").
		Updates(upd)
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected > 0 {
		return nil
	}

	err := r.gdb.
		Table(creRow.TableName()).
		Create(creRow).Error
	if err != nil {
		return err
	}

	return nil
}

// `SavePendingDel` saves a pending delete message
func (r *ossMsgRepoImpl) SavePendingDel(msg *aggr.OssDelMsg) repo_iface.RepoErr {
	creRow := entity.NewOssDelMsgCreRowFromAggr(msg)

	err := r.gdb.
		Table(creRow.TableName()).
		Create(creRow).Error
	if err != nil {
		return err
	}

	return nil
}

// `MarkCmpl` marks one message completed by id
func (r *ossMsgRepoImpl) MarkCmpl(id string) repo_iface.RepoErr {
	upd := entity.NewOssMsgMarkCmplUpdRow()

	err := r.gdb.
		Table(entity.OSS_MSG_TABLE).
		Where("id = ?", id).
		Select("status", "updated_at").
		Updates(upd).Error
	if err != nil {
		return err
	}

	return nil
}

// `MarkCmplByRes` marks one pending create message completed by resource identity
func (r *ossMsgRepoImpl) MarkCmplByRes(ty enum.OssResTyp, resId string) repo_iface.RepoErr {
	var row entity.OssMsgRow

	err := r.gdb.
		Table(entity.OSS_MSG_TABLE).
		Where(
			"resource_type = ? AND resource_id = ? AND operation = ? AND status = ?",
			string(ty),
			resId,
			string(enum.OssOpCre),
			string(enum.OssMsgStatePend),
		).
		Order("created_at DESC").
		First(&row).Error
	if err != nil {
		return err
	}

	upd := entity.NewOssMsgMarkCmplUpdRow()

	err = r.gdb.
		Table(entity.OSS_MSG_TABLE).
		Where("id = ?", row.Id).
		Select("status", "updated_at").
		Updates(upd).Error
	if err != nil {
		return err
	}

	return nil
}

// `ClaimPending` claims one pending message by operation and marks it processing.
func (r *ossMsgRepoImpl) ClaimPending(op enum.OssOp) (*aggr.OssMsg, repo_iface.RepoErr) {
	var row entity.OssMsgRow

	now := time.Now()

	res := r.gdb.
		Table(entity.OSS_MSG_TABLE).
		Where("status = ? AND operation = ? AND visible_at <= ?", string(enum.OssMsgStatePend), string(op), now).
		Order("visible_at ASC").
		Limit(1).
		Find(&row)
	if res.Error != nil {
		return nil, res.Error
	}

	if res.RowsAffected == 0 {
		return nil, nil
	}

	upd := entity.NewOssMsgMarkProcUpdRow(now)

	res = r.gdb.
		Table(entity.OSS_MSG_TABLE).
		Where("id = ? AND status = ?", row.Id, string(enum.OssMsgStatePend)).
		Select("status", "processing_at", "updated_at").
		Updates(upd)
	if res.Error != nil {
		return nil, res.Error
	}

	if res.RowsAffected == 0 {
		// Message is claimed by another worker.
		return nil, nil
	}

	row.Status = string(enum.OssMsgStateProc)
	row.ProcAt = &now

	return row.ToOssMsgAggr(), nil
}

// `MarkPending` resets one processing message to pending and bumps attempt count.
func (r *ossMsgRepoImpl) MarkPending(id string, errMsg string, nextVisibleAt time.Time) repo_iface.RepoErr {
	now := time.Now()
	upd := entity.NewOssMsgMarkPendingUpdRow(errMsg, nextVisibleAt, now)

	// Reset all pending fields in a single update to reduce partial-failure window.
	err := r.gdb.
		Table(entity.OSS_MSG_TABLE).
		Where("id = ?", id).
		Select("status", "processing_at", "visible_at", "last_error", "updated_at").
		Updates(upd).Error
	if err != nil {
		return err
	}

	// Increment attempt_count atomically with a column expression.
	err = r.gdb.
		Table(entity.OSS_MSG_TABLE).
		Where("id = ?", id).
		UpdateColumn("attempt_count", gorm.Expr("attempt_count + ?", 1)).Error
	if err != nil {
		return err
	}

	return nil
}

// `ResetStuck` resets stuck processing messages to pending state
func (r *ossMsgRepoImpl) ResetStuck(bef time.Time) repo_iface.RepoErr {
	upd := entity.NewOssMsgResetStuckUpdRow()

	err := r.gdb.
		Table(entity.OSS_MSG_TABLE).
		Where(
			"status = ? AND processing_at IS NOT NULL AND processing_at <= ?",
			string(enum.OssMsgStateProc),
			bef,
		).
		Select("status", "processing_at", "updated_at").
		Updates(upd).Error
	if err != nil {
		return err
	}

	return nil
}

// `CleanCmpl` deletes completed messages created before cutoff time
func (r *ossMsgRepoImpl) CleanCmpl(bef time.Time) repo_iface.RepoErr {
	err := r.gdb.
		Table(entity.OSS_MSG_TABLE).
		Where("status = ? AND created_at <= ?", string(enum.OssMsgStateCmpl), bef).
		Delete(&entity.OssMsgRow{}).Error
	if err != nil {
		return err
	}

	return nil
}
