package val

type CommentVal struct {
	Id string `json:"id"`

	TeamId string `json:"team_id"`
	UserId string `json:"user_id"`

	Content string `json:"content"`

	CreatedAt int64 `json:"created_at"`
}

type CreateCommentArgs struct {
	TeamId string `json:"team_id"`

	Content string `json:"content"`
}
