package entity

import (
	"time"

	"poprako-s/internal/domain/model/aggr"
)

// `WORKSET_TABLE` is the table name for the workset entity.
const WORKSET_TABLE = "t_workset"

// `WorksetRow` maps a full workset record for read queries.
type WorksetRow struct {
	// `Id` is the primary key.
	Id string `gorm:"column:id;primaryKey"`

	// `TeamId` is the foreign key to the owning team.
	TeamId string `gorm:"column:team_id"`
	// `Team` is included when `includes` contains `team`.
	Team *TeamRow `gorm:"foreignKey:TeamId"`

	// `Index` is the team-scoped position of this workset.
	Index int `gorm:"column:index"`

	// `Name` is the display title.
	Name string `gorm:"column:name"`

	// `Desc` is the optional description.
	Desc *string `gorm:"column:description"`

	// `ComicCount` is the denormalised count of comics in this workset.
	ComicCount int `gorm:"column:comic_count"`

	// `CreatedAt` is the creation timestamp.
	CreatedAt time.Time `gorm:"column:created_at"`

	// `UpdatedAt` is the last-modification timestamp.
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// `TableName` returns the table name used by GORM for `WorksetRow`.
func (*WorksetRow) TableName() string {
	return WORKSET_TABLE
}

// `ToWorksetAggr` converts a `WorksetRow` into a `Workset` aggregate.
func (r *WorksetRow) ToWorksetAggr() *aggr.Workset {
	if r == nil {
		// Nil rows may appear in optional preload paths; return nil safely.
		return nil
	}

	var team *aggr.Team

	if r.Team != nil {
		team = r.Team.ToTeamAggr()
	}

	return &aggr.Workset{
		Id:         r.Id,
		TeamId:     r.TeamId,
		Team:       team,
		Index:      r.Index,
		Name:       r.Name,
		Desc:       r.Desc,
		ComicCount: r.ComicCount,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}
}

// `WorksetCreRow` is the write model used when inserting a new workset.
type WorksetCreRow struct {
	// `Id` is the generated primary key.
	Id string `gorm:"column:id;primaryKey"`

	// `TeamId` is the owning team's identifier.
	TeamId string `gorm:"column:team_id"`

	// `Index` is the team-scoped position.
	Index int `gorm:"column:index"`

	// `Name` is the display title.
	Name string `gorm:"column:name"`

	// `Desc` is the optional description.
	Desc *string `gorm:"column:description"`
}

// `TableName` returns the table name used by GORM for `WorksetCreRow`.
func (*WorksetCreRow) TableName() string {
	return WORKSET_TABLE
}

// `NewWorksetCreRowFromAggr` constructs a `WorksetCreRow` from a `WorksetCre` aggregate.
func NewWorksetCreRowFromAggr(cre *aggr.WorksetCre) *WorksetCreRow {
	return &WorksetCreRow{
		Id:     cre.Id,
		TeamId: cre.TeamId,
		Index:  cre.Index,
		Name:   cre.Name,
		Desc:   cre.Desc,
	}
}

// `WorksetUpdRow` is the write model used when updating an existing workset.
// Only the explicitly selected columns are written.
type WorksetUpdRow struct {
	// `Name` is the new display title.
	Name string `gorm:"column:name"`

	// `Desc` is the new optional description.
	Desc *string `gorm:"column:description"`

	// `UpdatedAt` is set to `now` at the time of the update.
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// `TableName` returns the table name used by GORM for `WorksetUpdRow`.
func (*WorksetUpdRow) TableName() string {
	return WORKSET_TABLE
}
