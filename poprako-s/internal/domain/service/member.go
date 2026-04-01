package service

import (
	"errors"
	"time"

	"poprako-s/internal/domain/model"
	"poprako-s/internal/domain/repo"
)

// MemberService 定义成员领域相关的业务能力
type MemberService interface {
	// NewCreation 根据用户 ID、团队 ID 和角色列表创建 MemberCreation 领域模型
	// 仅超级管理员可以直接创建成员
	NewCreation(currUser *model.UserInfo, userID, teamID string, roles ...model.Role) (*model.MemberCreation, error)
	// NewUpdate 根据当前成员信息和目标角色掩码生成 MemberUpdate
	// 已有角色保留原时间戳，新增角色使用当前时间，移除的角色清空时间戳
	// 仅团队管理员可以更新成员角色
	NewUpdate(currUserID string, id string, current *model.MemberInfo, targetRoles model.RoleMask) (*model.MemberUpdate, error)
}

type memberServiceImpl struct {
	memberRepo repo.MemberRepo
}

// NewMemberService 返回 MemberService 的默认实现
func NewMemberService(memberRepo repo.MemberRepo) MemberService {
	return &memberServiceImpl{memberRepo: memberRepo}
}

// NewCreation 构造一个带有 service 生成 ID 的 MemberCreation
// 仅超级管理员可以直接创建成员
func (s *memberServiceImpl) NewCreation(currUser *model.UserInfo, userID, teamID string, roles ...model.Role) (*model.MemberCreation, error) {
	if !currUser.IsSuperAdmin {
		return nil, errors.New("仅超级管理员可以直接创建成员")
	}
	c := &model.MemberCreation{
		ID:     GenID("member"),
		UserID: userID,
		TeamID: teamID,
	}

	for _, role := range roles {
		switch role {
		case model.RoleRawProvider:
			c.ToBeRawProvider = true
		case model.RoleTranslator:
			c.ToBeTranslator = true
		case model.RoleProofreader:
			c.ToBeProofreader = true
		case model.RoleTypesetter:
			c.ToBeTypesetter = true
		case model.RoleReviewer:
			c.ToBeReviewer = true
		case model.RolePublisher:
			c.ToBePublisher = true
		case model.RoleAdmin:
			c.ToBeAdmin = true
		}
	}

	return c, nil
}

// NewUpdate 根据目标角色掩码和当前成员信息构造 MemberUpdate
// 仅团队管理员可以更新成员角色
func (s *memberServiceImpl) NewUpdate(
	currUserID string,
	id string,
	current *model.MemberInfo,
	targetRoles model.RoleMask,
) (*model.MemberUpdate, error) {
	member, err := s.memberRepo.Get(model.MemberQueryOpt{
		UserID: &currUserID,
		TeamID: &current.TeamID,
	})
	if err != nil || !member.HasAnyRole(model.RoleAdmin) {
		return nil, errors.New("仅团队管理员可以更新成员角色")
	}
	now := time.Now()

	resolve := func(currentAt *time.Time, role model.Role) *time.Time {
		if targetRoles&model.RoleMask(role) == 0 {
			return nil
		}
		if currentAt != nil {
			t := *currentAt
			return &t
		}
		t := now
		return &t
	}

	return &model.MemberUpdate{
		ID:                    id,
		AssignedRawProviderAt: resolve(current.AssignedRawProviderAt, model.RoleRawProvider),
		AssignedTranslatorAt:  resolve(current.AssignedTranslatorAt, model.RoleTranslator),
		AssignedProofreaderAt: resolve(current.AssignedProofreaderAt, model.RoleProofreader),
		AssignedTypesetterAt:  resolve(current.AssignedTypesetterAt, model.RoleTypesetter),
		AssignedReviewerAt:    resolve(current.AssignedReviewerAt, model.RoleReviewer),
		AssignedPublisherAt:   resolve(current.AssignedPublisherAt, model.RolePublisher),
		AssignedAdminAt:       resolve(current.AssignedAdminAt, model.RoleAdmin),
	}, nil
}
