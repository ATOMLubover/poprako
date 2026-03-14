package service

import "labelplus-next-web-be/internal/value"

type TeamListIncludeSpec struct {
	NeedMembers    bool
	NeedMemberUser bool
}

// ResolveTeamListIncludeSpec 解析 teams 列表接口的 includes[]。
// 支持：member, member.user
func ResolveTeamListIncludeSpec(includes []string) TeamListIncludeSpec {
	spec := TeamListIncludeSpec{}

	for _, include := range includes {
		if value.CheckResourceNeccessary(include, value.IncludeMember) {
			spec.NeedMembers = true
		}
		if value.CheckResourceNeccessary(include, value.IncludeMember, value.IncludeUser) {
			spec.NeedMembers = true
			spec.NeedMemberUser = true
		}
	}

	return spec
}
