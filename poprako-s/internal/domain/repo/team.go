package repo

import (
	"context"

	"poprako-s/internal/domain/model"
)

// TeamRepo 是汉化组仓库的接口
type TeamRepo interface {
	// GetByID 根据汉化组 ID 获取汉化组信息；若不存在返回 error
	GetByID(id string) (*model.TeamInfo, error)
	// List 根据筛选条件返回汉化组信息列表
	List(opt model.TeamQueryOpt) ([]model.TeamInfo, error)

	// Create 持久化一个新的汉化组
	Create(c *model.TeamCreation) (*model.TeamInfo, error)
	// Update 更新汉化组信息（PUT 语义）
	Update(u *model.TeamUpdate) error
	// Delete 删除汉化组（硬删除）
	Delete(id string) error

	// PreFillAvatarOSSKey 预填充汉化组头像的 OSS Key
	PreFillAvatarOSSKey(id string, avatarOSSKey string) error
	// ConfirmAvatarUploaded 确认汉化组头像已上传
	ConfirmAvatarUploaded(id string) error

	// FromTxnCx 从上下文中获取事务，并分离出一个带事务的 TeamRepo 实例
	FromTxnCx(cx context.Context) (TeamRepo, error)
}
