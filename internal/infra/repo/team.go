package repo_infra

import (
	"fmt"
	"time"

	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/query"
	repo_iface "poprako-s/internal/domain/repo"
	"poprako-s/internal/infra/repo/entity"

	"gorm.io/gorm"
)

// `teamRepoImpl` implements `repo_iface.TeamRepo` with `gorm`.
type teamRepoImpl struct {
	gdb *gorm.DB
}

// `NewTeamRepo` creates a non-transaction-scoped `TeamRepo`.
func NewTeamRepo(gdb *gorm.DB) repo_iface.TeamRepo {
	return &teamRepoImpl{gdb: gdb}
}

// `GetById` retrieves one `Team` aggregate by its id.
func (r *teamRepoImpl) GetById(id string) (*aggr.Team, repo_iface.RepoErr) {
	var row entity.TeamRow

	// Query the team row by id from `t_team`.
	err := r.gdb.
		Table(entity.TEAM_TABLE).
		Where("id = ?", id).
		First(&row).Error
	if err != nil {
		return nil, err
	}

	// Convert the row into `aggr.Team`.
	return row.ToTeamAggr(), nil
}

// `List` returns teams by list options and pagination.
func (r *teamRepoImpl) List(opt *query.ListTeamOpt) ([]*aggr.Team, repo_iface.RepoErr) {
	var rows []entity.TeamRow

	qry := r.gdb.
		Table(entity.TEAM_TABLE)

	if opt != nil && opt.Id != nil {
		qry = qry.Where("id = ?", *opt.Id)
	}

	if opt != nil && opt.Pagi.Offset > 0 {
		qry = qry.Offset(opt.Pagi.Offset)
	}

	if opt != nil && opt.Pagi.Limit > 0 {
		qry = qry.Limit(opt.Pagi.Limit)
	}

	err := qry.
		Order("created_at DESC").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	teams := make([]*aggr.Team, len(rows))
	for i := range rows {
		teams[i] = rows[i].ToTeamAggr()
	}

	return teams, nil
}

// `Create` inserts one team row and returns created aggregate.
func (r *teamRepoImpl) Create(cre *aggr.TeamCre) (*aggr.Team, repo_iface.RepoErr) {
	row := entity.NewTeamCreRowFromAggr(cre)

	err := r.gdb.
		Table(entity.TEAM_TABLE).
		Create(row).Error
	if err != nil {
		return nil, err
	}

	return r.GetById(row.Id)
}

// `Update` applies put-style update to one team row.
func (r *teamRepoImpl) Update(upd *aggr.TeamUpd) repo_iface.RepoErr {
	now := time.Now()

	updRow := &entity.TeamUpdRow{
		Name:      upd.Name,
		Desc:      upd.Desc,
		UpdatedAt: now,
	}

	updRe := r.gdb.
		Table(entity.TEAM_TABLE).
		Where("id = ?", upd.Id).
		Select("name", "description", "updated_at").
		Updates(updRow)
	if updRe.Error != nil {
		return updRe.Error
	}

	if updRe.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// `Delete` executes hard delete on one team row by id.
func (r *teamRepoImpl) Delete(id string) repo_iface.RepoErr {
	queryRe := r.gdb.
		Table(entity.TEAM_TABLE).
		Where("id = ?", id).
		Delete(&entity.TeamRow{})
	if queryRe.Error != nil {
		return queryRe.Error
	}

	if queryRe.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// `PrefillAvatarKey` writes one reserved avatar key before upload starts.
func (r *teamRepoImpl) PrefillAvatarKey(id string, key string) repo_iface.RepoErr {
	now := time.Now()

	updRe := r.gdb.
		Table(entity.TEAM_TABLE).
		Where("id = ?", id).
		Updates(map[string]any{
			"avatar_key":      key,
			"avatar_uploaded": false,
			"updated_at":      now,
		})
	if updRe.Error != nil {
		return updRe.Error
	}

	if updRe.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// `MarkAvatarUploaded` marks avatar upload status as completed.
func (r *teamRepoImpl) MarkAvatarUploaded(id string) repo_iface.RepoErr {
	now := time.Now()

	updRe := r.gdb.
		Table(entity.TEAM_TABLE).
		Where("id = ?", id).
		Updates(map[string]any{
			"avatar_uploaded": true,
			"updated_at":      now,
		})
	if updRe.Error != nil {
		return updRe.Error
	}

	if updRe.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// `IncrementWorksetNextIndex` allocates one next workset index from one team row.
func (r *teamRepoImpl) IncrementWorksetNextIndex(id string) (int, repo_iface.RepoErr) {
	type nextIndexRow struct {
		NextIndex int `gorm:"column:workset_next_index"`
	}

	var row nextIndexRow

	// Atomically increment `workset_next_index` and fetch the new value.
	updRe := r.gdb.Raw(
		fmt.Sprintf(
			"UPDATE %s SET workset_next_index = workset_next_index + 1 WHERE id = ? RETURNING workset_next_index",
			entity.TEAM_TABLE,
		),
		id,
	).Scan(&row)
	if updRe.Error != nil {
		return 0, updRe.Error
	}

	if updRe.RowsAffected == 0 {
		return 0, gorm.ErrRecordNotFound
	}

	return row.NextIndex - 1, nil
}
