package repo

import "poprako-s/internal/domain/model"

// UnitRepo 是翻译单元仓库的接口
type UnitRepo interface {
	// List 根据筛选条件返回翻译单元列表
	List(opt model.UnitQueryOpt) ([]model.UnitInfo, error)

	// CreateBatch 批量创建翻译单元
	CreateBatch(units []*model.UnitCreation) error
	// PatchBatch 批量部分更新翻译单元；若某行未找到不报错，符合协作场景下的幂等要求
	PatchBatch(patches []*model.UnitPatch) error
	// DeleteBatch 批量删除翻译单元；若某行不存在不报错
	DeleteBatch(unitIDs []string) error
}
