package entity

import (
	"time"

	"poprako-s/internal/domain/model/aggr"
)

const TEAM_TABLE = "t_team"

type TeamRow struct {
	Id string `gorm:"column:id;primaryKey"`

	Name string `gorm:"column:name;unique"`
	Desc string `gorm:"column:description"`

	AvatarKey      string `gorm:"column:avatar_key"`
	AvatarUploaded bool   `gorm:"column:avatar_uploaded"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// `TeamCreRow` maps immutable columns for one team create operation.
type TeamCreRow struct {
	Id string `gorm:"column:id;primaryKey"`

	Name string `gorm:"column:name;unique"`
	Desc string `gorm:"column:description"`

	AvatarKey      string `gorm:"column:avatar_key"`
	AvatarUploaded bool   `gorm:"column:avatar_uploaded"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// `TableName` returns table name for `TeamCreRow`.
func (*TeamCreRow) TableName() string {
	return TEAM_TABLE
}

// `NewTeamCreRowFromAggr` converts one team create aggregate into insert row.
func NewTeamCreRowFromAggr(cre *aggr.TeamCre) *TeamCreRow {
	now := time.Now()

	return &TeamCreRow{
		Id:             cre.Id,
		Name:           cre.Name,
		Desc:           cre.Desc,
		AvatarKey:      "",
		AvatarUploaded: false,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// `TeamUpdRow` maps mutable columns for one team update operation.
type TeamUpdRow struct {
	Name string `gorm:"column:name"`
	Desc string `gorm:"column:description"`

	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (*TeamRow) TableName() string {
	return TEAM_TABLE
}

func (r *TeamRow) ToTeamAggr() *aggr.Team {
	if r == nil {
		// Do not treat nil as an error, as it may be used in includes.
		return nil
	}

	return &aggr.Team{
		Id: r.Id,

		Name: r.Name,
		Desc: r.Desc,

		AvatarKey:      r.AvatarKey,
		AvatarUploaded: r.AvatarUploaded,

		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}
