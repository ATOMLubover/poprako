package service

import "labelplus-next-web-be/internal/value"

type MemberListIncludeSpec struct {
	NeedUser bool
}

// ResolveMemberListIncludeSpec 解析 includes[]，输出成员列表所需的关联展开规格。
func ResolveMemberListIncludeSpec(includes []string) MemberListIncludeSpec {
	spec := MemberListIncludeSpec{}

	for _, include := range includes {
		if value.CheckResourceNeccessary(include, value.IncludeUser) {
			spec.NeedUser = true
		}
	}

	return spec
}

type MyMemberListIncludeSpec struct {
	NeedTeam bool
}

// ResolveMyMemberListIncludeSpec 解析 includes[]，输出"我的成员"列表所需的关联展开规格。
func ResolveMyMemberListIncludeSpec(includes []string) MyMemberListIncludeSpec {
	spec := MyMemberListIncludeSpec{}

	for _, include := range includes {
		if value.CheckResourceNeccessary(include, value.IncludeTeam) {
			spec.NeedTeam = true
		}
	}

	return spec
}
