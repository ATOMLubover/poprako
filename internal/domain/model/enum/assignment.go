package enum

// `AssignmentIncl` is the include selector for assignment relation preloads
type AssignmentIncl string

const (
	// `AssignmentInclUser` includes the `user` relation in assignment queries
	AssignmentInclUser AssignmentIncl = "user"

	// `AssignmentInclChapter` includes the `chapter` relation in assignment queries
	AssignmentInclChapter AssignmentIncl = "chapter"

	// `AssignmentInclChapterComic` includes nested `chapter.comic` relation
	AssignmentInclChapterComic AssignmentIncl = "chapter.comic"

	// `AssignmentInclChapterComicWorkset` includes nested `chapter.comic.workset` relation
	AssignmentInclChapterComicWorkset AssignmentIncl = "chapter.comic.workset"

	// `AssignmentInclChapterComicWorksetTeam` includes nested `chapter.comic.workset.team` relation
	AssignmentInclChapterComicWorksetTeam AssignmentIncl = "chapter.comic.workset.team"

	// `AssignmentInclChapterCreator` includes nested `chapter.creator` relation
	AssignmentInclChapterCreator AssignmentIncl = "chapter.creator"

	// `AssignmentInclChapterComicCreator` includes nested `chapter.comic.creator` relation
	AssignmentInclChapterComicCreator AssignmentIncl = "chapter.comic.creator"
)
