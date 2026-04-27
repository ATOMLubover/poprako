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

// `TeamUpdArgs` holds mutable fields for team update.
type TeamUpdArgs struct {
	Id string `json:"id"`

	Name string `json:"name"`
	Desc string `json:"description"`
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
