package aggr

import "poprako-s/internal/domain/model/enum"

// `WithRoles` is an interface for entities that have roles,
// such as `Member`, `MemberInv` or `Assignment`.
type WithRoles interface {
	HasAnyRole(r ...enum.Role) bool
	ToRoleMask() RoleMask
	ToRoleArr() []enum.Role
	FromRoleMask(m RoleMask)
	FromRoleArr(arr []enum.Role)
}
