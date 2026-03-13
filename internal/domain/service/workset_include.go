package service

import "labelplus-next-web-be/internal/value"

type WorksetListIncludeSpec struct {
	NeedTeam bool
}

// ResolveWorksetListIncludeSpec 解析 includes[]，输出 workset 列表所需的关联展开规格。
func ResolveWorksetListIncludeSpec(includes []string) WorksetListIncludeSpec {
	spec := WorksetListIncludeSpec{}

	for _, include := range includes {
		if value.CheckResourceNeccessary(include, value.IncludeTeam) {
			spec.NeedTeam = true
		}
	}

	return spec
}
