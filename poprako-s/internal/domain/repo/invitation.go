package repo

import "poprako-s/internal/domain/model"

// InvitationRepo 是邀请仓库的接口
type InvitationRepo interface {
	// GetByInviteeQQ 根据被邀请求者 QQ 获取邀请信息；若不存在返回 error
	GetByInviteeQQ(id string) (*model.InvitationInfo, error)
	// List 根据筛选条件返回邀请信息列表
	List(opt model.InvitationQueryOpt) ([]model.InvitationInfo, error)

	// Create 持久化一个新的邀请
	Create(c *model.InvitationCreation) (*model.InvitationInfo, error)
	// Update 更新邀请信息（修改待分配角色）
	Update(u *model.InvitationUpdate) error
	// Invalidate 使邀请失效（标记为非 Pending）
	Invalidate(id string) error
	// Delete 删除邀请（硬删除）
	Delete(id string) error
}
