package repo_infra

import (
	"time"

	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	"poprako-s/internal/domain/model/query"
	repo_iface "poprako-s/internal/domain/repo"
	"poprako-s/internal/infra/repo/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// `worksetRepoImpl` is the GORM-backed implementation of `repo_iface.WorksetRepo`.
type worksetRepoImpl struct {
	// `gdb` is the underlying GORM database handle.
	gdb *gorm.DB
}

// `NewWorksetRepo` creates a non-transaction-scoped `WorksetRepo`.
func NewWorksetRepo(gdb *gorm.DB) repo_iface.WorksetRepo {
	return &worksetRepoImpl{gdb: gdb}
}

// `GetById` retrieves a single workset by its primary key.
func (r *worksetRepoImpl) GetById(id string, inc ...enum.WorksetIncl) (*aggr.Workset, repo_iface.RepoErr) {
	var row entity.WorksetRow

	q := r.gdb.
		Table(entity.WORKSET_TABLE).
		Where("id = ?", id)

	q = withWorksetIncl(q, inc...)

	err := q.First(&row).Error
	if err != nil {
		return nil, err
	}

	return row.ToWorksetAggr(), nil
}

// `List` returns all worksets matching the given options, ordered by `index` ascending.
func (r *worksetRepoImpl) List(opt *query.ListWorksetOpt, inc ...enum.WorksetIncl) ([]*aggr.Workset, repo_iface.RepoErr) {
	var rows []entity.WorksetRow

	q := r.gdb.
		Table(entity.WORKSET_TABLE)

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

// `Count` returns the number of worksets matching the given options.
func (r *worksetRepoImpl) Count(opt *query.ListWorksetOpt) (int64, repo_iface.RepoErr) {
	var count int64

	q := r.gdb.
		Table(entity.WORKSET_TABLE)

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
		Where("id = ?", upd.Id).
		Select("name", "description", "updated_at").
		Updates(updRow).Error

	return err
}

// `UpdateComicCount` applies delta to `comic_count` and refreshes `updated_at`.
func (r *worksetRepoImpl) UpdateComicCount(id string, delta int) repo_iface.RepoErr {
	now := time.Now()

	err := r.gdb.
		Table(entity.WORKSET_TABLE).
		Where("id = ?", id).
		Updates(map[string]any{
			"comic_count": gorm.Expr("GREATEST(comic_count + ?, 0)", delta),
			"updated_at":  now,
		}).Error

	return err
}

// `IncrementComicNextIndex` allocates one next comic index from one workset row.
func (r *worksetRepoImpl) IncrementComicNextIndex(id string) (int, repo_iface.RepoErr) {
	type nextIndexRow struct {
		NextIndex int `gorm:"column:comic_next_index"`
	}

	var row nextIndexRow

	updRe := r.gdb.
		Table(entity.WORKSET_TABLE).
		Where("id = ?", id).
		Select("comic_next_index").
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "comic_next_index"}}}).
		Updates(map[string]any{"comic_next_index": gorm.Expr("comic_next_index + 1")}).
		Scan(&row)
	if updRe.Error != nil {
		return 0, updRe.Error
	}

	if updRe.RowsAffected == 0 {
		return 0, gorm.ErrRecordNotFound
	}

	return row.NextIndex - 1, nil
}

// `Delete` hard-deletes one workset row by id.
func (r *worksetRepoImpl) Delete(id string) repo_iface.RepoErr {
	return r.gdb.
		Table(entity.WORKSET_TABLE).
		Where("id = ?", id).
		Delete(&entity.WorksetRow{}).Error
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
