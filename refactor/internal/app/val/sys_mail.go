package val

// `SysMailVal` is the app-facing value object for one system mail.
type SysMailVal struct {
	Id string `json:"id"`

	Title   string `json:"title"`
	Content string `json:"content"`

	Read bool `json:"read"`

	CreatedAt int64 `json:"created_at"`
}

// `ListSysMailArgs` holds list arguments for current-user system mails.
type ListSysMailArgs struct {
	Offset int `url:"offset"`
	Limit  int `url:"limit"`
}
