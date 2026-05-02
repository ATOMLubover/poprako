package query

// `ListTeamOpt` defines filters and pagination for team listing.
type ListTeamOpt struct {
	Id *string

	Pagi PagiOpt
}
