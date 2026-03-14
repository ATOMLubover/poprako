package entity

import (
	"time"

	"labelplus-next-web-be/internal/domain/model"
)

const ComicTable = "comic_table"

// ComicInfoRow 用于 List、Get，映射完整漫画信息。
// team_id 通过 workset_table JOIN 别名填充，workset 别名列仅在 IncludeWorksetInfo() 时填充。
type ComicInfoRow struct {
	ID        string `gorm:"column:id"`
	WorksetID string `gorm:"column:workset_id"`

	// team_id 通过 JOIN workset_table 获取，并以 team_id 别名返回（SELECT workset_table.team_id AS team_id）
	TeamID string `gorm:"column:team_id"`

	Index       int    `gorm:"column:index"`
	Title       string `gorm:"column:title"`
	Author      string `gorm:"column:author"`
	Description string `gorm:"column:description"`

	CoverURL string `gorm:"column:cover_url"`

	ChapterCount int    `gorm:"column:chapter_count"`
	CreatorID    string `gorm:"column:creator_id"`

	LastActiveAt time.Time `gorm:"column:last_active_at"`

	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`

	// Workset 别名列（IncludeWorksetInfo() 时填充）
	WorksetName       string    `gorm:"column:workset_name"`
	WorksetIndex      int       `gorm:"column:workset_index"`
	WorksetComicCount int       `gorm:"column:workset_comic_count"`
	WorksetCreatedAt  time.Time `gorm:"column:workset_created_at"`
	WorksetUpdatedAt  time.Time `gorm:"column:workset_updated_at"`

	// Creator 别名列（IncludeCreatorInfo() 时填充）
	CreatorName             string    `gorm:"column:creator_name"`
	CreatorQQ               string    `gorm:"column:creator_qq"`
	CreatorAvatarOSSKey     string    `gorm:"column:creator_avatar_oss_key"`
	CreatorIsAvatarUploaded bool      `gorm:"column:creator_is_avatar_uploaded"`
	CreatorIsSuperAdmin     bool      `gorm:"column:creator_is_super_admin"`
	CreatorCreatedAt        time.Time `gorm:"column:creator_created_at"`
	CreatorUpdatedAt        time.Time `gorm:"column:creator_updated_at"`
}

func (ComicInfoRow) TableName() string { return ComicTable }

// ComicInsertRow 用于 Create，仅包含写入所需字段。
type ComicInsertRow struct {
	ID          string `gorm:"column:id"`
	WorksetID   string `gorm:"column:workset_id"`
	Index       int    `gorm:"column:index"`
	Title       string `gorm:"column:title"`
	Author      string `gorm:"column:author"`
	Description string `gorm:"column:description"`
	CoverURL    string `gorm:"column:cover_url"`
	CreatorID   string `gorm:"column:creator_id"`
}

func (ComicInsertRow) TableName() string { return ComicTable }

func ToComicInfo(row ComicInfoRow) model.ComicInfo {
	var worksetInfo *model.WorksetInfo
	if !row.WorksetCreatedAt.IsZero() {
		info := model.NewWorksetInfo(
			row.WorksetID,
			row.TeamID,
			nil,
			row.WorksetIndex,
			row.WorksetName,
			"",
			row.WorksetComicCount,
			row.WorksetCreatedAt,
			row.WorksetUpdatedAt,
		)
		worksetInfo = &info
	}

	var creatorInfo *model.UserInfo
	if !row.CreatorCreatedAt.IsZero() {
		info := model.NewUserInfo(
			row.CreatorID,
			row.CreatorName,
			row.CreatorQQ,
			row.CreatorAvatarOSSKey,
			row.CreatorIsAvatarUploaded,
			row.CreatorIsSuperAdmin,
			row.CreatorCreatedAt,
			row.CreatorUpdatedAt,
		)
		creatorInfo = &info
	}

	return model.NewComicInfo(
		row.ID,
		row.WorksetID,
		row.TeamID,
		worksetInfo,
		row.Index,
		row.Title,
		row.Author,
		row.Description,
		row.CoverURL,
		row.ChapterCount,
		row.CreatorID,
		creatorInfo,
		row.LastActiveAt,
		row.CreatedAt,
		row.UpdatedAt,
	)
}
