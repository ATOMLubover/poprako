package service

import (
	"errors"

	"poprako-s/internal/domain/model"
)

// TeamService 定义汉化组领域相关的业务能力
type TeamService interface {
	// NewCreation 根据业务参数创建 TeamCreation 领域模型（ID 由 service 内部生成）
	// 仅超级管理员可以创建汉化组
	NewCreation(
		currUser *model.UserInfo,
		name string,
		desc string,
	) (*model.TeamCreation, error)
	// GenAvatarOSSKey 根据汉化组 ID 生成头像的 OSS Key
	GenAvatarOSSKey(
		teamID string,
	) string
}

// teamServiceImpl 是 TeamService 的具体实现 无内禀状态
type teamServiceImpl struct{}

// NewTeamService 返回 TeamService 的默认实现
func NewTeamService() TeamService {
	return &teamServiceImpl{}
}

// NewCreation 构造一个带有 service 生成 ID 的 TeamCreation
// 仅超级管理员可以创建汉化组
func (s *teamServiceImpl) NewCreation(
	currUser *model.UserInfo,
	name string,
	desc string,
) (*model.TeamCreation, error) {
	// 校验超级管理员权限
	if !currUser.IsSuperAdmin {
		return nil, errors.New("仅超级管理员可以创建汉化组")
	}

	// 返回创建载荷
	return &model.TeamCreation{
		ID:   GenID("team"),
		Name: name,
		Desc: desc,
	}, nil
}

// GenAvatarOSSKey 以目录式命名生成汉化组头像的 OSS Key，格式为 team_{teamID}/avatar
func (s *teamServiceImpl) GenAvatarOSSKey(
	teamID string,
) string {
	// 返回目录式对象 Key，避免跨实体命名冲突
	return "team_" + teamID + "/avatar"
}
