package repo_entity

import (
	"time"

	"poprako-s/internal/domain/model"
)

const (
	UserTable      = "user_table"
	UserStatsTable = "user_stats_table"
)

type UserInfoRow struct {
	ID string `gorm:"column:id"`

	Name string `gorm:"column:name"`
	QQ   string `gorm:"column:qq"`

	AvatarOSSKey     string `gorm:"column:avatar_oss_key"`
	IsAvatarUploaded bool   `gorm:"column:is_avatar_uploaded"`

	PasswordHash string `gorm:"column:password_hash"`

	IsSuperAdmin bool `gorm:"column:is_super_admin"`

	LastLoginAt *time.Time `gorm:"column:last_login_at"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at"`
}

type UserCredsRow struct {
	QQ           string `gorm:"column:qq"`
	PasswordHash string `gorm:"column:password_hash"`
}

type UserStatsRow struct {
	ID string `gorm:"column:id"`

	UserID string `gorm:"column:user_id"`

	TotalAssignmentCount    int `gorm:"column:total_assignment_count"`
	ActiveAssignmentCount   int `gorm:"column:active_assignment_count"`
	FinishedAssignmentCount int `gorm:"column:finished_assignment_count"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func ToUserInfo(row UserInfoRow) model.UserInfo {
	info := model.UserInfo{
		ID:               row.ID,
		Name:             row.Name,
		QQ:               row.QQ,
		AvatarKey:        row.AvatarOSSKey,
		IsAvatarUploaded: row.IsAvatarUploaded,
		IsSuperAdmin:     row.IsSuperAdmin,
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        row.UpdatedAt,
	}

	if row.LastLoginAt != nil {
		info.LastLoginAt = *row.LastLoginAt
	}

	return info
}

func ToUserCreds(row UserCredsRow) model.UserCreds {
	return model.UserCreds{
		QQ:      row.QQ,
		PwdHash: row.PasswordHash,
	}
}

func ToUserStats(row UserStatsRow) model.UserStats {
	return model.UserStats{
		UserID:                  row.UserID,
		TotalAssignmentCount:    row.TotalAssignmentCount,
		ActiveAssignmentCount:   row.ActiveAssignmentCount,
		FinishedAssignmentCount: row.FinishedAssignmentCount,
	}
}
