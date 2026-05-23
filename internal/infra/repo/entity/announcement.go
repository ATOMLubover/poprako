package entity

import (
	"time"

	"poprako-s/internal/domain/model/aggr"
)

// `ANNOUNCEMENT_TABLE` is the table name for team announcements.
const ANNOUNCEMENT_TABLE = "t_announcement"

// `AnnouncementRow` maps one announcement record for read queries.
type AnnouncementRow struct {
	Id string `gorm:"column:id;primaryKey"`

	TeamId string   `gorm:"column:team_id"`
	UserId string   `gorm:"column:user_id"`
	User   *UserRow `gorm:"foreignKey:UserId"`

	Title   string `gorm:"column:title"`
	Content string `gorm:"column:content"`

	CreatedAt time.Time `gorm:"column:created_at"`
}

// `TableName` returns the table name of `AnnouncementRow`.
func (*AnnouncementRow) TableName() string {
	return ANNOUNCEMENT_TABLE
}

// `ToAnnouncementAggr` converts a row into `Announcement` aggregate.
func (r *AnnouncementRow) ToAnnouncementAggr() *aggr.Announcement {
	if r == nil {
		return nil
	}

	var user *aggr.User
	if r.User != nil {
		user = r.User.ToUserAggr()
	}

	return &aggr.Announcement{
		Id: r.Id,

		TeamId: r.TeamId,
		UserId: r.UserId,
		User:   user,

		Title:   r.Title,
		Content: r.Content,

		CreatedAt: r.CreatedAt,
	}
}

// `AnnouncementCreRow` maps one announcement record for create queries.
type AnnouncementCreRow struct {
	Id string `gorm:"column:id;primaryKey"`

	TeamId string `gorm:"column:team_id"`
	UserId string `gorm:"column:user_id"`

	Title   string `gorm:"column:title"`
	Content string `gorm:"column:content"`

	CreatedAt time.Time `gorm:"column:created_at"`
}

// `TableName` returns the table name of `AnnouncementCreRow`.
func (*AnnouncementCreRow) TableName() string {
	return ANNOUNCEMENT_TABLE
}

// `NewAnnouncementCreRowFromAggr` builds a create row from `AnnouncementCre`.
func NewAnnouncementCreRowFromAggr(cre *aggr.AnnouncementCre) *AnnouncementCreRow {
	return &AnnouncementCreRow{
		Id: cre.Id,

		TeamId: cre.TeamId,
		UserId: cre.UserId,

		Title:   cre.Title,
		Content: cre.Content,

		CreatedAt: time.Now(),
	}
}
