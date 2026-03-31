package service

import "poprako-s/internal/domain/model"

// WorksetService 定义作品集领域相关的业务能力
type WorksetService interface {
	// NewCreation 根据业务参数创建 WorksetCreation 领域模型（ID 由 service 内部生成）
	// index 通常由调用方在事务中通过 Count 获取
	NewCreation(teamID string, index int, name, description string) *model.WorksetCreation
}

type worksetServiceImpl struct{}

// NewWorksetService 返回 WorksetService 的默认实现
func NewWorksetService() WorksetService {
	return &worksetServiceImpl{}
}

// NewCreation 构造一个带有 service 生成 ID 的 WorksetCreation
func (s *worksetServiceImpl) NewCreation(
	teamID string,
	index int,
	name, description string,
) *model.WorksetCreation {
	return &model.WorksetCreation{
		ID:          GenID("workset"),
		TeamID:      teamID,
		Index:       index,
		Name:        name,
		Description: description,
	}
}
