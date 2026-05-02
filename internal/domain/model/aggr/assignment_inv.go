package aggr

import "time"

// `AssignmentInv` represents one chapter assignment invitation.
type AssignmentInv struct {
	Id string

	ChapterId string
	InviterId string

	InviteeQid string
	InvCode    string

	Pending bool

	RoleMask RoleMask

	CreatedAt time.Time
	UpdatedAt time.Time
}

// `AssignmentInvCre` is the create payload for assignment invitation.
type AssignmentInvCre struct {
	Id string

	ChapterId string
	InviterId string

	InviteeQid string
	InvCode    string

	RoleMask RoleMask
}
