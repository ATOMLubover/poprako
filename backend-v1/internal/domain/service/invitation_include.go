package service

import "labelplus-next-web-be/internal/value"

type InvitationListIncludeSpec struct {
	NeedInvitor bool
}

// ResolveInvitationListIncludeSpec 解析邀请列表接口的 includes[]。
func ResolveInvitationListIncludeSpec(includes []string) InvitationListIncludeSpec {
	spec := InvitationListIncludeSpec{}

	for _, include := range includes {
		if value.CheckResourceNeccessary(include, value.IncludeInvitor) {
			spec.NeedInvitor = true
		}
	}

	return spec
}
