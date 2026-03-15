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

// AssignmentWithUserRow 用于 ListWithUserInfo（JOIN user_table）。
type AssignmentWithUserRow struct {
	AssignmentInfoRow

	UserName             string    `gorm:"column:user_name"`
	UserQQ               string    `gorm:"column:user_qq"`
	UserAvatarOSSKey     string    `gorm:"column:user_avatar_oss_key"`
	UserIsAvatarUploaded bool      `gorm:"column:user_is_avatar_uploaded"`
	UserIsSuperAdmin     bool      `gorm:"column:user_is_super_admin"`
	UserCreatedAt        time.Time `gorm:"column:user_created_at"`
	UserUpdatedAt        time.Time `gorm:"column:user_updated_at"`
}

func ToAssignmentInfoWithUser(row AssignmentWithUserRow) model.AssignmentInfo {
	user := model.UserInfo{
		ID:               row.UserID,
		Name:             row.UserName,
		QQ:               row.UserQQ,
		AvatarOSSKey:     row.UserAvatarOSSKey,
		IsAvatarUploaded: row.UserIsAvatarUploaded,
		IsSuperAdmin:     row.UserIsSuperAdmin,
		CreatedAt:        row.UserCreatedAt,
		UpdatedAt:        row.UserUpdatedAt,
	}

	assignmentInfo := ToAssignmentInfo(row.AssignmentInfoRow)
	assignmentInfo.User = &user

	return assignmentInfo
}

// AssignmentWithChapterAndComicRow 用于 ListWithChapterInfo（JOIN chapter_table + workset_table + comic_table）。
type AssignmentWithChapterAndComicRow struct {
	AssignmentInfoRow

	ChapterComicID   string    `gorm:"column:chapter_comic_id"`
	ChapterIndex     int       `gorm:"column:chapter_index"`
	ChapterSubtitle  string    `gorm:"column:chapter_subtitle"`
	ChapterPageCount int       `gorm:"column:chapter_page_count"`
	ChapterCoverURL  string    `gorm:"column:chapter_cover_url"`
	ChapterCreatedAt time.Time `gorm:"column:chapter_created_at"`
	ChapterUpdatedAt time.Time `gorm:"column:chapter_updated_at"`

	ComicWorksetID    string    `gorm:"column:comic_workset_id"`
	ComicTeamID       string    `gorm:"column:comic_team_id"`
	ComicIndex        int       `gorm:"column:comic_index"`
	ComicTitle        string    `gorm:"column:comic_title"`
	ComicAuthor       string    `gorm:"column:comic_author"`
	ComicDescription  string    `gorm:"column:comic_description"`
	ComicCoverURL     string    `gorm:"column:comic_cover_url"`
	ComicChapterCount int       `gorm:"column:comic_chapter_count"`
	ComicCreatorID    string    `gorm:"column:comic_creator_id"`
	ComicLastActiveAt time.Time `gorm:"column:comic_last_active_at"`
	ComicCreatedAt    time.Time `gorm:"column:comic_created_at"`
	ComicUpdatedAt    time.Time `gorm:"column:comic_updated_at"`
}

func ToAssignmentInfoWithChapter(row AssignmentWithChapterAndComicRow) model.AssignmentInfo {
	comic := model.NewComicInfo(
		row.ChapterComicID,
		row.ComicWorksetID,
		row.ComicTeamID,
		nil,
		row.ComicIndex,
		row.ComicTitle,
		row.ComicAuthor,
		row.ComicDescription,
		row.ComicCoverURL,
		row.ComicChapterCount,
		row.ComicCreatorID,
		nil,
		row.ComicLastActiveAt,
		row.ComicCreatedAt,
		row.ComicUpdatedAt,
	)

	chapter := model.NewChapterDetail(
		row.ChapterID,
		row.ChapterComicID,
		row.ChapterIndex,
		row.ChapterSubtitle,
		row.ChapterPageCount,
		0,
		0,
		0,
		row.ChapterCoverURL,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		"",
		row.ChapterCreatedAt,
		row.ChapterUpdatedAt,
	)
	chapter.Comic = &comic

	assignmentInfo := ToAssignmentInfo(row.AssignmentInfoRow)
	assignmentInfo.Chapter = &chapter

	return assignmentInfo
}
