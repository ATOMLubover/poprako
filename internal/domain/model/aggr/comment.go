package aggr

import "time"

// `Comment` represents one team board comment and optional included relations.
type Comment struct {
	Id string

	TeamId string
	UserId string
	User   *User

	Content string

	CreatedAt time.Time
}

// `CommentCre` is the create payload for one team board comment.
type CommentCre struct {
	Id string

	TeamId string
	UserId string

	Content string
}
