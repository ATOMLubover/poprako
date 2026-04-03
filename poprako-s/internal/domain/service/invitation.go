package service

import (
	"errors"

	"poprako-s/internal/domain/model"
	"poprako-s/internal/domain/repo"
)

// InvitationService 定义邀请领域相关的业务能力
type InvitationService interface {
	// NewCreation 根据业务参数创建 InvitationCreation 领域模型（ID 和邀请码由 service 内部生成）
	// 仅团队管理员可以创建邀请
	NewCreation(
		mr repo.MemberRepo,
		invitorID string,
		targetTeamID string,
		inviteeQQ string,
		roles ...model.Role,
	) (*model.InvitationCreation, error)
	// GenInvitationCode 生成一个 6 位邀请码
	GenInvitationCode() string
}

// invitationServiceImpl 是 InvitationService 的具体实现 无内禀状态
type invitationServiceImpl struct{}

// NewInvitationService 返回 InvitationService 的默认实现
func NewInvitationService() InvitationService {
	// 返回无状态实现
	return &invitationServiceImpl{}
}

// NewCreation 构造一个带有 service 生成 ID 和邀请码的 InvitationCreation
// 仅团队管理员可以创建邀请
func (s *invitationServiceImpl) NewCreation(
	mr repo.MemberRepo,
	invitorID string,
	targetTeamID string,
	inviteeQQ string,
	roles ...model.Role,
) (*model.InvitationCreation, error) {
	// 查询邀请发起人在团队中的成员记录
	member, err := mr.Get(model.MemberQueryOpt{
		UserID: &invitorID,
		TeamID: &targetTeamID,
	})
	if err != nil || !member.HasAnyRole(model.RoleAdmin) {
		// 返回权限错误
		return nil, errors.New("仅团队管理员可以创建邀请")
	}

	// 构造邀请创建载荷
	c := &model.InvitationCreation{
		ID:             GenID("invitation"),
		InvitorID:      invitorID,
		TargetTeamID:   targetTeamID,
		InviteeQQ:      inviteeQQ,
		InvitationCode: s.GenInvitationCode(),
	}

	// 映射角色到布尔字段
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

	// 返回创建载荷
	return c, nil
}

// GenInvitationCode 取 UUID 后 6 位作为邀请码
func (s *invitationServiceImpl) GenInvitationCode() string {
	// 生成带前缀 ID
	id := GenID("invitation")

	// 截取后 6 位作为邀请码
	if len(id) >= 6 {
		return id[len(id)-6:]
	}

	// 长度不足时直接返回
	return id
}
