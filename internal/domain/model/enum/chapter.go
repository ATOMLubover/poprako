package enum

// `ChapterIncl` is the include selector for chapter relation preloads
type ChapterIncl string

const (
	// `ChapterInclComic` includes `comic` relation in chapter queries
	ChapterInclComic ChapterIncl = "comic"

	// `ChapterInclComicWorkset` includes nested `comic.workset` relation
	ChapterInclComicWorkset ChapterIncl = "comic.workset"

	// `ChapterInclComicWorksetTeam` includes nested `comic.workset.team` relation
	ChapterInclComicWorksetTeam ChapterIncl = "comic.workset.team"

	// `ChapterInclComicCreator` includes nested `comic.creator` relation
	ChapterInclComicCreator ChapterIncl = "comic.creator"

	// `ChapterInclCreator` includes `creator` relation in chapter queries
	ChapterInclCreator ChapterIncl = "creator"
)
