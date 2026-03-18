package service

import "labelplus-next-web-be/internal/value"

type ComicListIncludeSpec struct {
	NeedWorkset bool
	NeedCreator bool
}

// ResolveComicListIncludeSpec 解析 includes[]，输出漫画列表所需的关联展开规格。
func ResolveComicListIncludeSpec(includes []string) ComicListIncludeSpec {
	spec := ComicListIncludeSpec{}

	for _, include := range includes {
		if value.CheckResourceNeccessary(include, value.IncludeWorkset) {
			spec.NeedWorkset = true
		}
		if value.CheckResourceNeccessary(include, value.IncludeCreator) {
			spec.NeedCreator = true
		}
	}

	return spec
}
