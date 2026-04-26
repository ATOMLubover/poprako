package entity

import (
	"time"

	"poprako-s/internal/domain/model/aggr"
)

const USER_TABLE = "t_user"

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

func (*UserRow) TableName() string {
	return USER_TABLE
}

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

type UserCredsRow struct {
	Id string `gorm:"column:id;primaryKey"`

	PwdHash string `gorm:"column:password_hash"`
}

func (r *UserCredsRow) TableName() string {
	return "t_user_creds"
}

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

type UserRegRow struct {
	Id string `gorm:"column:id;primaryKey"`

	Nickname string `gorm:"column:nickname"`
	Qid      string `gorm:"column:qid;unique"`

	PwdHash string `gorm:"column:password_hash"`
}

func NewUserRegRowFromAggr(reg *aggr.UserReg) *UserRegRow {
	return &UserRegRow{
		Id:       reg.Id,
		Nickname: reg.Nickname,
		Qid:      reg.Qid,
		PwdHash:  reg.PwdHash,
	}
}

func (r *UserRegRow) TableName() string {
	return USER_TABLE
}
