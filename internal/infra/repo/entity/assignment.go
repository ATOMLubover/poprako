package entity

import (
	"time"

	"poprako-s/internal/domain/model/aggr"
)

// `ASSIGNMENT_TABLE` is table name for assignment.
const ASSIGNMENT_TABLE = "t_assignment"

// `AssignmentRow` maps one assignment row for read queries.
type AssignmentRow struct {
	Id string `gorm:"column:id;primaryKey"`

	ChapterId string `gorm:"column:chapter_id"`
	UserId    string `gorm:"column:user_id"`
	// `User` is included when `includes` contains `user`.
	User *UserRow `gorm:"foreignKey:UserId"`
	// `Chapter` is included when includes contain chapter paths.
	Chapter *ChapterRow `gorm:"foreignKey:ChapterId"`

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

// `TableName` returns table name for `AssignmentRow`.
func (*AssignmentRow) TableName() string {
	return ASSIGNMENT_TABLE
}

// `ToAssignmentAggr` converts row to assignment aggregate.
func (r *AssignmentRow) ToAssignmentAggr() *aggr.Assignment {
	if r == nil {
		return nil
	}

	var user *aggr.User
	if r.User != nil {
		user = r.User.ToUserAggr()
	}

	var chapter *aggr.Chapter
	if r.Chapter != nil {
		chapter = r.Chapter.ToChapterAggr()
	}

	return &aggr.Assignment{
		Id:        r.Id,
		ChapterId: r.ChapterId,
		UserId:    r.UserId,
		User:      user,
		Chapter:   chapter,
		TimedRoles: aggr.TimedRoles{
			AssignedRawProviderAt: r.AssignedRawProviderAt,
			AssignedTranslatorAt:  r.AssignedTranslatorAt,
			AssignedProofreaderAt: r.AssignedProofreaderAt,
			AssignedTypesetterAt:  r.AssignedTypesetterAt,
			AssignedRedrawerAt:    r.AssignedRedrawerAt,
			AssignedReviewerAt:    r.AssignedReviewerAt,
			AssignedPublisherAt:   r.AssignedPublisherAt,
		},
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

// `AssignmentCreRow` maps assignment columns for create queries.
type AssignmentCreRow struct {
	Id string `gorm:"column:id;primaryKey"`

	ChapterId string `gorm:"column:chapter_id"`
	UserId    string `gorm:"column:user_id"`

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

// `TableName` returns table name for `AssignmentCreRow`.
func (*AssignmentCreRow) TableName() string {
	return ASSIGNMENT_TABLE
}

// `NewAssignmentCreRowFromAggr` builds create row from assignment create payload.
func NewAssignmentCreRowFromAggr(cre *aggr.AssignmentCre) *AssignmentCreRow {
	now := time.Now()

	return &AssignmentCreRow{
		Id:                    cre.Id,
		ChapterId:             cre.ChapterId,
		UserId:                cre.UserId,
		AssignedRawProviderAt: cre.TimedRoles.AssignedRawProviderAt,
		AssignedTranslatorAt:  cre.TimedRoles.AssignedTranslatorAt,
		AssignedProofreaderAt: cre.TimedRoles.AssignedProofreaderAt,
		AssignedTypesetterAt:  cre.TimedRoles.AssignedTypesetterAt,
		AssignedRedrawerAt:    cre.TimedRoles.AssignedRedrawerAt,
		AssignedReviewerAt:    cre.TimedRoles.AssignedReviewerAt,
		AssignedPublisherAt:   cre.TimedRoles.AssignedPublisherAt,
		CreatedAt:             now,
		UpdatedAt:             now,
	}
}

// `AssignmentPutRow` maps assignment columns for put update.
type AssignmentPutRow struct {
	AssignedRawProviderAt *time.Time `gorm:"column:assigned_raw_provider_at"`
	AssignedTranslatorAt  *time.Time `gorm:"column:assigned_translator_at"`
	AssignedProofreaderAt *time.Time `gorm:"column:assigned_proofreader_at"`
	AssignedTypesetterAt  *time.Time `gorm:"column:assigned_typesetter_at"`
	AssignedRedrawerAt    *time.Time `gorm:"column:assigned_redrawer_at"`
	AssignedReviewerAt    *time.Time `gorm:"column:assigned_reviewer_at"`
	AssignedPublisherAt   *time.Time `gorm:"column:assigned_publisher_at"`

	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// `TableName` returns table name for `AssignmentPutRow`.
func (*AssignmentPutRow) TableName() string {
	return ASSIGNMENT_TABLE
}

// `NewAssignmentPutRowFromAggr` builds put update row from assignment put payload.
func NewAssignmentPutRowFromAggr(put *aggr.AssignmentPut) *AssignmentPutRow {
	return &AssignmentPutRow{
		AssignedRawProviderAt: put.TimedRoles.AssignedRawProviderAt,
		AssignedTranslatorAt:  put.TimedRoles.AssignedTranslatorAt,
		AssignedProofreaderAt: put.TimedRoles.AssignedProofreaderAt,
		AssignedTypesetterAt:  put.TimedRoles.AssignedTypesetterAt,
		AssignedRedrawerAt:    put.TimedRoles.AssignedRedrawerAt,
		AssignedReviewerAt:    put.TimedRoles.AssignedReviewerAt,
		AssignedPublisherAt:   put.TimedRoles.AssignedPublisherAt,
		UpdatedAt:             time.Now(),
	}
}
