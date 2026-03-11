package value

import (
	"errors"

	"labelplus-next-web-be/internal/domain/model"
	"labelplus-next-web-be/internal/util"

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

func NewMemberProfile(userInfo UserInfo, roles ...model.RoleFlag) MemberProfile {
	return MemberProfile{
		UserInfo: userInfo,
		Roles:    model.MaskRoles(roles),
	}
}

// FIXME: 有没有更好的方法自动生成 avatarURL
func NewMemberProfileFromModel(mp model.MemberWithUserInfo, avatarURL string) MemberProfile {
	if mp.UserInfo.ID == "" {
		zap.L().Warn("NewMemberProfileFromModel: mp 为空或用户信息为空")
		return MemberProfile{}
	}

	userInfo := NewUserInfoFromModel(mp.UserInfo, avatarURL)

	return NewMemberProfile(userInfo, mp.Roles()...)
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

func NewCreateMemberResult(memberID string) CreateMemberResult {
	return CreateMemberResult{
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

type MemberWithTeamInfo struct {
	ID string `json:"id"`

	UserID string   `json:"user_id"`
	Team   TeamInfo `json:"team"`

	AssignedRawProviderAt *int64 `json:"assigned_raw_provider_at,omitempty"`
	AssignedTranslatorAt  *int64 `json:"assigned_translator_at,omitempty"`
	AssignedProofreaderAt *int64 `json:"assigned_proofreader_at,omitempty"`
	AssignedTypesetterAt  *int64 `json:"assigned_typesetter_at,omitempty"`
	AssignedReviewerAt    *int64 `json:"assigned_reviewer_at,omitempty"`
	AssignedPublishererAt *int64 `json:"assigned_publisher_at,omitempty"`
	AssignedAdminAt       *int64 `json:"assigned_admin_at,omitempty"`
}

func NewMemberWithTeamInfoFromModel(m model.MemberWithTeamInfo, avatarURL string) MemberWithTeamInfo {
	return MemberWithTeamInfo{
		ID: m.ID,

		UserID: m.UserID,
		Team:   NewTeamInfoFromModel(m.Team, avatarURL),

		AssignedRawProviderAt: util.ToUnixPtr(m.AssignedRawProviderAt),
		AssignedTranslatorAt:  util.ToUnixPtr(m.AssignedTranslatorAt),
		AssignedProofreaderAt: util.ToUnixPtr(m.AssignedProofreaderAt),
		AssignedTypesetterAt:  util.ToUnixPtr(m.AssignedTypesetterAt),
		AssignedReviewerAt:    util.ToUnixPtr(m.AssignedReviewerAt),
		AssignedPublishererAt: util.ToUnixPtr(m.AssignedPublisherAt),
		AssignedAdminAt:       util.ToUnixPtr(m.AssignedAdminAt),
	}
}
