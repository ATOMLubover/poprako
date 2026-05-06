package val

import "poprako-s/internal/domain/model/enum"

// `ComicVal` is the app-facing value object for comic data
// It keeps nullable description semantics and timestamp millisecond format
// Relation fields are omitted until include-based assembly is introduced
type ComicVal struct {
	Id string `json:"id"`

	WorksetId string      `json:"workset_id"`
	Workset   *WorksetVal `json:"workset,omitempty"`

	Index int `json:"index"`

	Title         string  `json:"title"`
	Author        string  `json:"author"`
	Desc          *string `json:"description,omitempty"`
	IsCompleted   bool    `json:"is_completed"`
	CoverUrl      string  `json:"cover_url"`
	CoverUploaded bool    `json:"cover_uploaded"`

	ChapterCount int `json:"chapter_count"`

	CreatorId string   `json:"creator_id"`
	Creator   *UserVal `json:"creator,omitempty"`

	LastActiveAt int64 `json:"last_active_at"`
	CreatedAt    int64 `json:"created_at"`
	UpdatedAt    int64 `json:"updated_at"`
}

// `ListComicArgs` carries query args for comic list API
// `WorksetId` comes from route path and pagination is explicit
// `FuzzyTitle` is optional and empty means no fuzzy filter
type ListComicArgs struct {
	// `WorksetId` is the target workset identifier
	WorksetId string `url:"workset_id"`

	// `FuzzyTitle` is optional fuzzy title keyword
	FuzzyTitle string `url:"fuzzy_title"`

	// Workflow phase filters based on pinned chapter progress.
	UploadPhase *enum.WorkflowPhase `url:"upload_phase"`

	TranslatePhase *enum.WorkflowPhase `url:"translate_phase"`

	ProofreadPhase *enum.WorkflowPhase `url:"proofread_phase"`

	TypesetPhase *enum.WorkflowPhase `url:"typeset_phase"`

	ReviewPhase *enum.WorkflowPhase `url:"review_phase"`

	PublishPhase *enum.WorkflowPhase `url:"publish_phase"`

	Includes []enum.ComicIncl `url:"includes"`

	Offset int `url:"offset"`
	Limit  int `url:"limit"`
}

// `CreateComicArgs` holds input for creating a comic
// `WorksetId` is required because comic belongs to exactly one workset
// `Desc` keeps nullable semantics and nil means SQL NULL
// These args are used by both app and HTTP transport mapping
type CreateComicArgs struct {
	// `WorksetId` is the owning workset id
	WorksetId string `json:"workset_id"`

	// `Title` is the comic title
	Title string `json:"title"`
	// `Author` is the comic author
	Author string `json:"author"`
	// `Desc` is optional description
	Desc *string `json:"description"`
}

// `ComicCreatedRes` is returned after successful comic creation
// It only exposes generated id for follow-up queries
// Response shape stays stable with existing app response wrapper
type ComicCreatedRes struct {
	// `Id` is the generated comic identifier
	Id string `json:"id"`
}

// `GetComicByIdArgs` carries query args for comic detail API.
type GetComicByIdArgs struct {
	// `ComicId` identifies the target comic.
	ComicId string `url:"comic_id"`

	// `Includes` controls relation assembly fields.
	Includes []enum.ComicIncl `url:"includes"`
}

// `ComicUpdArgs` holds mutable fields for comic put update
// Nil `Desc` means writing SQL NULL
// `Id` is bound from path by HTTP layer
type ComicUpdArgs struct {
	// `Id` identifies the target comic
	Id string `json:"id"`

	// `Title` is new comic title
	Title string `json:"title"`
	// `Author` is new comic author
	Author string `json:"author"`
	// `Desc` is new optional description
	Desc *string `json:"description"`
}

// `ResvComicCoverArgs` holds parameters for reserving comic cover upload.
type ResvComicCoverArgs struct {
	ComicId string `json:"comic_id"`
	FileExt string `json:"file_extension"`
}

// `ResvComicCoverBody` is the transport body for comic-cover reservation endpoint.
type ResvComicCoverBody struct {
	FileExt string `json:"file_extension"`
}

// `ResvComicCoverRes` returns signed put url for cover upload.
type ResvComicCoverRes struct {
	PutUrl string `json:"put_url"`
}
