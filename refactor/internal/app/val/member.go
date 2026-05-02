package val

import (
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
)

// `MemberVal` is app-facing value object for one team member.
type MemberVal struct {
	Id string `json:"id"`

	UserId string `json:"user_id"`
	TeamId string `json:"team_id"`

	RoleMask aggr.RoleMask `json:"role_mask"`

	User *UserVal `json:"user,omitempty"`
	Team *TeamVal `json:"team,omitempty"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

// `CreateMemberArgs` carries create args for member.
type CreateMemberArgs struct {
	UserId string `json:"user_id"`
	TeamId string `json:"team_id"`

	RoleMask aggr.RoleMask `json:"role_mask"`
}

// `CreateMemberRes` is member creation response payload.
type CreateMemberRes struct {
	Id string `json:"id"`
}

// `ListMemberByTeamArgs` carries list args for members under one team.
type ListMemberByTeamArgs struct {
	TeamId string `url:"team_id"`

	Includes []enum.MemberIncl `url:"includes"`

	Offset int `url:"offset"`
	Limit  int `url:"limit"`
}

// `ListMyMemberArgs` carries list args for current user memberships.
type ListMyMemberArgs struct {
	Includes []enum.MemberIncl `url:"includes"`

	Offset int `url:"offset"`
	Limit  int `url:"limit"`
}

// `MemberRoleUpdArgs` carries put-style role update args for one member.
type MemberRoleUpdArgs struct {
	Id string `json:"id"`

	RoleMask aggr.RoleMask `json:"role_mask"`
}

// `JoinTeamArgs` carries invitation code for joining one team.
type JoinTeamArgs struct {
	InvCode string `json:"invitation_code"`
}
