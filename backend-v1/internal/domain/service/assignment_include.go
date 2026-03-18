package service

import "labelplus-next-web-be/internal/value"

// AssignmentListIncludeSpec 是 assignment 列表接口的关联展开规格。
type AssignmentListIncludeSpec struct {
	NeedUser           bool
	NeedChapter        bool
	NeedChapterComic   bool
	NeedChapterCreator bool
}

// ResolveAssignmentListIncludeSpec 解析 includes[]，输出章节分配列表所需的关联展开规格。
// 支持的 include 值：user, chapter, chapter.comic, chapter.creator
func ResolveAssignmentListIncludeSpec(includes []string) AssignmentListIncludeSpec {
	spec := AssignmentListIncludeSpec{}

	for _, include := range includes {
		if value.CheckResourceNeccessary(include, value.IncludeUser) {
			spec.NeedUser = true
		}

		if value.CheckResourceNeccessary(include, value.IncludeChapter) {
			spec.NeedChapter = true
		}

		if value.CheckResourceNeccessary(include, value.IncludeChapter, value.IncludeComic) {
			spec.NeedChapter = true
			spec.NeedChapterComic = true
		}

		if value.CheckResourceNeccessary(include, value.IncludeChapter, value.IncludeCreator) {
			spec.NeedChapter = true
			spec.NeedChapterCreator = true
		}
	}

	return spec
}

// AssignmentMyListIncludeSpec 是 ListMyAssignments 的关联展开规格。
type AssignmentMyListIncludeSpec = AssignmentListIncludeSpec

// ResolveAssignmentMyListIncludeSpec 解析 includes[]，输出我的分配列表所需的关联展开规格。
// 支持的 include 值：chapter, chapter.comic, chapter.creator
func ResolveAssignmentMyListIncludeSpec(includes []string) AssignmentMyListIncludeSpec {
	spec := AssignmentMyListIncludeSpec{}

	for _, include := range includes {
		if value.CheckResourceNeccessary(include, value.IncludeChapter) {
			spec.NeedChapter = true
		}

		if value.CheckResourceNeccessary(include, value.IncludeChapter, value.IncludeComic) {
			spec.NeedChapter = true
			spec.NeedChapterComic = true
		}

		if value.CheckResourceNeccessary(include, value.IncludeChapter, value.IncludeCreator) {
			spec.NeedChapter = true
			spec.NeedChapterCreator = true
		}
	}

	return spec
}
