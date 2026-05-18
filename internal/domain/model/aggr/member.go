package aggr

import (
	"time"

	"poprako-s/internal/domain/model/enum"
)

// `Member` represents a user's membership in a team, including their assigned roles
type Member struct {
	// Embedded for roles ops.
	TimedRoles

	Id string

	UserId string
	// `UserNickname` is the redundant nickname copied from `t_user.nickname`.
	UserNickname string
	// `User` is only filled when `includes` option includes `user`.
	User *User

	TeamId string
	// `Team` is only filled when `includes` option includes `team`.
	Team *Team

	CreatedAt time.Time
	UpdatedAt time.Time
}

// `MemberCre` holds the input needed to create a new team membership record
type MemberCre struct {
	Id string

	UserId string
	// `UserNickname` is used to fill `t_member.user_nickname` during insert.
	UserNickname string
	TeamId       string

	RoleMask RoleMask
}

// `HasAnyRole` reports whether the creation input carries at least one of the given roles
func (c *MemberCre) HasAnyRole(r ...enum.Role) bool {
	return c.RoleMask.HasAnyRole(r...)
}

// `ToRoleMask` returns the role mask of the creation input
func (c *MemberCre) ToRoleMask() RoleMask {
	return c.RoleMask
}

// `ToRoleArr` expands the role mask of the creation input into a role slice
func (c *MemberCre) ToRoleArr() []enum.Role {
	return c.RoleMask.ToRoleArr()
}

// `FromRoleMask` sets the role mask of the creation input
func (c *MemberCre) FromRoleMask(m RoleMask) {
	c.RoleMask = m
}

// `MemberRoleUpd` holds the input needed to update the role set of an existing membership
type MemberRoleUpd struct {
	Id string

	RoleMask RoleMask
}

// `HasAnyRole` reports whether the update input carries at least one of the given roles
func (u *MemberRoleUpd) HasAnyRole(r ...enum.Role) bool {
	return u.RoleMask.HasAnyRole(r...)
}

// `ToRoleMask` returns the role mask of the update input
func (u *MemberRoleUpd) ToRoleMask() RoleMask {
	return u.RoleMask
}

// `ToRoleArr` expands the role mask of the update input into a role slice
func (u *MemberRoleUpd) ToRoleArr() []enum.Role {
	return u.RoleMask.ToRoleArr()
}

// `FromRoleMask` sets the role mask of the update input
func (u *MemberRoleUpd) FromRoleMask(m RoleMask) {
	u.RoleMask = m
}
