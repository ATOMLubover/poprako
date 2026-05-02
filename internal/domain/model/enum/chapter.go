package enum

// `ChapterIncl` is the include selector for chapter relation preloads
type ChapterIncl string

const (
	// `ChapterInclComic` includes `comic` relation in chapter queries
	ChapterInclComic ChapterIncl = "comic"
)
