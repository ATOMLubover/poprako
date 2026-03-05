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

	PermissionTeamListAll Permission = PrefixPermissionTeam + "list:all"
	PermissionTeamCreate  Permission = PrefixPermissionTeam + "create"
	PermissionTeamUpdate  Permission = PrefixPermissionTeam + "update"
	PermissionTeamDelete  Permission = PrefixPermissionTeam + "delete"
)

const (
	PrefixPermissionMember Permission = "permission:member:"

	PermissionMemberList   Permission = PrefixPermissionMember + "list"
	PermissionMemberCreate Permission = PrefixPermissionMember + "create"
	PermissionMemberUpdate Permission = PrefixPermissionMember + "update"
	PermissionMemberDelete Permission = PrefixPermissionMember + "delete"
)
