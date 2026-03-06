package value

import (
	"errors"

	"labelplus-next-web-be/internal/domain/model"

	"go.uber.org/zap"
)

type UpdateMemberRoleArgs struct {
	ID string `json:"id"`

	Roles model.RoleMask `json:"roles"`
}

func (uma *UpdateMemberRoleArgs) Validate() error {
	if uma == nil {
		return errors.New("参数不能为空")
	}

	if uma.ID == "" {
		return errors.New("成员 ID 不能为空")
	}

	if uma.Roles == 0 {
		return errors.New("分工角色不能为空")
	}

	return nil
}

// MemberProfile 不仅包含成员的汉化组信息，
// 还包含成员的用户信息（如昵称、头像等）
type MemberProfile struct {
	UserInfo
	Roles model.RoleMask `json:"roles"`
}

func NewMemberProfile(userInfo UserInfo, roles ...model.RoleFlag) *MemberProfile {
	return &MemberProfile{
		UserInfo: userInfo,
		Roles:    model.MaskRoles(roles),
	}
}

func NewMemberProfileFromModel(mp *model.MemberProfile) *MemberProfile {
	if mp == nil {
		zap.L().Warn("NewMemberProfileFromModel: mp 为空")
		return nil
	}

	userInfo := NewUserInfoFromModel(mp.UserInfo)
	if userInfo == nil {
		zap.L().Warn("NewMemberProfileFromModel: 用户信息为空")
		return nil
	}

	return NewMemberProfile(*userInfo, mp.Roles()...)
}

type ListTeamMemberArgs struct {
	TeamID string `url:"team_id"`
	PaginationParams
}

func (ltma *ListTeamMemberArgs) Validate() error {
	if ltma == nil {
		return errors.New("参数不能为空")
	}

	if ltma.TeamID == "" {
		return errors.New("汉化组 ID 不能为空")
	}

	if err := ltma.PaginationParams.Validate(); err != nil {
		return errors.New("分页参数无效: " + err.Error())
	}

	return nil
}

type CreateMemberArgs struct {
	UserID string         `json:"user_id"`
	TeamID string         `json:"team_id"`
	Roles  model.RoleMask `json:"roles"`
}

func (cma *CreateMemberArgs) Validate() error {
	if cma == nil {
		return errors.New("参数不能为空")
	}

	if cma.UserID == "" {
		return errors.New("用户 ID 不能为空")
	}

	if cma.TeamID == "" {
		return errors.New("汉化组 ID 不能为空")
	}

	if cma.Roles == 0 {
		return errors.New("分工角色不能为空")
	}

	return nil
}

type CreateMemberResult struct {
	MemberID string `json:"member_id"`
}

func NewCreateMemberResult(memberID string) *CreateMemberResult {
	return &CreateMemberResult{
		MemberID: memberID,
	}
}

type JoinTeamArgs struct {
	InvitationCode string `json:"invitation_code"`
}

func (jta *JoinTeamArgs) Validate() error {
	if jta == nil {
		return errors.New("参数不能为空")
	}

	if jta.InvitationCode == "" {
		return errors.New("邀请码不能为空")
	}

	return nil
}
