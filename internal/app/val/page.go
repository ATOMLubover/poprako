package val

// `PageVal` is the app-facing value object for page data.
type PageVal struct {
	Id string `json:"id"`

	ChapterId string `json:"chapter_id"`

	Index int `json:"index"`

	ImageUrl string `json:"image_url"`

	ImageUploaded bool `json:"image_uploaded"`

	TotalUnitCount int `json:"total_unit_count"`

	TranslatedUnitCount int `json:"translated_unit_count"`

	ProofreadUnitCount int `json:"proofread_unit_count"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

// `ResvChapterPagesArgs` holds input for reserving chapter pages.
type ResvChapterPagesArgs struct {
	ChapterId string `json:"chapter_id"`

	PageCount int `json:"page_count"`

	FileExt string `json:"file_extension"`
}

// `ResvChapterPagesRes` holds created pages and their put urls.
type ResvChapterPagesRes struct {
	Creations []PageCreationRes `json:"creations,omitempty"`
}

// `PageCreationRes` holds one reserved page upload result.
type PageCreationRes struct {
	PageId string `json:"page_id"`

	PutUrl string `json:"put_url"`
}

// `ListChapterPageArgs` carries query args for page list API.
type ListChapterPageArgs struct {
	ChapterId string `url:"chapter_id"`

	Offset int `url:"offset"`
	Limit  int `url:"limit"`
}

// `MarkPageImageUploadedArgs` holds input for upload confirmation.
type MarkPageImageUploadedArgs struct {
	PageId string `json:"page_id"`
}
