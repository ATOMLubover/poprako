package val

// `ChapterExportVal` is the app-facing export object for one chapter.
type ChapterExportVal struct {
	ChapterId string `json:"chapter_id"`

	ChapterIndex int `json:"chapter_index"`

	ChapterSubtitle *string `json:"chapter_subtitle,omitempty"`

	ComicId string `json:"comic_id"`

	ComicTitle string `json:"comic_title"`

	Pages []PageExportVal `json:"pages"`
}

// `PageExportVal` is the app-facing export object for one page.
type PageExportVal struct {
	PageId string `json:"page_id"`

	PageIndex int `json:"page_index"`

	ImageUrl string `json:"image_url"`

	IsUploaded bool `json:"is_uploaded"`

	Units []UnitExportVal `json:"units"`
}

// `UnitExportVal` is the app-facing export object for one unit.
type UnitExportVal struct {
	UnitId string `json:"unit_id"`

	UnitIndex int `json:"unit_index"`

	PageId string `json:"page_id"`

	PageIndex int `json:"page_index"`

	XCoord float64 `json:"x_coord"`
	YCoord float64 `json:"y_coord"`

	IsBubble bool `json:"is_bubble"`

	TranslatedText *string `json:"translated_text"`

	TranslatorId *string `json:"translator_id"`

	TranslatorComment *string `json:"translator_comment"`

	IsProofread bool `json:"is_proofread"`

	ProofreadText *string `json:"proofread_text"`

	ProofreaderId *string `json:"proofreader_id"`

	ProofreaderComment *string `json:"proofreader_comment"`
}

// `ImportChapterArgs` holds import input for one chapter.
type ImportChapterArgs struct {
	ChapterId string `json:"chapter_id"`

	Format string `json:"format"`

	Content string `json:"content"`
}

// `ImportChapterBody` is the transport body for chapter import endpoint.
type ImportChapterBody struct {
	Format string `json:"format"`

	Content string `json:"content"`
}

// `ImportChapterRes` is the import summary returned to caller.
type ImportChapterRes struct {
	ImportedPageCount int `json:"imported_page_count"`

	ImportedUnitCount int `json:"imported_unit_count"`
}
