package enum

// `CommentIncl` is the include selector for team comment relation preloads.
type CommentIncl string

const (
	// `CommentInclUser` includes the `user` relation in team comment queries.
	CommentInclUser CommentIncl = "user"
)
