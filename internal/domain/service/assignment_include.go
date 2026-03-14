package service

import "labelplus-next-web-be/internal/value"

type AssignmentListIncludeSpec struct {
	NeedUser bool
}

// ResolveAssignmentListIncludeSpec 解析 includes[]，输出章节分配列表所需的关联展开规格。
func ResolveAssignmentListIncludeSpec(includes []string) AssignmentListIncludeSpec {
	spec := AssignmentListIncludeSpec{}

	for _, include := range includes {
		if value.CheckResourceNeccessary(include, value.IncludeUser) {
			spec.NeedUser = true
		}
	}

	return spec
}
