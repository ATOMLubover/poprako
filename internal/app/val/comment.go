package val

type CommentVal struct {
	Id string `json:"id"`

	TeamId string   `json:"team_id"`
	UserId string   `json:"user_id"`
	User   *UserVal `json:"user"`

	Content string `json:"content"`

	CreatedAt int64 `json:"created_at"`
}

type ListCommentArgs struct {
	TeamId string `url:"team_id"`

	Offset int `url:"offset"`
	Limit  int `url:"limit"`
}

type CreateCommentArgs struct {
	TeamId string `json:"team_id"`

	Content string `json:"content"`
}

type CommentCreatedRes struct {
	Id string `json:"id"`
}
