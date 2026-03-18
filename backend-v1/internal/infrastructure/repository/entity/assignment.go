package entity

import (
	"time"

	"labelplus-next-web-be/internal/domain/model"
)

const AssignmentTable = "assignment_table"

type AssignmentInfoRow struct {
	ID string `gorm:"column:id"`

	ChapterID string `gorm:"column:chapter_id"`
	UserID    string `gorm:"column:user_id"`

	AssignedRawProviderAt *time.Time `gorm:"column:assigned_raw_provider_at"`
	AssignedTranslatorAt  *time.Time `gorm:"column:assigned_translator_at"`
	AssignedProofreaderAt *time.Time `gorm:"column:assigned_proofreader_at"`
	AssignedTypesetterAt  *time.Time `gorm:"column:assigned_typesetter_at"`
	AssignedRedrawerAt    *time.Time `gorm:"column:assigned_redrawer_at"`
	AssignedReviewerAt    *time.Time `gorm:"column:assigned_reviewer_at"`
	AssignedPublisherAt   *time.Time `gorm:"column:assigned_publisher_at"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (AssignmentInfoRow) TableName() string { return AssignmentTable }

func ToAssignmentInfo(row AssignmentInfoRow) model.AssignmentInfo {
	return model.NewAssignmentInfo(
		row.ID,
		row.ChapterID,
		row.UserID,
		row.AssignedRawProviderAt,
		row.AssignedTranslatorAt,
		row.AssignedProofreaderAt,
		row.AssignedTypesetterAt,
		row.AssignedRedrawerAt,
		row.AssignedReviewerAt,
		row.AssignedPublisherAt,
		row.CreatedAt,
		row.UpdatedAt,
	)
}

// AssignmentInsertRow 用于 Create，仅包含写入所需字段。
type AssignmentInsertRow struct {
	ID        string `gorm:"column:id"`
	ChapterID string `gorm:"column:chapter_id"`
	UserID    string `gorm:"column:user_id"`

	AssignedRawProviderAt *time.Time `gorm:"column:assigned_raw_provider_at"`
	AssignedTranslatorAt  *time.Time `gorm:"column:assigned_translator_at"`
	AssignedProofreaderAt *time.Time `gorm:"column:assigned_proofreader_at"`
	AssignedTypesetterAt  *time.Time `gorm:"column:assigned_typesetter_at"`
	AssignedRedrawerAt    *time.Time `gorm:"column:assigned_redrawer_at"`
	AssignedReviewerAt    *time.Time `gorm:"column:assigned_reviewer_at"`
	AssignedPublisherAt   *time.Time `gorm:"column:assigned_publisher_at"`
}

func (AssignmentInsertRow) TableName() string { return AssignmentTable }

// AssignmentIncludeRow 是统一的聚合行类型，可同时携带 user 和 chapter+comic 别名列。
// User 列与 Chapter/Comic 列均可选，由 query option 决定是否 JOIN 及 SELECT。
type AssignmentIncludeRow struct {
	AssignmentInfoRow

	// User 别名列（IncludeUserInfo() 时填充）
	UserName             string    `gorm:"column:user_name"`
	UserQQ               string    `gorm:"column:user_qq"`
	UserAvatarOSSKey     string    `gorm:"column:user_avatar_oss_key"`
	UserIsAvatarUploaded bool      `gorm:"column:user_is_avatar_uploaded"`
	UserIsSuperAdmin     bool      `gorm:"column:user_is_super_admin"`
	UserCreatedAt        time.Time `gorm:"column:user_created_at"`
	UserUpdatedAt        time.Time `gorm:"column:user_updated_at"`

	// Chapter 别名列（IncludeChapterInfo() 时填充）
	ChapterComicID   string    `gorm:"column:chapter_comic_id"`
	ChapterIndex     int       `gorm:"column:chapter_index"`
	ChapterSubtitle  string    `gorm:"column:chapter_subtitle"`
	ChapterPageCount int       `gorm:"column:chapter_page_count"`
	ChapterCreatorID string    `gorm:"column:chapter_creator_id"`
	ChapterCreatedAt time.Time `gorm:"column:chapter_created_at"`
	ChapterUpdatedAt time.Time `gorm:"column:chapter_updated_at"`

	// Comic 别名列（IncludeChapterInfo() 时填充，隐含 chapter.comic）
	ComicWorksetID    string    `gorm:"column:comic_workset_id"`
	ComicIndex        int       `gorm:"column:comic_index"`
	ComicTitle        string    `gorm:"column:comic_title"`
	ComicAuthor       string    `gorm:"column:comic_author"`
	ComicDescription  string    `gorm:"column:comic_description"`
	ComicChapterCount int       `gorm:"column:comic_chapter_count"`
	ComicCreatorID    string    `gorm:"column:comic_creator_id"`
	ComicLastActiveAt time.Time `gorm:"column:comic_last_active_at"`
	ComicCreatedAt    time.Time `gorm:"column:comic_created_at"`
	ComicUpdatedAt    time.Time `gorm:"column:comic_updated_at"`

	// ChapterCreator 别名列（IncludeRelationInfo() 含 chapter.creator 时填充）
	ChapterCreatorName             string    `gorm:"column:chapter_creator_name"`
	ChapterCreatorQQ               string    `gorm:"column:chapter_creator_qq"`
	ChapterCreatorAvatarOSSKey     string    `gorm:"column:chapter_creator_avatar_oss_key"`
	ChapterCreatorIsAvatarUploaded bool      `gorm:"column:chapter_creator_is_avatar_uploaded"`
	ChapterCreatorIsSuperAdmin     bool      `gorm:"column:chapter_creator_is_super_admin"`
	ChapterCreatorCreatedAt        time.Time `gorm:"column:chapter_creator_created_at"`
	ChapterCreatorUpdatedAt        time.Time `gorm:"column:chapter_creator_updated_at"`
}

// ToAssignmentInfoFromIncludeRow 将聚合行转换为 model.AssignmentInfo。
// 通过检查时间零值判断对应 include 是否实际被查询。
func ToAssignmentInfoFromIncludeRow(row AssignmentIncludeRow) model.AssignmentInfo {
	assignmentInfo := ToAssignmentInfo(row.AssignmentInfoRow)

	if !row.UserCreatedAt.IsZero() {
		user := model.NewUserInfo(
			row.UserID,
			row.UserName,
			row.UserQQ,
			row.UserAvatarOSSKey,
			row.UserIsAvatarUploaded,
			row.UserIsSuperAdmin,
			row.UserCreatedAt,
			row.UserUpdatedAt,
		)
		assignmentInfo.User = &user
	}

	if !row.ChapterCreatedAt.IsZero() {
		chapter := model.NewChapterDetail(
			row.ChapterID,
			row.ChapterComicID,
			row.ChapterIndex,
			row.ChapterSubtitle,
			row.ChapterPageCount,
			0,
			0,
			0,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			row.ChapterCreatorID,
			row.ChapterCreatedAt,
			row.ChapterUpdatedAt,
		)

		if !row.ComicCreatedAt.IsZero() {
			comic := model.NewComicInfo(
				row.ChapterComicID,
				row.ComicWorksetID,
				nil,
				row.ComicIndex,
				row.ComicTitle,
				row.ComicAuthor,
				row.ComicDescription,
				row.ComicChapterCount,
				row.ComicCreatorID,
				nil,
				row.ComicLastActiveAt,
				row.ComicCreatedAt,
				row.ComicUpdatedAt,
			)
			chapter.Comic = &comic
		}

		if !row.ChapterCreatorCreatedAt.IsZero() {
			creator := model.NewUserInfo(
				row.ChapterCreatorID,
				row.ChapterCreatorName,
				row.ChapterCreatorQQ,
				row.ChapterCreatorAvatarOSSKey,
				row.ChapterCreatorIsAvatarUploaded,
				row.ChapterCreatorIsSuperAdmin,
				row.ChapterCreatorCreatedAt,
				row.ChapterCreatorUpdatedAt,
			)
			chapter.Creator = &creator
		}

		assignmentInfo.Chapter = &chapter
	}

	return assignmentInfo
}
