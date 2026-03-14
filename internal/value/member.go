package value

import (
	"errors"

	"labelplus-next-web-be/internal/domain/model"
	"labelplus-next-web-be/internal/util"
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

// MemberInfo 是成员的统一展示模型，User 和 Team 字段均为可选，仅在 includes 指定时填充。
type MemberInfo struct {
	ID string `json:"id"`

	UserID string    `json:"user_id"`
	User   *UserInfo `json:"user,omitempty"`

	TeamID string    `json:"team_id"`
	Team   *TeamInfo `json:"team,omitempty"`

	Roles model.RoleMask `json:"roles"`

	AssignedRawProviderAt *int64 `json:"assigned_raw_provider_at,omitempty"`
	AssignedTranslatorAt  *int64 `json:"assigned_translator_at,omitempty"`
	AssignedProofreaderAt *int64 `json:"assigned_proofreader_at,omitempty"`
	AssignedTypesetterAt  *int64 `json:"assigned_typesetter_at,omitempty"`
	AssignedReviewerAt    *int64 `json:"assigned_reviewer_at,omitempty"`
	AssignedPublisherAt   *int64 `json:"assigned_publisher_at,omitempty"`
	AssignedAdminAt       *int64 `json:"assigned_admin_at,omitempty"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

func NewMemberInfoFromModel(m model.MemberWithInfo) MemberInfo {
	result := MemberInfo{
		ID:                    m.ID,
		UserID:                m.UserID,
		TeamID:                m.TeamID,
		Roles:                 model.MaskRoles(m.Roles()),
		AssignedRawProviderAt: util.ToUnixPtr(m.AssignedRawProviderAt),
		AssignedTranslatorAt:  util.ToUnixPtr(m.AssignedTranslatorAt),
		AssignedProofreaderAt: util.ToUnixPtr(m.AssignedProofreaderAt),
		AssignedTypesetterAt:  util.ToUnixPtr(m.AssignedTypesetterAt),
		AssignedReviewerAt:    util.ToUnixPtr(m.AssignedReviewerAt),
		AssignedPublisherAt:   util.ToUnixPtr(m.AssignedPublisherAt),
		AssignedAdminAt:       util.ToUnixPtr(m.AssignedAdminAt),
		CreatedAt:             m.CreatedAt.UnixMilli(),
		UpdatedAt:             m.UpdatedAt.UnixMilli(),
	}

	if m.User != nil {
		userInfo := NewUserInfoFromModel(*m.User, "")
		result.User = &userInfo
	}

	if m.Team != nil {
		teamInfo := NewTeamInfoFromModel(*m.Team, "")
		result.Team = &teamInfo
	}

	return result
}

type ListTeamMemberArgs struct {
	TeamID   string   `url:"team_id"`
	Includes []string `url:"includes[]"`
	PaginationParams
}

type ListMyMemberArgs struct {
	Includes []string `url:"includes[]"`
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

func (lmma *ListMyMemberArgs) Validate() error {
	if lmma == nil {
		return errors.New("参数不能为空")
	}

	if err := lmma.PaginationParams.Validate(); err != nil {
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
