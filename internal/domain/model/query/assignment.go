package query

// `ListAssignmentOpt` defines list filters for assignment.
type ListAssignmentOpt struct {
	ChapterId *string
	UserId    *string

	Pagi PagiOpt
}
