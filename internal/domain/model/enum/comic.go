package enum

// `ComicIncl` is the include selector for comic relation preloads
type ComicIncl string

const (
	// `ComicInclWorkset` includes the `workset` relation in comic queries
	ComicInclWorkset ComicIncl = "workset"

	// `ComicInclWorksetTeam` includes nested `workset.team` relation in comic queries
	ComicInclWorksetTeam ComicIncl = "workset.team"

	// `ComicInclCreator` includes the `creator` relation in comic queries
	ComicInclCreator ComicIncl = "creator"
)
