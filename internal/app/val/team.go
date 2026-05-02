package val

// `TeamVal` is the app-facing value object for `Team`.
type TeamVal struct {
	Id string `json:"id"`

	Name string `json:"name"`
	Desc string `json:"description"`

	AvatarUrl      string `json:"avatar_url"`
	AvatarUploaded bool   `json:"avatar_uploaded"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

// `TeamCreArgs` holds immutable fields for team creation.
type TeamCreArgs struct {
	Name string `json:"name"`
	Desc string `json:"description"`
}

// `TeamCreRes` returns the created team id.
type TeamCreRes struct {
	Id string `json:"id"`
}

// `TeamUpdArgs` holds mutable fields for team update.
type TeamUpdArgs struct {
	Id string `json:"id"`

	Name string `json:"name"`
	Desc string `json:"description"`
}

// `ListTeamArgs` carries list args for listing all teams.
type ListTeamArgs struct {
	Offset int `url:"offset"`
	Limit  int `url:"limit"`
}

// `ResvTeamAvatarArgs` holds parameters for reserving team avatar upload.
type ResvTeamAvatarArgs struct {
	TeamId  string `json:"team_id"`
	FileExt string `json:"file_extension"`
}

// `ResvTeamAvatarRes` returns the signed put url for avatar upload.
type ResvTeamAvatarRes struct {
	PutUrl string `json:"put_url"`
}
