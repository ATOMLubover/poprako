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
	NewCreation(currUserID, teamID string, index int, name, description string) (*model.WorksetCreation, error)
}

type worksetServiceImpl struct {
	memberRepo repo.MemberRepo
}

// NewWorksetService 返回 WorksetService 的默认实现
func NewWorksetService(memberRepo repo.MemberRepo) WorksetService {
	return &worksetServiceImpl{memberRepo: memberRepo}
}

// NewCreation 构造一个带有 service 生成 ID 的 WorksetCreation
// 仅团队管理员可以创建作品集
func (s *worksetServiceImpl) NewCreation(
	currUserID, teamID string,
	index int,
	name, description string,
) (*model.WorksetCreation, error) {
	member, err := s.memberRepo.Get(model.MemberQueryOpt{
		UserID: &currUserID,
		TeamID: &teamID,
	})
	if err != nil || !member.HasAnyRole(model.RoleAdmin) {
		return nil, errors.New("仅团队管理员可以创建作品集")
	}

	return &model.WorksetCreation{
		ID:          GenID("workset"),
		TeamID:      teamID,
		Index:       index,
		Name:        name,
		Description: description,
	}, nil
}
