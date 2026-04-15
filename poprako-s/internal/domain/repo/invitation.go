package repo

import (
	"context"

	"poprako-s/internal/domain/model"
)

// MemberInvitationRepo 是邀请仓库的接口
type MemberInvitationRepo interface {
	// GetByInviteeQQ 根据被邀请求者 QQ 获取邀请信息；若不存在返回 error
	GetByInviteeQQ(id string) (*model.MemberInvitationInfo, error)
	// List 根据筛选条件返回邀请信息列表
	List(opt model.MemberInvitationQueryOpt) ([]model.MemberInvitationInfo, error)

	// Create 持久化一个新的邀请
	Create(c *model.MemberInvitationCreation) (*model.MemberInvitationInfo, error)
	// Update 更新邀请信息（修改待分配角色）
	Update(u *model.MemberInvitationUpdate) error
	// Invalidate 使邀请失效（标记为非 Pending）
	Invalidate(id string) error
	// Delete 删除邀请（硬删除）
	Delete(id string) error

	// FromTxnCx 从上下文中获取事务，并分离出一个带事务的 InvitationRepo 实例
	FromTxnCx(cx context.Context) (MemberInvitationRepo, error)
}
