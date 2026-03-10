package entity

import (
	"time"

	"labelplus-next-web-be/internal/domain/model"
)

const UserTable = "user_table"

// UserInfoRow 用于 List、GetInfoByID，不含密码字段
type UserInfoRow struct {
	ID               string    `gorm:"column:id"`
	Name             string    `gorm:"column:name"`
	QQ               string    `gorm:"column:qq"`
	AvatarOSSKey     string    `gorm:"column:avatar_oss_key"`
	IsAvatarUploaded bool      `gorm:"column:is_avatar_uploaded"`
	IsSuperAdmin     bool      `gorm:"column:is_super_admin"`
	CreatedAt        time.Time `gorm:"column:created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at"`
}

func (UserInfoRow) TableName() string { return UserTable }

func ToUserInfo(row UserInfoRow) model.UserInfo {
	return model.UserInfo{
		ID:               row.ID,
		Name:             row.Name,
		QQ:               row.QQ,
		AvatarOSSKey:     row.AvatarOSSKey,
		IsAvatarUploaded: row.IsAvatarUploaded,
		IsSuperAdmin:     row.IsSuperAdmin,
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        row.UpdatedAt,
	}
}

// UserCredentialsRow 用于 GetCredentialsByQQ，仅含鉴权所需的两个字段
type UserCredentialsRow struct {
	ID           string `gorm:"column:id"`
	PasswordHash string `gorm:"column:password_hash"`
}

func (UserCredentialsRow) TableName() string { return UserTable }

// UserInsertRow 用于 Create，含所有需要写入的字段
type UserInsertRow struct {
	ID               string `gorm:"column:id"`
	Name             string `gorm:"column:name"`
	QQ               string `gorm:"column:qq"`
	AvatarOSSKey     string `gorm:"column:avatar_oss_key"`
	IsAvatarUploaded bool   `gorm:"column:is_avatar_uploaded"`
	PasswordHash     string `gorm:"column:password_hash"`
	IsSuperAdmin     bool   `gorm:"column:is_super_admin"`
}

func (UserInsertRow) TableName() string { return UserTable }
