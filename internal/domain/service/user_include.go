package service

import "labelplus-next-web-be/internal/value"

type UserDetailIncludeSpec struct {
	NeedMembers    bool
	NeedMemberTeam bool
}

// ResolveUserDetailIncludeSpec 解析 users 详情接口的 includes[]。
// 支持：member, member.team
func ResolveUserDetailIncludeSpec(includes []string) UserDetailIncludeSpec {
	spec := UserDetailIncludeSpec{}

	for _, include := range includes {
		if value.CheckResourceNeccessary(include, value.IncludeMember) {
			spec.NeedMembers = true
		}
		if value.CheckResourceNeccessary(include, value.IncludeMember, value.IncludeTeam) {
			spec.NeedMembers = true
			spec.NeedMemberTeam = true
		}
	}

	return spec
}
