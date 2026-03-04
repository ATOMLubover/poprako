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
