package service

import (
	"errors"

	"poprako-s/internal/domain/model"
	"poprako-s/internal/domain/repo"
)

// WorksetService 定义作品集领域相关的业务能力
type WorksetService interface {
	// NewCreation 根据业务参数创建 WorksetCreation 领域模型（ID 由 service 内部生成）
	// index 通常由调用方在事务中通过 Count 获取
	// 仅团队管理员可以创建作品集
	NewCreation(
		mr repo.MemberRepo,
		currUserID string,
		teamID string,
		index int,
		name string,
		description string,
	) (*model.WorksetCreation, error)
}

// worksetServiceImpl 是 WorksetService 的具体实现 无内禀状态
type worksetServiceImpl struct{}

// NewWorksetService 返回 WorksetService 的默认实现
func NewWorksetService() WorksetService {
	// 返回无状态实现
	return &worksetServiceImpl{}
}

// NewCreation 构造一个带有 service 生成 ID 的 WorksetCreation
// 仅团队管理员可以创建作品集
func (s *worksetServiceImpl) NewCreation(
	mr repo.MemberRepo,
	currUserID string,
	teamID string,
	index int,
	name string,
	description string,
) (*model.WorksetCreation, error) {
	// 查询当前用户在团队中的成员记录 用于鉴权
	member, err := mr.Get(model.MemberQueryOpt{
		UserID: &currUserID,
		TeamID: &teamID,
	})
	if err != nil || !member.HasAnyRole(model.RoleAdmin) {
		// 返回权限错误
		return nil, errors.New("仅团队管理员可以创建作品集")
	}

	// 返回创建载荷
	return &model.WorksetCreation{
		ID:          GenID("workset"),
		TeamID:      teamID,
		Index:       index,
		Name:        name,
		Description: description,
	}, nil
}
