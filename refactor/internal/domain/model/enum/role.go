package enum

type Role uint32

const (
	RoleRawProvider Role = 1 << iota
	RoleTranslator
	RoleProofreader
	RoleTypesetter
	RoleRedrawer
	RoleReviewer
	RolePublisher

	// `RoleAdmin` is a special role that does not
	// exist in chapter assignment. It only
	// works in team level.
	RoleAdmin

	// Range henlpers.
	RoleInf Role = RoleRawProvider
	RoleSup Role = RoleAdmin
)
