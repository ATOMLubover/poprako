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

type unitRepoImpl struct {
	gdb *gorm.DB
}

func NewUnitRepo(gdb *gorm.DB) iface.UnitRepo {
	return &unitRepoImpl{gdb: gdb}
}

func NewUnitRepoFromCx(cx context.Context) (iface.UnitRepo, error) {
	gdb, ok := cx.Value(txnKey).(*gorm.DB)
	if !ok {
		return nil, errors.New("[NewUnitRepoFromCx]: 无法从上下文中获取事务数据库连接")
	}

	return &unitRepoImpl{gdb: gdb}, nil
}

func (r *unitRepoImpl) FromTxnCx(cx context.Context) (iface.UnitRepo, error) {
	return NewUnitRepoFromCx(cx)
}

func (r *unitRepoImpl) List(opt model.UnitQueryOpt) ([]model.UnitInfo, error) {
	db := r.gdb.Table(entity.UnitTable)
	if opt.PageID != "" {
		db = db.Where("page_id = ?", opt.PageID)
	}

	var rows []entity.UnitInfoRow
	if err := db.Order("index ASC").Find(&rows).Error; err != nil {
		return nil, err
	}

	items := make([]model.UnitInfo, 0, len(rows))
	for _, row := range rows {
		items = append(items, entity.ToUnitInfo(row))
	}

	return items, nil
}

func (r *unitRepoImpl) CreateBatch(units []*model.UnitCreation) error {
	if len(units) == 0 {
		return nil
	}

	now := time.Now()
	rows := make([]map[string]any, 0, len(units))
	for _, unit := range units {
		rows = append(rows, map[string]any{
			"id":                  unit.ID,
			"page_id":             unit.PageID,
			"x_coord":             float64(unit.XCoord),
			"y_coord":             float64(unit.YCoord),
			"index":               unit.Index,
			"in_bubble":           unit.IsBubble,
			"is_proofread":        unit.IsProofread,
			"translated_text":     unit.TranslatedText,
			"translator_id":       unit.TranslatorID,
			"translator_comment":  unit.TranslatorComment,
			"proofreader_text":    unit.ProofreadText,
			"proofreader_id":      unit.ProofreaderID,
			"proofreader_comment": unit.ProofreaderComment,
			"created_at":          now,
			"updated_at":          now,
		})
	}

	return r.gdb.Table(entity.UnitTable).Create(rows).Error
}

func (r *unitRepoImpl) PatchBatch(patches []*model.UnitPatch) error {
	for _, patch := range patches {
		updates := map[string]any{
			"updated_at": time.Now(),
		}

		if patch.Index != nil {
			updates["index"] = *patch.Index
		}
		if patch.XCoord != nil {
			updates["x_coord"] = float64(*patch.XCoord)
		}
		if patch.YCoord != nil {
			updates["y_coord"] = float64(*patch.YCoord)
		}
		if patch.IsBubble != nil {
			updates["in_bubble"] = *patch.IsBubble
		}
		if patch.TranslatedText != nil {
			updates["translated_text"] = *patch.TranslatedText
		}
		if patch.TranslatorID != nil {
			updates["translator_id"] = *patch.TranslatorID
		}
		if patch.TranslatorComment != nil {
			updates["translator_comment"] = *patch.TranslatorComment
		}
		if patch.IsProofread != nil {
			updates["is_proofread"] = *patch.IsProofread
		}
		if patch.ProofreadText != nil {
			updates["proofreader_text"] = *patch.ProofreadText
		}
		if patch.ProofreaderID != nil {
			updates["proofreader_id"] = *patch.ProofreaderID
		}
		if patch.ProofreaderComment != nil {
			updates["proofreader_comment"] = *patch.ProofreaderComment
		}

		if err := r.gdb.Table(entity.UnitTable).Where("id = ?", patch.ID).Updates(updates).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *unitRepoImpl) DeleteBatch(unitIDs []string) error {
	if len(unitIDs) == 0 {
		return nil
	}

	return r.gdb.Table(entity.UnitTable).Where("id IN ?", unitIDs).Delete(nil).Error
}
