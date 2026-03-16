package repository

import "labelplus-next-web-be/internal/domain/model"

type UnitRepository interface {
	Transactor
	// LockByPageID 通过悲观锁锁定指定 page 下的所有行，用于并发 Save 场景。
	LockByPageID(executor Executor, pageID string) error
	List(executor Executor, options ...QueryOption) ([]model.UnitInfo, error)
	CreateBatch(executor Executor, units []model.UnitCreation) error
	// PatchBatch 批量 patch 更新，若某行未找到不报错，符合协作场景下的幂等要求。
	PatchBatch(executor Executor, patches []model.UnitPatch) error
	// DeleteBatch 批量删除，若某行不存在不报错。
	DeleteBatch(executor Executor, unitIDs []string) error
}
