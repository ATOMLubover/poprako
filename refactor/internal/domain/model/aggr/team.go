package aggr

import "time"

type Team struct {
	Id string

	Name string
	Desc string

	AvatarKey      string
	AvatarUploaded bool

	CreatedAt time.Time
	UpdatedAt time.Time
}
