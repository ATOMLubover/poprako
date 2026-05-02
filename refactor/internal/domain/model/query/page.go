package query

// `ListPageOpt` defines filters and pagination for page listing.
type ListPageOpt struct {
	// `ChapterId` restricts pages to one chapter.
	ChapterId *string

	// `Pagi` specifies explicit pagination options.
	Pagi PagiOpt
}
