package query

// `ListCommentOpt` defines filters and pagination for team comment listing.
type ListCommentOpt struct {
	// `TeamId` is required, as comments are listed by team scope.
	TeamId string

	Pagi PagiOpt
}
