package service

import "labelplus-next-web-be/internal/value"

// AssignmentListIncludeSpec 是 ListChapterAssignments 的关联展开规格。
type AssignmentListIncludeSpec struct {
	NeedUser bool
}

// ResolveAssignmentListIncludeSpec 解析 includes[]，输出章节分配列表所需的关联展开规格。
// 支持的 include 值：user
func ResolveAssignmentListIncludeSpec(includes []string) AssignmentListIncludeSpec {
	spec := AssignmentListIncludeSpec{}

	for _, include := range includes {
		if value.CheckResourceNeccessary(include, value.IncludeUser) {
			spec.NeedUser = true
		}
	}

	return spec
}

// AssignmentMyListIncludeSpec 是 ListMyAssignments 的关联展开规格。
type AssignmentMyListIncludeSpec struct {
	NeedChapter bool
}

// ResolveAssignmentMyListIncludeSpec 解析 includes[]，输出我的分配列表所需的关联展开规格。
// 支持的 include 值：chapter（chapter 隐含 chapter.comic）
func ResolveAssignmentMyListIncludeSpec(includes []string) AssignmentMyListIncludeSpec {
	spec := AssignmentMyListIncludeSpec{}

	for _, include := range includes {
		if value.CheckResourceNeccessary(include, value.IncludeChapter) {
			spec.NeedChapter = true
		}
	}

	return spec
}
