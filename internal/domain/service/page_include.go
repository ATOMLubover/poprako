package service

import "labelplus-next-web-be/internal/value"

type PageListIncludeSpec struct {
	NeedCreator bool
}

// ResolvePageListIncludeSpec 解析页面列表接口的 includes[]。
func ResolvePageListIncludeSpec(includes []string) PageListIncludeSpec {
	spec := PageListIncludeSpec{}

	for _, include := range includes {
		if value.CheckResourceNeccessary(include, value.IncludeCreator) {
			spec.NeedCreator = true
		}
	}

	return spec
}
