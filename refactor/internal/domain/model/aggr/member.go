package aggr

import (
	"time"

	"poprako-s/internal/domain/model/enum"
)

type Member struct {
	// Embedded for roles ops.
	TimedRoles

	Id string

	UserId string
	// `User` is only filled when `includes` option includes `user`.
	User *User

	TeamId string
	// `Team` is only filled when `includes` option includes `team`.
	Team *Team

	CreatedAt time.Time
	UpdatedAt time.Time
}

type MemberCre struct {
	Id string

	UserId string
	TeamId string

	RoleMask RoleMask
}

func (c *MemberCre) HasAnyRole(r ...enum.Role) bool {
	return c.RoleMask.HasAnyRole(r...)
}

func (c *MemberCre) ToRoleMask() RoleMask {
	return c.RoleMask
}

func (c *MemberCre) ToRoleArr() []enum.Role {
	return c.RoleMask.ToRoleArr()
}

func (c *MemberCre) FromRoleMask(m RoleMask) {
	c.RoleMask = m
}

type MemberRoleUpd struct {
	Id string

	RoleMask RoleMask
}

func (u *MemberRoleUpd) HasAnyRole(r ...enum.Role) bool {
	return u.RoleMask.HasAnyRole(r...)
}

func (u *MemberRoleUpd) ToRoleMask() RoleMask {
	return u.RoleMask
}

func (u *MemberRoleUpd) ToRoleArr() []enum.Role {
	return u.RoleMask.ToRoleArr()
}

func (u *MemberRoleUpd) FromRoleMask(m RoleMask) {
	u.RoleMask = m
}
