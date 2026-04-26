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
