package repo_infra

import (
	"context"
	"errors"
	"time"

	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	"poprako-s/internal/domain/model/query"
	repo_iface "poprako-s/internal/domain/repo"
	"poprako-s/internal/infra/repo/entity"

	"gorm.io/gorm"
)

// `worksetRepoImpl` is the GORM-backed implementation of `repo_iface.WorksetRepo`.
type worksetRepoImpl struct {
	// `gdb` is the underlying GORM database handle.
	gdb *gorm.DB
}

// `TxnWorksetRepo` creates a transaction-scoped `WorksetRepo` extracted from `cx`.
func TxnWorksetRepo(cx context.Context) (repo_iface.WorksetRepo, repo_iface.RepoErr) {
	gdb := takeTxnGdb(cx)
	if gdb == nil {
		return nil, errors.New("[TxnWorksetRepo] no transaction context found for WorksetRepo")
	}

	return &worksetRepoImpl{gdb: gdb}, nil
}

// `NewWorksetRepo` creates a non-transaction-scoped `WorksetRepo`.
func NewWorksetRepo(gdb *gorm.DB) repo_iface.WorksetRepo {
	return &worksetRepoImpl{gdb: gdb}
}

// `GetById` retrieves a single active workset by its primary key.
func (r *worksetRepoImpl) GetById(id string, inc ...enum.WorksetIncl) (*aggr.Workset, repo_iface.RepoErr) {
	var row entity.WorksetRow

	q := r.gdb.
		Table(entity.WORKSET_TABLE).
		Where("id = ? AND deleted_at IS NULL", id)

	q = withWorksetIncl(q, inc...)

	err := q.First(&row).Error
	if err != nil {
		return nil, err
	}

	return row.ToWorksetAggr(), nil
}

// `List` returns all active worksets matching the given options, ordered by `index` ascending.
func (r *worksetRepoImpl) List(opt *query.ListWorksetOpt, inc ...enum.WorksetIncl) ([]*aggr.Workset, repo_iface.RepoErr) {
	var rows []entity.WorksetRow

	q := r.gdb.
		Table(entity.WORKSET_TABLE).
		Where("deleted_at IS NULL")

	// Apply optional filters.
	if opt != nil && opt.TeamId != nil {
		q = q.Where("team_id = ?", *opt.TeamId)
	}

	if opt != nil && opt.Pagi.Offset > 0 {
		q = q.Offset(opt.Pagi.Offset)
	}

	if opt != nil && opt.Pagi.Limit > 0 {
		q = q.Limit(opt.Pagi.Limit)
	}

	q = withWorksetIncl(q, inc...)

	err := q.Order("index ASC").Find(&rows).Error
	if err != nil {
		return nil, err
	}

	// Convert each row to an aggregate.
	result := make([]*aggr.Workset, len(rows))

	for i := range rows {
		result[i] = rows[i].ToWorksetAggr()
	}

	return result, nil
}

// `Count` returns the number of active worksets matching the given options.
func (r *worksetRepoImpl) Count(opt *query.ListWorksetOpt) (int64, repo_iface.RepoErr) {
	var count int64

	q := r.gdb.
		Table(entity.WORKSET_TABLE).
		Where("deleted_at IS NULL")

	// Apply optional filters.
	if opt != nil && opt.TeamId != nil {
		q = q.Where("team_id = ?", *opt.TeamId)
	}

	err := q.Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

// `Create` inserts a new workset row and returns the fully populated aggregate.
func (r *worksetRepoImpl) Create(cre *aggr.WorksetCre) (*aggr.Workset, repo_iface.RepoErr) {
	row := entity.NewWorksetCreRowFromAggr(cre)

	err := r.gdb.
		Table(row.TableName()).
		Create(row).Error
	if err != nil {
		return nil, err
	}

	// Re-fetch to return the full aggregate including server-side defaults.
	return r.GetById(row.Id)
}

// `Update` applies the mutable fields of `upd` to the matching workset row.
func (r *worksetRepoImpl) Update(upd *aggr.WorksetUpd) repo_iface.RepoErr {
	now := time.Now()

	updRow := &entity.WorksetUpdRow{
		Name:      upd.Name,
		Desc:      upd.Desc,
		UpdatedAt: now,
	}

	err := r.gdb.
		Table(entity.WORKSET_TABLE).
		Where("id = ? AND deleted_at IS NULL", upd.Id).
		Select("name", "description", "updated_at").
		Updates(updRow).Error

	return err
}

// `UpdateComicCount` applies delta to `comic_count` and refreshes `updated_at`.
func (r *worksetRepoImpl) UpdateComicCount(id string, delta int) repo_iface.RepoErr {
	now := time.Now()

	err := r.gdb.
		Table(entity.WORKSET_TABLE).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{
			"comic_count": gorm.Expr("GREATEST(comic_count + ?, 0)", delta),
			"updated_at":  now,
		}).Error

	return err
}

// `Remove` marks the workset as deleted by setting `deleted_at` to the current time.
func (r *worksetRepoImpl) Remove(id string) repo_iface.RepoErr {
	now := time.Now()

	err := r.gdb.
		Table(entity.WORKSET_TABLE).
		Where("id = ? AND deleted_at IS NULL", id).
		UpdateColumn("deleted_at", now).Error

	return err
}

// `withWorksetIncl` applies typed include options to the base query.
func withWorksetIncl(q *gorm.DB, inc ...enum.WorksetIncl) *gorm.DB {
	for _, i := range inc {
		switch i {
		case enum.WorksetInclTeam:
			q = q.Preload("Team")
		}
	}

	return q
}
