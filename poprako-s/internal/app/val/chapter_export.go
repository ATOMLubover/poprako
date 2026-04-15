package val

// ChapterExport 代表章节导出信息 VO
type ChapterExport struct {
	// ChapterID 表示章节 ID
	ChapterID string `json:"chapter_id"`
	// ChapterIndex 表示章节序号
	ChapterIndex int `json:"chapter_index"`
	// ChapterSubtitle 表示章节副标题
	ChapterSubtitle *string `json:"chapter_subtitle,omitempty"`

	// ComicID 表示漫画 ID
	ComicID string `json:"comic_id"`
	// ComicTitle 表示漫画标题
	ComicTitle string `json:"comic_title"`

	// Pages 表示章节中的页面列表
	Pages []PageExport `json:"pages"`
}

// PageExport 代表页面导出信息
type PageExport struct {
	// PageID 表示页面 ID
	PageID string `json:"page_id"`
	// PageIndex 表示页面序号
	PageIndex int `json:"page_index"`

	// ImageURL 表示页面图片地址
	ImageURL string `json:"image_url"`
	// IsUploaded 表示页面是否已上传
	IsUploaded bool `json:"is_uploaded"`

	// Units 表示页面中的翻译单元列表
	Units []UnitExport `json:"units"`
}

// UnitExport 代表翻译单元导出信息
type UnitExport struct {
	// UnitID 表示翻译单元 ID
	UnitID string `json:"unit_id"`
	// UnitIndex 表示翻译单元序号
	UnitIndex int `json:"unit_index"`

	// PageID 表示所属页面 ID
	PageID string `json:"page_id"`
	// PageIndex 表示所属页面序号
	PageIndex int `json:"page_index"`

	// XCoord 表示翻译单元的 X 坐标
	XCoord int `json:"x_coord"`
	// YCoord 表示翻译单元的 Y 坐标
	YCoord int `json:"y_coord"`

	// IsBubble 表示该单元是否是气泡框
	IsBubble bool `json:"is_bubble"`

	// TranslatedText 表示翻译后的文本
	TranslatedText *string `json:"translated_text"`
	// TranslatorID 表示翻译者的用户 ID
	TranslatorID *string `json:"translator_id"`
	// TranslatorComment 表示翻译者的备注
	TranslatorComment *string `json:"translator_comment"`

	// IsProofread 表示该单元是否已经校对
	IsProofread bool `json:"is_proofread"`
	// ProofreadText 表示校对后的文本
	ProofreadText *string `json:"proofread_text"`
	// ProofreaderID 表示校对者的用户 ID
	ProofreaderID *string `json:"proofreader_id"`
	// ProofreaderComment 表示校对者的备注
	ProofreaderComment *string `json:"proofreader_comment"`
}
