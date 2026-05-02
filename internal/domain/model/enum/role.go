package enum

// `Role` is a bitmask flag that identifies a single team or assignment role
type Role uint32

// Role flag constants, each occupying one bit position in a `RoleMask`
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

	// Range helpers.
	RoleInf Role = RoleRawProvider
	RoleSup Role = RoleAdmin
)
