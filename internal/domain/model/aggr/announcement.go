package aggr

import "time"

// `Announcement` represents one team announcement and optional included relations.
type Announcement struct {
	Id string

	TeamId string
	UserId string
	User   *User

	Title   string
	Content string

	CreatedAt time.Time
}

// `AnnouncementCre` is the create payload for one team announcement.
type AnnouncementCre struct {
	Id string

	TeamId string
	UserId string

	Title   string
	Content string
}
