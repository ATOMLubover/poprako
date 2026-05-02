package aggr

import "time"

// `Team` represents a translation team that owns worksets and members
type Team struct {
	Id string

	Name string
	Desc string

	AvatarKey      string
	AvatarUploaded bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

// `TeamCre` holds immutable fields for one team creation.
type TeamCre struct {
	Id string

	Name string
	Desc string
}

// `TeamUpd` holds mutable fields for one team put update.
type TeamUpd struct {
	Id string

	Name string
	Desc string
}
