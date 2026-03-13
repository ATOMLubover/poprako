package value

import "strings"

type GetIncludeOption string

// 在关联查询时，确定需要哪些内嵌字段
// 仅在针对类似 Assignment.Chapter 这种反向查询时有意义
// 注意，应该根据业务，允许嵌套，比如 Assignment.Chapter.Comic
const (
	IncludeUser    GetIncludeOption = "user"
	IncludeCreator GetIncludeOption = "creator"
	IncludeTeam    GetIncludeOption = "team"
	IncludeWorkset GetIncludeOption = "workset"
	IncludeComic   GetIncludeOption = "comic"
	IncludeChapter GetIncludeOption = "chapter"
)

// 在传入形如 "chapter", "comic" 时，会检查是否存在
// chapter.comic 这样的结构，这样的设计会使得 AssignmentInfo 携带 ChapterInfo
// 且 ChapterInfo 又携带 ComicInfo，形成一个树形结构
func CheckResourceNeccessary(rawOption string, options ...GetIncludeOption) bool {
	if rawOption == "" || len(options) == 0 {
		return false
	}

	segments := strings.Split(rawOption, ".")
	if len(segments) > len(options) {
		return false
	}

	for i, segment := range segments {
		if segment == "" {
			return false
		}

		if segment != string(options[i]) {
			return false
		}
	}

	return true
}
