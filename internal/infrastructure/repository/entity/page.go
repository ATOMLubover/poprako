package entity

import (
	"time"

	"labelplus-next-web-be/internal/domain/model"
)

const (
	ChapterTable = "chapter_table"
	PageTable    = "page_table"
)

type ChapterInfoRow struct {
	ID string `gorm:"column:id"`

	ComicID   string `gorm:"column:comic_id"`
	Index     int    `gorm:"column:index"`
	ChapterNo string `gorm:"column:subtitle"`

	PageCount int `gorm:"column:page_count"`

	CoverURL string `gorm:"column:cover_url"`

	UploadedAt     *time.Time `gorm:"column:uploaded_at"`
	TransalatingAt *time.Time `gorm:"column:transalating_at"`
	TranslatedAt   *time.Time `gorm:"column:translated_at"`
	ProofreadingAt *time.Time `gorm:"column:proofreading_at"`
	ProofreadAt    *time.Time `gorm:"column:proofread_at"`
	TypesettingAt  *time.Time `gorm:"column:typesetting_at"`
	TypesetAt      *time.Time `gorm:"column:typeset_at"`
	ReviewedAt     *time.Time `gorm:"column:reviewed_at"`
	PublishedAt    *time.Time `gorm:"column:published_at"`

	CreatorID string `gorm:"column:creator_id"`

	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

func (ChapterInfoRow) TableName() string { return ChapterTable }

type ChapterInsertRow struct {
	ID        string `gorm:"column:id"`
	ComicID   string `gorm:"column:comic_id"`
	Index     int    `gorm:"column:index"`
	ChapterNo string `gorm:"column:subtitle"`
	CreatorID string `gorm:"column:creator_id"`
}

func (ChapterInsertRow) TableName() string { return ChapterTable }

type PageInfoRow struct {
	ID string `gorm:"column:id"`

	ChapterID string `gorm:"column:chapter_id"`
	Index     int    `gorm:"column:index"`
	OSSKey    string `gorm:"column:oss_key"`
	CreatorID string `gorm:"column:creator_id"`

	IsUploaded bool `gorm:"column:uploaded"`

	TotalUnitCount      int `gorm:"column:total_unit_count"`
	TranslatedUnitCount int `gorm:"column:translated_unit_count"`
	ProofreadUnitCount  int `gorm:"column:proofread_unit_count"`

	// Creator 别名列（IncludeCreatorInfo() 时填充）
	CreatorName             string    `gorm:"column:creator_name"`
	CreatorQQ               string    `gorm:"column:creator_qq"`
	CreatorAvatarOSSKey     string    `gorm:"column:creator_avatar_oss_key"`
	CreatorIsAvatarUploaded bool      `gorm:"column:creator_is_avatar_uploaded"`
	CreatorIsSuperAdmin     bool      `gorm:"column:creator_is_super_admin"`
	CreatorCreatedAt        time.Time `gorm:"column:creator_created_at"`
	CreatorUpdatedAt        time.Time `gorm:"column:creator_updated_at"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (PageInfoRow) TableName() string { return PageTable }

type PageInsertRow struct {
	ID string `gorm:"column:id"`

	ChapterID string `gorm:"column:chapter_id"`
	Index     int    `gorm:"column:index"`
	OSSKey    string `gorm:"column:oss_key"`
	CreatorID string `gorm:"column:creator_id"`
}

func (PageInsertRow) TableName() string { return PageTable }

func ToPageInfo(row PageInfoRow) model.PageInfo {
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

	return model.NewPageInfo(
		row.ID,
		row.ChapterID,
		row.Index,
		row.OSSKey,
		row.IsUploaded,
		row.CreatorID,
		creatorInfo,
		row.TotalUnitCount,
		row.TranslatedUnitCount,
		row.ProofreadUnitCount,
		row.CreatedAt,
		row.UpdatedAt,
	)
}

func ToChapterInfo(row ChapterInfoRow) model.ChapterInfo {
	return model.NewChapterDetail(
		row.ID,
		row.ComicID,
		row.Index,
		row.ChapterNo,
		row.PageCount,
		0,
		0,
		0,
		row.CoverURL,
		row.UploadedAt,
		row.TransalatingAt,
		row.TranslatedAt,
		row.ProofreadingAt,
		row.ProofreadAt,
		row.TypesettingAt,
		row.TypesetAt,
		row.ReviewedAt,
		row.PublishedAt,
		row.CreatorID,
		row.CreatedAt,
		row.UpdatedAt,
	)
}

// ChapterWithInfoRow 是统一的聚合行类型，可携带 creator 别名列。
// Creator 列由 IncludeCreatorInfo() query option 决定是否 JOIN。
type ChapterWithInfoRow struct {
	ChapterInfoRow

	// Creator 别名列（IncludeCreatorInfo() 时填充）
	CreatorName             string    `gorm:"column:creator_name"`
	CreatorQQ               string    `gorm:"column:creator_qq"`
	CreatorAvatarOSSKey     string    `gorm:"column:creator_avatar_oss_key"`
	CreatorIsAvatarUploaded bool      `gorm:"column:creator_is_avatar_uploaded"`
	CreatorIsSuperAdmin     bool      `gorm:"column:creator_is_super_admin"`
	CreatorCreatedAt        time.Time `gorm:"column:creator_created_at"`
	CreatorUpdatedAt        time.Time `gorm:"column:creator_updated_at"`
}

func ToChapterWithInfo(row ChapterWithInfoRow) model.ChapterInfo {
	chapter := ToChapterInfo(row.ChapterInfoRow)

	if !row.CreatorCreatedAt.IsZero() {
		creatorInfo := model.NewUserInfo(
			row.CreatorID,
			row.CreatorName,
			row.CreatorQQ,
			row.CreatorAvatarOSSKey,
			row.CreatorIsAvatarUploaded,
			row.CreatorIsSuperAdmin,
			row.CreatorCreatedAt,
			row.CreatorUpdatedAt,
		)
		chapter.Creator = &creatorInfo
	}

	return chapter
}
