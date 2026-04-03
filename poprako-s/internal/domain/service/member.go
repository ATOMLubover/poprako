package service

import (
	"errors"
	"time"

	"poprako-s/internal/domain/model"
	"poprako-s/internal/domain/repo"
)

// MemberService 定义成员领域相关的业务能力
type MemberService interface {
	// NewCreation 根据当前用户、目标用户 ID、团队 ID 和角色列表创建 MemberCreation 领域模型
	// 仅超级管理员可以直接创建成员
	NewCreation(
		currUser *model.UserInfo,
		userID string,
		teamID string,
		roles ...model.Role,
	) (*model.MemberCreation, error)

	// NewCreationFromInvitation 根据邀请信息创建 MemberCreation 领域模型
	// 角色由邀请信息决定，不需要超级管理员权限
	NewCreationFromInvitation(
		currUser *model.UserInfo,
		i *model.InvitationInfo,
	) (*model.MemberCreation, error)

	// NewUpdate 根据当前操作用户权限和目标角色掩码生成 MemberUpdate
	// 已有角色保留原时间戳，新增角色使用当前时间，移除的角色清空时间戳
	// 仅团队管理员可以更新成员角色
	NewUpdate(
		mr repo.MemberRepo,
		currUserID string,
		teamID string,
		targetRoles model.RoleMask,
	) (*model.MemberUpdate, error)
}

// memberServiceImpl 是 MemberService 的具体实现，无内禀状态
type memberServiceImpl struct{}

// NewMemberService 返回 MemberService 的默认实现
func NewMemberService() MemberService {
	// 返回无状态实现
	return &memberServiceImpl{}
}

// NewCreation 构造一个带 service 生成 ID 的 MemberCreation
func (s *memberServiceImpl) NewCreation(
	currUser *model.UserInfo,
	userID string,
	teamID string,
	roles ...model.Role,
) (*model.MemberCreation, error) {
	// 仅超级管理员可以直接创建成员
	if !currUser.IsSuperAdmin {
		// 返回权限不足错误
		return nil, errors.New("仅超级管理员可以直接创建成员")
	}

	// 构造成员创建载荷
	c := &model.MemberCreation{
		ID:     GenID("member"),
		UserID: userID,
		TeamID: teamID,
	}

	// 将角色列表映射到布尔字段
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

	// 返回构造结果
	return c, nil
}

// NewCreationFromInvitation 根据邀请信息构造 MemberCreation，角色由邀请决定
func (s *memberServiceImpl) NewCreationFromInvitation(
	currUser *model.UserInfo,
	i *model.InvitationInfo,
) (*model.MemberCreation, error) {
	// 构造成员创建载荷，直接映射邀请中指定的角色
	c := &model.MemberCreation{
		ID:              GenID("member"),
		UserID:          currUser.ID,
		TeamID:          i.TeamID,
		ToBeRawProvider: i.ToBeRawProvider,
		ToBeTranslator:  i.ToBeTranslator,
		ToBeProofreader: i.ToBeProofreader,
		ToBeTypesetter:  i.ToBeTypesetter,
		ToBeReviewer:    i.ToBeReviewer,
		ToBePublisher:   i.ToBePublisher,
		ToBeAdmin:       i.ToBeAdmin,
	}

	// 返回构造结果
	return c, nil
}

// NewUpdate 根据目标角色掩码和当前操作用户权限构造 MemberUpdate
func (s *memberServiceImpl) NewUpdate(
	mr repo.MemberRepo,
	currUserID string,
	teamID string,
	targetRoles model.RoleMask,
) (*model.MemberUpdate, error) {
	// 查询当前操作用户的成员记录，用于鉴权
	currMem, err := mr.Get(model.MemberQueryOpt{
		UserID: &currUserID,
		TeamID: &teamID,
	})
	if err != nil || !currMem.HasAnyRole(model.RoleAdmin) {
		// 返回权限不足错误
		return nil, errors.New("仅团队管理员可以更新成员角色")
	}

	// 记录当前时间，用于新增角色的时间戳
	now := time.Now()

	// resolve 决定某个角色字段的时间戳：
	// - 目标掩码中不含该角色 → nil（清除）
	// - 目标掩码含该角色且原有时间戳存在 → 保留原时间戳
	// - 目标掩码含该角色但无原时间戳 → 使用当前时间
	resolve := func(currMemt *time.Time, role model.Role) *time.Time {
		if targetRoles&model.RoleMask(role) == 0 {
			return nil
		}

		if currMemt != nil {
			t := *currMemt
			return &t
		}

		t := now

		return &t
	}

	// 构造并返回 MemberUpdate，保留已有时间戳
	return &model.MemberUpdate{
		ID:                    teamID,
		AssignedRawProviderAt: resolve(currMem.AssignedRawProviderAt, model.RoleRawProvider),
		AssignedTranslatorAt:  resolve(currMem.AssignedTranslatorAt, model.RoleTranslator),
		AssignedProofreaderAt: resolve(currMem.AssignedProofreaderAt, model.RoleProofreader),
		AssignedTypesetterAt:  resolve(currMem.AssignedTypesetterAt, model.RoleTypesetter),
		AssignedReviewerAt:    resolve(currMem.AssignedReviewerAt, model.RoleReviewer),
		AssignedPublisherAt:   resolve(currMem.AssignedPublisherAt, model.RolePublisher),
		AssignedAdminAt:       resolve(currMem.AssignedAdminAt, model.RoleAdmin),
	}, nil
}
