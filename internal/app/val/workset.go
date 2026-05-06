package val

// `WorksetVal` is the app-facing value object for a `Workset`.
type WorksetVal struct {
	Id     string   `json:"id"`
	TeamId string   `json:"team_id"`
	Team   *TeamVal `json:"team,omitempty"`

	Index int     `json:"index"`
	Name  string  `json:"name"`
	Desc  *string `json:"description"`

	ComicCount int `json:"comic_count"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

type ListWorksetArgs struct {
	// `TeamId` is the identifier of the team whose worksets are to be listed.
	TeamId string `url:"team_id"`

	Offset int `url:"offset"`
	Limit  int `url:"limit"`
}

// `CreateWorksetArgs` holds the input for creating a new workset.
type CreateWorksetArgs struct {
	// `TeamId` is the owning team's identifier.
	TeamId string `json:"team_id"`

	// `Name` is the display title of the new workset.
	Name string `json:"name"`

	// `Desc` is an optional description. nil means use empty string.
	Desc *string `json:"description"`
}

// `WorksetCreatedRes` is returned after a successful workset creation.
type WorksetCreatedRes struct {
	// `Id` is the generated identifier for the new workset.
	Id string `json:"id"`
}

// `WorksetUpdArgs` holds the mutable fields for a workset update.
type WorksetUpdArgs struct {
	// `Id` identifies the workset to update.
	Id string `json:"id"`

	// `Name` is the new display title.
	Name string `json:"name"`

	// `Desc` is the new description for `PUT` semantics.
	// Nil means writing SQL `NULL`.
	Desc *string `json:"description"`
}
