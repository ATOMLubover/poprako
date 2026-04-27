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
