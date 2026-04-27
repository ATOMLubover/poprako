package entity

import (
	"time"

	"poprako-s/internal/domain/model/aggr"
)

// `USER_TABLE` is the table name for the user entity
const USER_TABLE = "t_user"

// `UserRow` maps a full user record for read queries
type UserRow struct {
	Id string `gorm:"column:id;primaryKey"`

	Nickname string `gorm:"column:nickname"`
	Qid      string `gorm:"column:qid;unique"`

	AvatarKey      string `gorm:"column:avatar_key"`
	AvatarUploaded bool   `gorm:"column:avatar_uploaded"`

	IsSuperAdmin bool `gorm:"column:is_super_admin"`

	LastActiveAt time.Time `gorm:"column:last_active_at"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// `TableName` returns the table name of `UserRow`
func (*UserRow) TableName() string {
	return USER_TABLE
}

// `ToUserAggr` converts a `UserRow` into a `User` aggregate
func (r *UserRow) ToUserAggr() *aggr.User {
	if r == nil {
		// Do not treat nil as an error, as it may be used in includes.
		return nil
	}

	return &aggr.User{
		Id: r.Id,

		Nickname: r.Nickname,
		Qid:      r.Qid,

		AvatarKey:      r.AvatarKey,
		AvatarUploaded: r.AvatarUploaded,

		IsSuperAdmin: r.IsSuperAdmin,

		LastActiveAt: r.LastActiveAt,

		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

// `UserCredsRow` maps only the credential columns needed for password verification
type UserCredsRow struct {
	Id string `gorm:"column:id;primaryKey"`

	PwdHash string `gorm:"column:password_hash"`
}

// `TableName` returns the table name of `UserCredsRow`
func (r *UserCredsRow) TableName() string {
	return USER_TABLE
}

// `ToUserCredsAggr` converts a `UserCredsRow` into a `UserCreds` aggregate
func (r *UserCredsRow) ToUserCredsAggr() *aggr.UserCreds {
	if r == nil {
		// Do not treat nil as an error, as it may be used in includes.
		return nil
	}

	return &aggr.UserCreds{
		Id:      r.Id,
		PwdHash: r.PwdHash,
	}
}

// `UserRegRow` maps only the columns required for a new user insert
type UserRegRow struct {
	Id string `gorm:"column:id;primaryKey"`

	Nickname string `gorm:"column:nickname"`
	Qid      string `gorm:"column:qid;unique"`

	PwdHash string `gorm:"column:password_hash"`
}

// `NewUserRegRowFromAggr` builds a `UserRegRow` from a `UserReg` aggregate
func NewUserRegRowFromAggr(reg *aggr.UserReg) *UserRegRow {
	return &UserRegRow{
		Id:       reg.Id,
		Nickname: reg.Nickname,
		Qid:      reg.Qid,
		PwdHash:  reg.PwdHash,
	}
}

// `TableName` returns the table name of `UserRegRow`
func (r *UserRegRow) TableName() string {
	return USER_TABLE
}
