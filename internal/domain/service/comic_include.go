package service

import "labelplus-next-web-be/internal/value"

type ComicListIncludeSpec struct {
	NeedWorkset bool
}

// ResolveComicListIncludeSpec 解析 includes[]，输出漫画列表所需的关联展开规格。
func ResolveComicListIncludeSpec(includes []string) ComicListIncludeSpec {
	spec := ComicListIncludeSpec{}

	for _, include := range includes {
		if value.CheckResourceNeccessary(include, value.IncludeWorkset) {
			spec.NeedWorkset = true
		}
	}

	return spec
}
