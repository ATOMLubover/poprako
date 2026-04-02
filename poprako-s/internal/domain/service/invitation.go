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
	NewCreation(invitorID, targetTeamID, inviteeQQ string, roles ...model.Role) (*model.InvitationCreation, error)
	// GenInvitationCode 生成一个 6 位邀请码
	GenInvitationCode() string
}

type invitationServiceImpl struct {
	memberRepo repo.MemberRepo
}

// NewInvitationService 返回 InvitationService 的默认实现
func NewInvitationService(memberRepo repo.MemberRepo) InvitationService {
	return &invitationServiceImpl{memberRepo: memberRepo}
}

// NewCreation 构造一个带有 service 生成 ID 和邀请码的 InvitationCreation
// 仅团队管理员可以创建邀请
func (s *invitationServiceImpl) NewCreation(
	invitorID, targetTeamID, inviteeQQ string,
	roles ...model.Role,
) (*model.InvitationCreation, error) {
	member, err := s.memberRepo.Get(model.MemberQueryOpt{
		UserID: &invitorID,
		TeamID: &targetTeamID,
	})
	if err != nil {
		return nil, errors.New("仅团队管理员可以创建邀请")
	}
	if !member.HasAnyRole(model.RoleAdmin) {
		return nil, errors.New("仅团队管理员可以创建邀请")
	}

	c := &model.InvitationCreation{
		ID:             GenID("invitation"),
		InvitorID:      invitorID,
		TargetTeamID:   targetTeamID,
		InviteeQQ:      inviteeQQ,
		InvitationCode: s.GenInvitationCode(),
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

// GenInvitationCode 取 UUID 后 6 位作为邀请码
func (s *invitationServiceImpl) GenInvitationCode() string {
	id := GenID("invitation")
	if len(id) >= 6 {
		return id[len(id)-6:]
	}

	return id
}
