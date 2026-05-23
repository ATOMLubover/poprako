package entity

import (
	"time"

	"poprako-s/internal/domain/model/aggr"
)

// `COMMENT_TABLE` is the table name for team comments.
const COMMENT_TABLE = "t_comment"

// `CommentRow` maps one team comment record for read queries.
type CommentRow struct {
	Id string `gorm:"column:id;primaryKey"`

	TeamId string   `gorm:"column:team_id"`
	UserId string   `gorm:"column:user_id"`
	User   *UserRow `gorm:"foreignKey:UserId"`

	Content string `gorm:"column:content"`

	CreatedAt time.Time `gorm:"column:created_at"`
}

// `TableName` returns the table name of `CommentRow`.
func (*CommentRow) TableName() string {
	return COMMENT_TABLE
}

// `ToCommentAggr` converts a row into `Comment` aggregate.
func (r *CommentRow) ToCommentAggr() *aggr.Comment {
	if r == nil {
		return nil
	}

	var user *aggr.User
	if r.User != nil {
		user = r.User.ToUserAggr()
	}

	return &aggr.Comment{
		Id: r.Id,

		TeamId: r.TeamId,
		UserId: r.UserId,
		User:   user,

		Content: r.Content,

		CreatedAt: r.CreatedAt,
	}
}

// `CommentCreRow` maps one team comment record for create queries.
type CommentCreRow struct {
	Id string `gorm:"column:id;primaryKey"`

	TeamId string `gorm:"column:team_id"`
	UserId string `gorm:"column:user_id"`

	Content string `gorm:"column:content"`

	CreatedAt time.Time `gorm:"column:created_at"`
}

// `TableName` returns the table name of `CommentCreRow`.
func (*CommentCreRow) TableName() string {
	return COMMENT_TABLE
}

// `NewCommentCreRowFromAggr` builds a create row from `CommentCre`.
func NewCommentCreRowFromAggr(cre *aggr.CommentCre) *CommentCreRow {
	return &CommentCreRow{
		Id: cre.Id,

		TeamId: cre.TeamId,
		UserId: cre.UserId,

		Content: cre.Content,

		CreatedAt: time.Now(),
	}
}
