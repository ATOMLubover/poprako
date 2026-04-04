package val

import "poprako-s/internal/domain/model"

// ComicInfo 表示用于在应用层和接口层传递的漫画信息 VO
type ComicInfo struct {
	// ID 是漫画的唯一标识
	ID string `json:"id"`

	// WorksetID 是所属作品集 ID
	WorksetID string `json:"workset_id"`
	// Index 是漫画在作品集内的序号
	Index int `json:"index"`

	// Title 是漫画标题
	Title string `json:"title"`
	// Author 是漫画作者
	Author string `json:"author"`
	// Description 是漫画描述
	Description string `json:"description"`

	// ChapterCount 是漫画下章节数量
	ChapterCount int `json:"chapter_count"`

	// CoverURL 是漫画封面的可访问地址
	CoverURL string `json:"cover_url"`
	// IsCoverUploaded 表示封面是否已经上传
	IsCoverUploaded bool `json:"is_cover_uploaded"`

	// CreatorID 是漫画创建者 ID
	CreatorID string `json:"creator_id"`
	// Creator 是可选的创建者信息（仅在 includes 时填充）
	Creator *UserInfo `json:"creator,omitempty"`
	// Workset 是可选的作品集信息（仅在 includes 时填充）
	Workset *WorksetInfo `json:"workset,omitempty"`

	// LastActiveAt 是最近活跃时间的 Unix 毫秒时间戳
	LastActiveAt int64 `json:"last_active_at"`
	// CreatedAt 是记录创建时间的 Unix 毫秒时间戳
	CreatedAt int64 `json:"created_at"`
	// UpdatedAt 是记录最近一次更新时间的 Unix 毫秒时间戳
	UpdatedAt int64 `json:"updated_at"`
}

// ListComicArgs 表示列出漫画列表请求的参数
type ListComicArgs struct {
	// WorksetID 是目标作品集 ID
	WorksetID string `json:"workset_id" validate:"required"`
	// FuzzyTitle 是模糊搜索标题
	FuzzyTitle string `json:"fuzzy_title"`

	// 进度筛选
	UploadStatus    *model.WorkflowPhase `json:"upload_status"`
	TranslateStatus *model.WorkflowPhase `json:"translate_status"`
	ProofreadStatus *model.WorkflowPhase `json:"proofread_status"`
	TypesetStatus   *model.WorkflowPhase `json:"typeset_status"`
	ReviewStatus    *model.WorkflowPhase `json:"review_status"`
	PublishStatus   *model.WorkflowPhase `json:"publish_status"`

	// Includes 指定查询时要包含的反向单射数据
	Includes []model.ComicInclude `json:"includes"`

	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// CreateComicArgs 表示创建漫画请求的参数
type CreateComicArgs struct {
	// WorksetID 是目标作品集 ID
	WorksetID string `json:"workset_id" validate:"required"`
	// Title 是漫画标题
	Title string `json:"title" validate:"required"`
	// Author 是漫画作者
	Author string `json:"author" validate:"required"`
	// Description 是漫画描述
	Description string `json:"description"`
}

// CreateComicRes 表示创建漫画成功后的响应数据
type CreateComicRes struct {
	// ID 是新创建漫画的标识
	ID string `json:"id"`
}

// UpdateComicArgs 表示更新漫画请求的参数
type UpdateComicArgs struct {
	// ID 是要更新的漫画标识
	ID string `json:"id" validate:"required"`
	// Title 是更新后的标题
	Title string `json:"title" validate:"required"`
	// Author 是更新后的作者
	Author string `json:"author" validate:"required"`
	// Description 是更新后的描述
	Description string `json:"description"`
}

// ReserveComicCoverArgs 表示预留漫画封面上传接口的请求参数
type ReserveComicCoverArgs struct {
	// ComicID 是要上传封面的漫画标识
	ComicID string `json:"comic_id" validate:"required"`
	// FileName 是封面文件的原始名称，主要用于 OSS 存储时保留扩展名
	// 需要携带文件扩展名以便 OSS 正确识别文件类型，例如 "cover.jpg"
	FileName string `json:"file_name" validate:"required"`
}

// ReserveComicCoverRes 表示预留漫画封面上传接口的响应数据
type ReserveComicCoverRes struct {
	// PutURL 是用于上传封面的预签名 URL，客户端可以直接使用该 URL 上传封面文件
	PutURL string `json:"put_url"`
}
