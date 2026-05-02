package query

// `ListChapterOpt` holds optional filters for listing chapters.
// All filter fields are pointers nil means not applied.
// Pagination is explicit in `Pagi`.
type ListChapterOpt struct {
	// `ComicId` restricts result to one comic.
	ComicId *string

	// `Pagi` specifies pagination options.
	Pagi PagiOpt
}
