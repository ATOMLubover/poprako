package repo_infra

import (
	"poprako-s/internal/domain/model/aggr"
	repo_iface "poprako-s/internal/domain/repo"
	"poprako-s/internal/infra/repo/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// `userStatsRepoImpl` is GORM-backed implementation of `UserStatsRepo`.
type userStatsRepoImpl struct {
	gdb *gorm.DB
}

// `NewUserStatsRepo` creates a non-transaction-scoped `UserStatsRepo`.
func NewUserStatsRepo(gdb *gorm.DB) repo_iface.UserStatsRepo {
	return &userStatsRepoImpl{gdb: gdb}
}

// `EnsureGet` returns existing user stats row or creates one lazily.
func (r *userStatsRepoImpl) EnsureGet(id string) (*aggr.UserStats, repo_iface.RepoErr) {
	var row entity.UserStatsRow

	err := r.gdb.
		Table(entity.USER_STATS_TABLE).
		Where("user_id = ?", id).
		First(&row).Error
	if err == nil {
		return row.ToUserStatsAggr(), nil
	}

	if !IsNotFound(err) {
		return nil, err
	}

	cre := entity.NewUserStatsCreRow(id)

	err = r.gdb.
		Table(entity.USER_STATS_TABLE).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}},
			DoNothing: true,
		}).
		Create(cre).Error
	if err != nil {
		return nil, err
	}

	err = r.gdb.
		Table(entity.USER_STATS_TABLE).
		Where("user_id = ?", id).
		First(&row).Error
	if err != nil {
		return nil, err
	}

	return row.ToUserStatsAggr(), nil
}

// `Patch` applies delta update to one user stats row.
func (r *userStatsRepoImpl) Patch(pat *aggr.UserStatsPatch) repo_iface.RepoErr {
	if pat == nil || pat.UserId == "" {
		return gorm.ErrInvalidData
	}

	if _, err := r.EnsureGet(pat.UserId); err != nil {
		return err
	}

	res := r.gdb.
		Table(entity.USER_STATS_TABLE).
		Where("user_id = ?", pat.UserId).
		Updates(map[string]any{
			"total_assignment_count":    gorm.Expr("GREATEST(total_assignment_count + ?, 0)", pat.TotalAssignmentDelta),
			"active_assignment_count":   gorm.Expr("GREATEST(active_assignment_count + ?, 0)", pat.ActiveAssignmentDelta),
			"finished_assignment_count": gorm.Expr("GREATEST(finished_assignment_count + ?, 0)", pat.FinishedAssignmentDelta),
			"updated_at":                gorm.Expr("NOW()"),
		})
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
