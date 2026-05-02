package aggr

import (
	"time"
)

type MemberInv struct {
	Id string

	InvitorId string
	Invitor   *User
	TeamId    string

	// NOTE: as invitee may not be registered, we use `InviteeQid` instead of `InviteeId`.
	InviteeQid string
	Invitee    *User

	InvCode string
	Pending bool

	RoleMask RoleMask

	CreatedAt time.Time
}

type MemberInvCre struct {
	Id string

	InvitorId string
	TeamId    string

	InviteeQid string
	InvCode    string

	RoleMask RoleMask
}

// `MemberInvUpd` holds mutable fields for one member invitation put update.
type MemberInvUpd struct {
	Id string

	RoleMask RoleMask
}
