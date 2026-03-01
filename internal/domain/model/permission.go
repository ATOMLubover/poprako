package model

type Permission string

const (
	PrefixPermissionInvitation Permission = "invitations:"

	PermissionInvitationsList   Permission = PrefixPermissionInvitation + "list"
	PermissionInvitationsCreate Permission = PrefixPermissionInvitation + "create"
	PermissionInvitationsDelete Permission = PrefixPermissionInvitation + "delete"
	PermissionInvitationsPatch  Permission = PrefixPermissionInvitation + "patch"
)

const (
	PrefixPermissionUsers Permission = "users:"

	PermissionUsersRemove Permission = PrefixPermissionUsers + "remove"
)
