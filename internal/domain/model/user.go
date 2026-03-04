package model

import (
	"time"
)

type UserInfo struct {
	ID string

	Name      string
	QQ        string
	AvatarURL string

	IsSuperAdmin bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserCredentials 仅用于登录时校验，不对外暴露完整用户信息
type UserCredentials struct {
	UserID       string
	PasswordHash string
}

func NewUserCredentials(userID string, passwordHash string) *UserCredentials {
	return &UserCredentials{
		UserID:       userID,
		PasswordHash: passwordHash,
	}
}

type UserRegistration struct {
	Name     string
	QQ       string
	Password string

	Roles RoleMask
}

func NewUserRegistration(
	name string,
	qq string,
	password string,
	roles ...RoleFlag,
) *UserRegistration {
	registration := &UserRegistration{
		Name:     name,
		QQ:       qq,
		Password: password,
		Roles:    MaskRoles(roles),
	}

	return registration
}

type UserPatch struct {
	Name         string
	QQ           string
	PasswordHash string
}

func NewUserPatch(
	name string,
	qq string,
	passwordHash string,
) *UserPatch {
	return &UserPatch{
		Name:         name,
		QQ:           qq,
		PasswordHash: passwordHash,
	}
}
