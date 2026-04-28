package val

import "poprako-s/internal/domain/model/aggr"

// `AssignmentInvVal` is app-facing value object for assignment invitation.
type AssignmentInvVal struct {
	Id string `json:"id"`

	ChapterId string `json:"chapter_id"`
	InviterId string `json:"inviter_id"`

	InviteeQid string `json:"invitee_qid"`
	InvCode    string `json:"invitation_code"`

	Pending bool `json:"pending"`

	RoleMask aggr.RoleMask `json:"role_mask"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

// `ListAssignmentInvArgs` carries list args for assignment invitation.
type ListAssignmentInvArgs struct {
	ChapterId string `url:"chapter_id"`
	Pending   *bool  `url:"pending"`

	Offset int `url:"offset"`
	Limit  int `url:"limit"`
}

// `CreateAssignmentInvArgs` carries create args for assignment invitation.
type CreateAssignmentInvArgs struct {
	ChapterId  string        `json:"chapter_id"`
	InviteeQid string        `json:"invitee_qid"`
	RoleMask   aggr.RoleMask `json:"role_mask"`
}

// `CreateAssignmentInvRes` is create response of assignment invitation.
type CreateAssignmentInvRes struct {
	Id      string `json:"id"`
	InvCode string `json:"invitation_code"`
}

// `JoinAssignmentInvArgs` carries join args for assignment invitation.
type JoinAssignmentInvArgs struct {
	InvCode string `json:"invitation_code"`
}
