package query

// `ListAssignmentInvOpt` defines list filters for assignment invitation.
type ListAssignmentInvOpt struct {
	ChapterId string
	Pending   *bool

	Pagi PagiOpt
}
