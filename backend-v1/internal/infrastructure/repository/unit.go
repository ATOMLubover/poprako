package repository

import (
	"labelplus-next-web-be/internal/domain/model"
	intf "labelplus-next-web-be/internal/domain/repository"
	"labelplus-next-web-be/internal/infrastructure/repository/entity"
	"labelplus-next-web-be/internal/util"

	"gorm.io/gorm/clause"
)

type unitRepository struct {
	executor intf.Executor
}

func NewUnitRepository(executor intf.Executor) intf.UnitRepository {
	return &unitRepository{executor: executor}
}

func (r *unitRepository) withTransaction(executor intf.Executor) intf.Executor {
	if executor != nil {
		return executor
	}

	return r.executor
}

func (r *unitRepository) BeginTransaction() intf.Executor {
	return r.executor.Begin()
}

func (r *unitRepository) LockByPageID(executor intf.Executor, pageID string) error {
	executor = r.withTransaction(executor)

	var lockedIDs []string

	return executor.
		Table(entity.UnitTable).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("page_id = ?", pageID).
		Pluck("id", &lockedIDs).Error
}

func (r *unitRepository) List(executor intf.Executor, options ...intf.QueryOption) ([]model.UnitInfo, error) {
	executor = r.withTransaction(executor)
	executor = executor.Table(entity.UnitTable)

	for _, opt := range options {
		executor = opt(executor)
	}

	var rows []entity.UnitInfoRow
	if err := executor.Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]model.UnitInfo, len(rows))
	for i, row := range rows {
		result[i] = entity.ToUnitInfo(row)
	}

	return result, nil
}

func (r *unitRepository) CreateBatch(executor intf.Executor, units []model.UnitCreation) error {
	executor = r.withTransaction(executor)

	if len(units) == 0 {
		return nil
	}

	rows := make([]entity.UnitInsertRow, len(units))
	for i, unit := range units {
		rows[i] = entity.UnitInsertRow{
			ID:                 util.GenerateUUID(),
			PageID:             unit.PageID,
			Index:              unit.Index,
			XCoord:             float64(unit.XCoord),
			YCoord:             float64(unit.YCoord),
			InBubble:           unit.IsBubble,
			IsProofread:        unit.IsProofread,
			TranslatedText:     unit.TranslatedText,
			TranslatorID:       unit.TranslatorID,
			TranslatorComment:  unit.TranslatorComment,
			ProofreaderText:    unit.ProofreadText,
			ProofreaderID:      unit.ProofreaderID,
			ProofreaderComment: unit.ProofreaderComment,
		}
	}

	return executor.Create(&rows).Error
}

// PatchBatch 对每个 patch 独立执行 UPDATE，仅更新 Some 字段。
// 若某行未找到不报错，满足协作场景的幂等要求。
func (r *unitRepository) PatchBatch(executor intf.Executor, patches []model.UnitPatch) error {
	executor = r.withTransaction(executor)

	for _, patch := range patches {
		updates := buildUnitPatchMap(patch)
		if len(updates) == 0 {
			continue
		}

		if err := executor.
			Table(entity.UnitTable).
			Where("id = ?", patch.ID).
			Updates(updates).Error; err != nil {
			return err
		}
	}

	return nil
}

func buildUnitPatchMap(patch model.UnitPatch) map[string]any {
	updates := make(map[string]any)

	if patch.Index.State() == util.OptionSome {
		updates["index"] = patch.Index.Unwrap()
	}

	if patch.XCoord.State() == util.OptionSome {
		updates["x_coord"] = float64(patch.XCoord.Unwrap())
	}

	if patch.YCoord.State() == util.OptionSome {
		updates["y_coord"] = float64(patch.YCoord.Unwrap())
	}

	if patch.IsBubble.State() == util.OptionSome {
		updates["in_bubble"] = patch.IsBubble.Unwrap()
	}

	if patch.TranslatedText.State() == util.OptionSome {
		updates["translated_text"] = patch.TranslatedText.Unwrap()
	}

	if patch.TranslatorID.State() == util.OptionSome {
		updates["translator_id"] = patch.TranslatorID.Unwrap()
	}

	if patch.TranslatorComment.State() == util.OptionSome {
		updates["translator_comment"] = patch.TranslatorComment.Unwrap()
	}

	if patch.IsProofread.State() == util.OptionSome {
		updates["is_proofread"] = patch.IsProofread.Unwrap()
	}

	if patch.ProofreadText.State() == util.OptionSome {
		updates["proofreader_text"] = patch.ProofreadText.Unwrap()
	}

	if patch.ProofreaderID.State() == util.OptionSome {
		updates["proofreader_id"] = patch.ProofreaderID.Unwrap()
	}

	if patch.ProofreaderComment.State() == util.OptionSome {
		updates["proofreader_comment"] = patch.ProofreaderComment.Unwrap()
	}

	return updates
}

// DeleteBatch 批量删除指定 ID 的 units，若某 ID 不存在不报错。
func (r *unitRepository) DeleteBatch(executor intf.Executor, unitIDs []string) error {
	executor = r.withTransaction(executor)

	if len(unitIDs) == 0 {
		return nil
	}

	return executor.
		Table(entity.UnitTable).
		Where("id IN ?", unitIDs).
		Delete(nil).Error
}
