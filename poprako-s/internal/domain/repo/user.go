package repo

import (
	"time"

	"poprako-s/internal/domain/model"
)

// UserRepo 是用户仓库的接口
type UserRepo interface {
	// 根据 QQ 号获取用户凭证信息，这是一个特殊场景
	GetCredsByQQ(qq string) (*model.UserCreds, error)
	// 根据用户 ID 获取用户信息
	GetByID(id string) (*model.UserInfo, error)

	// 注册用户
	Reg(reg *model.UserReg) (*model.UserInfo, error)
	// 更新用户信息
	Update(update *model.UserUpdate) error
	// 删除用户信息（硬删除）
	Delete(id string) error

	// RefreshLastLogin 更新用户的最后登录时间
	RefreshLastLogin(qq string, t time.Time) error

	// 预填充用户头像信息（因为采用的是图片上传不走主服务器的逻辑
	// 所以需要客户端自行确认头像已上传）
	PreFillAvatarOSSKey(id string, avatarOSSKey string) error
	// 确认用户头像已上传
	ConfirmAvatarUploaded(id string) error
}
