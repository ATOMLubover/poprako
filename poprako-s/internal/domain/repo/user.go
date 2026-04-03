package repo

import (
	"context"
	"time"

	"poprako-s/internal/domain/model"
)

// UserRepo 是用户仓库的接口
type UserRepo interface {
	// GetCredsByQQ 根据 QQ 号获取用户凭证信息；若不存在返回 error
	GetCredsByQQ(qq string) (*model.UserCreds, error)
	// GetByID 根据用户 ID 获取用户信息；若不存在返回 error
	GetByID(id string) (*model.UserInfo, error)
	// GetByQQ 根据 QQ 号获取用户信息；若不存在返回 error
	GetByQQ(qq string) (*model.UserInfo, error)
	// List 根据筛选条件返回用户信息列表
	List(opt model.UserQueryOpt) ([]model.UserInfo, error)

	// Create 持久化一个新的用户
	Create(c *model.UserCreation) (*model.UserInfo, error)
	// Update 更新用户信息（PUT 语义）
	Update(u *model.UserUpdate) error
	// Remove 删除用户信息（软删除）
	Remove(id string) error

	// RefreshLastLogin 更新用户的最后登录时间
	RefreshLastLogin(qq string, t time.Time) error

	// PreFillAvatarOSSKey 预填充用户头像的 OSS Key（图片上传不走主服务器，
	// 客户端上传完成后再确认）
	PreFillAvatarOSSKey(id string, avatarOSSKey string) error
	// ConfirmAvatarUploaded 确认用户头像已上传
	ConfirmAvatarUploaded(id string) error

	// 获取或创建用户统计信息（注意，可能总是需要一个 UPDATE ON CONFLICT 的操作）
	GetOrCreateStats(userID string) (*model.UserStats, error)
	// 按 delta 更新用户统计信息
	PatchStats(stats *model.UserStatsPatch) error

	// FromTxnCx 从上下文中获取事务，并分离出一个带事务的 UserRepo 实例
	// 如果上下文没有事务信息，则返回错误
	FromTxnCx(cx context.Context) (UserRepo, error)
}
