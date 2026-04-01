package repo

import "poprako-s/internal/domain/model"

// AssignmentRepo 是分配记录仓库的接口
type AssignmentRepo interface {
	// GetByID 根据分配记录 ID 获取分配信息；若不存在返回 error
	GetByID(id string) (*model.AssignmentInfo, error)
	// Get 根据筛选条件获取单条分配信息；若不存在返回 error
	Get(opt model.AssignmentQueryOpt) (*model.AssignmentInfo, error)
	// List 根据筛选条件返回分配信息列表
	List(opt model.AssignmentQueryOpt) ([]model.AssignmentInfo, error)
	// Exist 根据筛选条件判断是否存在匹配的分配记录
	Exist(opt model.AssignmentQueryOpt) (bool, error)

	// Create 持久化一个新的分配记录
	Create(c *model.AssignmentCreation) (*model.AssignmentInfo, error)
	// Update 更新分配记录（PUT 语义，全量角色替换）
	Update(u *model.AssignmentUpdate) error
	// Delete 删除分配记录（硬删除）
	Delete(id string) error
}
