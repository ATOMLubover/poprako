package service

import (
	"strings"

	"poprako-s/internal/domain/model"
)

// TeamService 定义团队领域相关的业务能力
type TeamService interface {
	// NewCreation 根据业务参数创建 TeamCreation 领域模型（ID 由 service 内部生成）
	NewCreation(name, description string) *model.TeamCreation
	// GenAvatarOSSKey 根据团队 ID 生成头像的 OSS Key
	GenAvatarOSSKey(teamID string) string
}

type teamServiceImpl struct{}

// NewTeamService 返回 TeamService 的默认实现
func NewTeamService() TeamService {
	return &teamServiceImpl{}
}

// NewCreation 构造一个带有 service 生成 ID 的 TeamCreation
func (s *teamServiceImpl) NewCreation(name, description string) *model.TeamCreation {
	return &model.TeamCreation{
		ID:          GenID("team"),
		Name:        name,
		Description: description,
	}
}

// GenAvatarOSSKey 以固定前缀拼接团队 ID 作为头像的 OSS Key
func (s *teamServiceImpl) GenAvatarOSSKey(teamID string) string {
	return strings.Join([]string{"team-avatar", teamID}, "_")
}
