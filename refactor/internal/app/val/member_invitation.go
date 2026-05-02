package val

import "poprako-s/internal/domain/model/aggr"

// `MemberInvVal` is app-facing value object for one team invitation.
type MemberInvVal struct {
	Id string `json:"id"`

	InvitorId string `json:"invitor_id"`
	TeamId    string `json:"team_id"`

	InviteeQid string `json:"invitee_qid"`
	InvCode    string `json:"invitation_code"`

	Pending bool `json:"pending"`

	RoleMask aggr.RoleMask `json:"role_mask"`

	CreatedAt int64 `json:"created_at"`
}

// `ListMemberInvArgs` carries list args for team invitation listing.
type ListMemberInvArgs struct {
	TeamId string `url:"team_id"`

	Pending *bool `url:"pending"`

	Offset int `url:"offset"`
	Limit  int `url:"limit"`
}

// `CreateMemberInvArgs` carries create args for one team invitation.
type CreateMemberInvArgs struct {
	TeamId string `json:"team_id"`

	InviteeQid string `json:"invitee_qid"`

	RoleMask aggr.RoleMask `json:"role_mask"`
}

// `CreateMemberInvRes` is create response payload for one team invitation.
type CreateMemberInvRes struct {
	Id string `json:"id"`

	InvCode string `json:"invitation_code"`
}

// `MemberInvUpdArgs` carries put-style invitation role update args.
type MemberInvUpdArgs struct {
	Id string `json:"id"`

	RoleMask aggr.RoleMask `json:"role_mask"`
}
