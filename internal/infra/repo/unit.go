package repo_infra

import (
	"strings"

	"poprako-s/internal/domain/model/aggr"
	repo_iface "poprako-s/internal/domain/repo"
	"poprako-s/internal/infra/repo/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// `unitRepoImpl` is the GORM-backed implementation of `UnitRepo`.
type unitRepoImpl struct {
	// `gdb` is the underlying GORM handle.
	gdb *gorm.DB
}

// `unitCountRow` maps one page count query result.
type unitCountRow struct {
	// `Total` is total unit count.
	Total int `gorm:"column:total"`

	// `Translated` is translated unit count.
	Translated int `gorm:"column:translated"`

	// `Proofread` is proofread unit count.
	Proofread int `gorm:"column:proofread"`
}

// `NewUnitRepo` creates a non-transaction-scoped `UnitRepo`.
func NewUnitRepo(gdb *gorm.DB) repo_iface.UnitRepo {
	return &unitRepoImpl{gdb: gdb}
}

// `ListByPage` returns units under one page ordered by `Index` ascending.
func (r *unitRepoImpl) ListByPage(pageId string) ([]*aggr.Unit, repo_iface.RepoErr) {
	var rows []entity.UnitRow

	err := r.gdb.
		Table(entity.UNIT_TABLE).
		Where("page_id = ?", pageId).
		Order("index ASC").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]*aggr.Unit, len(rows))
	for i := range rows {
		result[i] = rows[i].ToUnitAggr()
	}

	return result, nil
}

// `ListIndicesByPage` returns page unit ids and indices.
func (r *unitRepoImpl) ListIndicesByPage(pageId string) ([]*aggr.UnitIndex, repo_iface.RepoErr) {
	var rows []entity.UnitRow

	err := r.gdb.
		Table(entity.UNIT_TABLE).
		Select("id", "index").
		Where("page_id = ?", pageId).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]*aggr.UnitIndex, len(rows))
	for i := range rows {
		result[i] = &aggr.UnitIndex{Id: rows[i].Id, Index: rows[i].Index}
	}

	return result, nil
}

// `CountByPage` returns total, translated, and proofread counts of one page.
func (r *unitRepoImpl) CountByPage(pageId string) (int, int, int, repo_iface.RepoErr) {
	var row unitCountRow

	err := r.gdb.Raw(`
		SELECT
			COUNT(*) AS total,
			COUNT(CASE WHEN NULLIF(translated_text, '') IS NOT NULL THEN 1 END) AS translated,
			COUNT(CASE WHEN is_proofread THEN 1 END) AS proofread
		FROM t_unit
		WHERE page_id = ?
	`, pageId).Scan(&row).Error
	if err != nil {
		return 0, 0, 0, err
	}

	return row.Total, row.Translated, row.Proofread, nil
}

// `Create` inserts one unit.
func (r *unitRepoImpl) Create(cre *aggr.UnitCre) repo_iface.RepoErr {
	row := entity.NewUnitCreRowFromAggr(cre)

	return r.gdb.Table(entity.UNIT_TABLE).Create(row).Error
}

// `Save` executes one unit upsert.
func (r *unitRepoImpl) Save(sv *aggr.UnitSave) repo_iface.RepoErr {
	row := entity.NewUnitUpdRowFromAggr(sv)

	return r.gdb.
		Table(entity.UNIT_TABLE).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).
		Create(row).Error
}

// `Reindex` applies batch index update to page units.
func (r *unitRepoImpl) Reindex(indices []*aggr.UnitIndex) repo_iface.RepoErr {
	if len(indices) == 0 {
		return nil
	}

	sql, args := mkUnitReindexStmt(indices)

	return r.gdb.Exec(sql, args...).Error
}

// `Delete` hard-deletes one unit.
func (r *unitRepoImpl) Delete(id string) repo_iface.RepoErr {
	return r.gdb.Table(entity.UNIT_TABLE).Where("id = ?", id).Delete(nil).Error
}

// `mkUnitReindexStmt` builds the batch reindex SQL and bind arguments.
func mkUnitReindexStmt(indices []*aggr.UnitIndex) (string, []any) {
	args := make([]any, 0, len(indices)*3+1)
	b := &strings.Builder{}

	b.WriteString("UPDATE ")
	b.WriteString(entity.UNIT_TABLE)
	b.WriteString(" SET \"index\" = CASE id")

	for i := range indices {
		b.WriteString(" WHEN ? THEN ?")
		args = append(args, indices[i].Id, indices[i].Index)
	}

	b.WriteString(" ELSE \"index\" END, updated_at = NOW() WHERE id IN (")

	for i := range indices {
		if i > 0 {
			b.WriteString(", ")
		}

		b.WriteString("?")
		args = append(args, indices[i].Id)
	}

	b.WriteString(")")

	return b.String(), args
}
