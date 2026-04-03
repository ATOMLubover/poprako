package repo

import (
	"context"

	"poprako-s/internal/domain/model"
)

// MemberRepo 是成员仓库的接口
type MemberRepo interface {
	// GetByID 根据成员 ID 获取成员信息；若不存在返回 error
	GetByID(id string) (*model.MemberInfo, error)
	// Get 根据筛选条件获取单条成员信息；若不存在返回 error
	Get(opt model.MemberQueryOpt) (*model.MemberInfo, error)
	// List 根据筛选条件返回成员信息列表
	List(opt model.MemberQueryOpt) ([]model.MemberInfo, error)
	// Exist 根据筛选条件判断是否存在匹配的成员
	Exist(opt model.MemberQueryOpt) (bool, error)

	// Create 持久化一个新的成员
	Create(c *model.MemberCreation) (*model.MemberInfo, error)
	// Update 更新成员信息（PUT 语义，全量角色替换）
	Update(u *model.MemberUpdate) error
	// Delete 删除成员（硬删除）
	Delete(id string) error

	// FromTxnCx 从上下文中获取事务，并分离出一个带事务的 MemberRepo 实例
	FromTxnCx(cx context.Context) (MemberRepo, error)
}
