package service

import "labelplus-next-web-be/internal/value"

type ChapterListIncludeSpec struct {
	NeedCreator bool
}

// ResolveChapterListIncludeSpec 解析 includes[]，输出章节列表所需的关联展开规格。
func ResolveChapterListIncludeSpec(includes []string) ChapterListIncludeSpec {
	spec := ChapterListIncludeSpec{}

	for _, include := range includes {
		if value.CheckResourceNeccessary(include, value.IncludeCreator) {
			spec.NeedCreator = true
		}
	}

	return spec
}
