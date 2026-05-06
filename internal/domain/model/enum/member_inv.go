package enum

// `MemberInvIncl` is the include selector for member invitation relation preloads
type MemberInvIncl string

const (
	// `MemberInvInclInvitor` includes the `invitor` relation in invitation queries
	MemberInvInclInvitor MemberInvIncl = "invitor"

	// `MemberInvInclInvitee` includes the `invitee` relation in invitation queries
	MemberInvInclInvitee MemberInvIncl = "invitee"
)
