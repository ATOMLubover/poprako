package aggr

import (
	"fmt"
	"time"
)

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

// `GenAvatarKey` returns the OSS object key for the team avatar with the given file extension.
func (t *Team) GenAvatarKey(ext string) string {
	return fmt.Sprintf("team_avatar/%s.%s", t.Id, ext)
}
