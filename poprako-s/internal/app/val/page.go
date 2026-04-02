package val

// PageInfo 表示用于在应用层和接口层传递的页面信息 VO
type PageInfo struct {
	// ID 是页面的唯一标识
	ID string `json:"id"`

	// ChapterID 是所属章节 ID
	ChapterID string `json:"chapter_id"`
	// Index 是页面在章节中的序号
	Index int `json:"index"`

	// ImageURL 是页面图片的可访问地址
	ImageURL string `json:"image_url"`
	// IsUploaded 表示图片是否已经上传
	IsUploaded bool `json:"is_uploaded"`

	// CreatorID 是页面创建者 ID
	CreatorID string `json:"creator_id"`
	// Creator 是可选的创建者信息（仅在 includes 时填充）
	Creator *UserInfo `json:"creator,omitempty"`

	// TotalUnitCount 是页面的翻译单元总数
	TotalUnitCount int `json:"total_unit_count"`
	// TranslatedUnitCount 是已翻译的单元数量
	TranslatedUnitCount int `json:"translated_unit_count"`
	// ProofreadUnitCount 是已校对的单元数量
	ProofreadUnitCount int `json:"proofread_unit_count"`

	// CreatedAt 是记录创建时间的 Unix 毫秒时间戳
	CreatedAt int64 `json:"created_at"`
	// UpdatedAt 是记录最近一次更新时间的 Unix 毫秒时间戳
	UpdatedAt int64 `json:"updated_at"`
}

// ReserveChapterPagesArgs 表示预留章节页面请求的参数
type ReserveChapterPagesArgs struct {
	// ChapterID 是目标章节 ID
	ChapterID string `json:"chapter_id" validate:"required"`
	// PageCount 是要预留的页面数量
	PageCount int `json:"page_count" validate:"required"`
	// Extension 是图片文件扩展名
	Extension string `json:"extension" validate:"required"`
}

// ReserveChapterPagesRes 表示预留章节页面成功后的响应数据
type ReserveChapterPagesRes struct {
	// Creations 是每个页面的创建结果
	Creations []PageCreationResult `json:"creations"`
}

// PageCreationResult 表示单个页面预留创建的结果
type PageCreationResult struct {
	// PageID 是新创建页面的标识
	PageID string `json:"page_id"`
	// PutURL 是用于上传的预签名 URL
	PutURL string `json:"put_url"`
}

// ListChapterPageArgs 表示列出章节页面请求的参数
type ListChapterPageArgs struct {
	// ChapterID 是目标章节 ID
	ChapterID string `json:"chapter_id" validate:"required"`
	Offset    int    `json:"offset"`
	Limit     int    `json:"limit"`
}

// UpdatePageArgs 表示更新页面请求的参数
type UpdatePageArgs struct {
	// ID 是要更新的页面标识
	ID string `json:"id" validate:"required"`
	// IsUploaded 表示是否标记为已上传
	IsUploaded bool `json:"is_uploaded"`
}
