package enum

// `ComicIncl` is the include selector for comic relation preloads
type ComicIncl string

const (
	// `ComicInclWorkset` includes the `workset` relation in comic queries
	ComicInclWorkset ComicIncl = "workset"
)
