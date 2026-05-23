package val

type AnnouncementVal struct {
	Id string `json:"id"`

	TeamId string   `json:"team_id"`
	UserId string   `json:"user_id"`
	User   *UserVal `json:"user"`

	Title   string `json:"title"`
	Content string `json:"content"`

	CreatedAt int64 `json:"created_at"`
}

type ListAnnouncementArgs struct {
	TeamId string `url:"team_id"`

	Offset int `url:"offset"`
	Limit  int `url:"limit"`
}

type CreateAnnouncementArgs struct {
	TeamId string `json:"team_id"`

	Title   string `json:"title"`
	Content string `json:"content"`
}

type AnnouncementCreatedRes struct {
	Id string `json:"id"`
}
