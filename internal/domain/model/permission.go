package model

type Permission string

const (
	PrefixPermissionInvitation Permission = "permission:invitation:"

	PermissionInvitationList   Permission = PrefixPermissionInvitation + "list"
	PermissionInvitationCreate Permission = PrefixPermissionInvitation + "create"
	PermissionInvitationDelete Permission = PrefixPermissionInvitation + "delete"
	PermissionInvitationUpdate Permission = PrefixPermissionInvitation + "update"
)

const (
	PrefixPermissionUser Permission = "permission:user:"

	PermissionUserRemove Permission = PrefixPermissionUser + "remove"
)

const (
	PrefixPermissionTeam Permission = "permission:team:"

	PermissionTeamListMine Permission = PrefixPermissionTeam + "list:mine"
	PermissionTeamListAll  Permission = PrefixPermissionTeam + "list:all"
	PermissionTeamCreate   Permission = PrefixPermissionTeam + "create"
	PermissionTeamUpdate   Permission = PrefixPermissionTeam + "update"
	PermissionTeamDelete   Permission = PrefixPermissionTeam + "delete"
)

const (
	PrefixPermissionMember Permission = "permission:member:"

	PermissionMemberList   Permission = PrefixPermissionMember + "list"
	PermissionMemberUpdate Permission = PrefixPermissionMember + "update"
	PermissionMemberDelete Permission = PrefixPermissionMember + "delete"
)
