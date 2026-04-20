package service

import (
	"errors"

	"poprako-s/internal/domain/model"
	"poprako-s/internal/domain/repo"
)

// MemberInvitationService 定义成员邀请领域相关的业务能力
type MemberInvitationService interface {
	// NewCreation 根据业务参数创建 InvitationCreation 领域模型（ID 和邀请码由 service 内部生成）
	// 仅汉化组管理员可以创建邀请
	NewCreation(
		mr repo.MemberRepo,
		invitorID string,
		targetTeamID string,
		inviteeQQ string,
		roles ...model.Role,
	) (*model.MemberInvitationCreation, error)
	// GenInvCode 生成一个 6 位邀请码
	GenInvCode() string
}

// memberInvitationServiceImpl 是 InvitationService 的具体实现 无内禀状态
type memberInvitationServiceImpl struct{}

// NewMemberInvitationService 返回 InvitationService 的默认实现
func NewMemberInvitationService() MemberInvitationService {
	// 返回无状态实现
	return &memberInvitationServiceImpl{}
}

// NewCreation 构造一个带有 service 生成 ID 和邀请码的 InvitationCreation
// 仅汉化组管理员可以创建邀请
func (s *memberInvitationServiceImpl) NewCreation(
	mr repo.MemberRepo,
	invitorID string,
	targetTeamID string,
	inviteeQQ string,
	roles ...model.Role,
) (*model.MemberInvitationCreation, error) {
	// 查询邀请发起人在汉化组中的成员记录
	member, err := mr.Get(model.MemberQueryOpt{
		UserID: &invitorID,
		TeamID: &targetTeamID,
	})
	if err != nil || !member.HasAnyRole(model.RoleAdmin) {
		// 返回权限错误
		return nil, errors.New("仅汉化组管理员可以创建邀请")
	}

	// 构造邀请创建载荷
	c := &model.MemberInvitationCreation{
		ID:             GenID("invitation"),
		InvitorID:      invitorID,
		TargetTeamID:   targetTeamID,
		InviteeQQ:      inviteeQQ,
		InvitationCode: s.GenInvCode(),
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
		case model.RoleRedrawer:
			return nil, errors.New("成员邀请暂不支持美工角色")
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

// GenInvCode 取 UUID 后 6 位作为邀请码
func (s *memberInvitationServiceImpl) GenInvCode() string {
	// 生成带前缀 ID
	id := GenID("invitation")

	// 截取后 6 位作为邀请码
	if len(id) >= 6 {
		return id[len(id)-6:]
	}

	// 长度不足时直接返回
	return id
}
