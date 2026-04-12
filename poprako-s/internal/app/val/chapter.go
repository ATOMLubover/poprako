package val

import "poprako-s/internal/domain/model"

// ChapterInfo 表示用于在应用层和接口层传递的章节信息 VO
type ChapterInfo struct {
	// ID 是章节的唯一标识
	ID string `json:"id"`

	// ComicID 是所属漫画 ID
	ComicID string `json:"comic_id"`
	// TODO
	Comic *ComicInfo `json:"comic,omitempty"`
	// IsPinned 表示章节是否被顶置
	IsPinned bool `json:"is_pinned"`

	// Index 是章节在漫画中的序号
	Index int `json:"index"`
	// Subtitle 是章节副标题
	Subtitle string `json:"subtitle"`

	// PageCount 是章节的页面数量
	PageCount int `json:"page_count"`
	// TotalUnitCount 是章节的翻译单元总数
	TotalUnitCount int `json:"total_unit_count"`
	// TranslatedUnitCount 是已翻译的单元数量
	TranslatedUnitCount int `json:"translated_unit_count"`
	// ProofreadUnitCount 是已校对的单元数量
	ProofreadUnitCount int `json:"proofread_unit_count"`

	// CreatorID 是章节创建者 ID
	CreatorID string `json:"creator_id"`
	// Creator 是可选的创建者信息（仅在 includes 时填充）
	Creator *UserInfo `json:"creator,omitempty"`

	// 工作流时间戳（Unix 毫秒），nil 表示尚未触发
	UploadedAt     *int64 `json:"uploaded_at"`
	TransalatingAt *int64 `json:"translating_at"`
	TranslatedAt   *int64 `json:"translated_at"`
	ProofreadingAt *int64 `json:"proofreading_at"`
	ProofreadAt    *int64 `json:"proofread_at"`
	TypesettingAt  *int64 `json:"typesetting_at"`
	TypesetAt      *int64 `json:"typeset_at"`
	ReviewedAt     *int64 `json:"reviewed_at"`
	PublishedAt    *int64 `json:"published_at"`

	// CreatedAt 是记录创建时间的 Unix 毫秒时间戳
	CreatedAt int64 `json:"created_at"`
	// UpdatedAt 是记录最近一次更新时间的 Unix 毫秒时间戳
	UpdatedAt int64 `json:"updated_at"`
}

// ListChapterArgs 表示列出章节列表请求的参数
type ListChapterArgs struct {
	// ComicID 是目标漫画 ID
	ComicID string `json:"comic_id" url:"comic_id" validate:"required"`
	Offset  int    `json:"offset" url:"offset"`
	Limit   int    `json:"limit" url:"limit"`
}

// CreateChapterArgs 表示创建章节请求的参数
type CreateChapterArgs struct {
	// ComicID 是目标漫画 ID
	ComicID string `json:"comic_id" validate:"required"`
	// Subtitle 是章节副标题（可选）
	Subtitle *string `json:"subtitle"`
}

// CreateChapterRes 表示创建章节成功后的响应数据
type CreateChapterRes struct {
	// ID 是新创建章节的标识
	ID string `json:"id"`
}

// UpdateChapterArgs 表示更新章节请求的参数
type UpdateChapterArgs struct {
	// ChapterID 是要更新的章节标识
	ChapterID string `json:"chapter_id" validate:"required"`
	// Subtitle 是更新后的副标题（可选）
	Subtitle *string `json:"subtitle"`
	// IsPinned 是更新后的顶置状态（可选）
	IsPinned *bool `json:"is_pinned"`

	// 工作流转换事件（可选）
	WorkflowTransition *model.WorkflowTransition `json:"workflow_transition"`
}

// InviteChapterAssigneeArgs 表示邀请章节协助者的请求参数
type InviteChapterAssigneeArgs struct {
	// ChapterID 是目标章节 ID
	ChapterID string `json:"chapter_id" validate:"required"`
	// InviteeQQ 是被邀请者的 QQ 号码
	// 因为可能有跨组邀请，因此用户只能通过 QQ 指定
	InviteeQQ string `json:"invitee_qq" validate:"required"`
	// Roles 是被邀请者的角色
	Roles model.RoleMask `json:"role" validate:"required"`
}

// InviteChapterAssigneeRes 表示邀请章节协助者成功后的响应数据
type InviteChapterAssigneeRes struct {
	// InvCode 是邀请代码，供被邀请者使用
	InvCode string `json:"invitation_code"`
}
