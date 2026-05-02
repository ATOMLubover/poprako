package query

// `ListWorksetOpt` holds optional filters for listing worksets.
// All fields are pointer types; nil means the field is not applied as a filter.
type ListWorksetOpt struct {
	// `TeamId` restricts the result to worksets that belong to this team.
	TeamId *string

	// `Pagi` specifies pagination options for the result.
	Pagi PagiOpt
}
