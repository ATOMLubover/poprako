package val

import "poprako-s/internal/domain/model"

// AssignmentInfo 表示用于在应用层和接口层传递的分配信息 VO
type AssignmentInfo struct {
	// ID 是分配记录的唯一标识
	ID string `json:"id"`

	// ChapterID 是所属章节 ID
	ChapterID string `json:"chapter_id"`
	// UserID 是被分配用户的 ID
	UserID string `json:"user_id"`

	// Roles 是该分配包含的角色掩码
	Roles model.RoleMask `json:"roles"`

	// Chapter 是可选的章节信息（仅在 includes 时填充）
	Chapter *ChapterInfo `json:"chapter,omitempty"`
	// User 是可选的用户信息（仅在 includes 时填充）
	User *UserInfo `json:"user,omitempty"`

	// CreatedAt 是记录创建时间的 Unix 毫秒时间戳
	CreatedAt int64 `json:"created_at"`
	// UpdatedAt 是记录最近一次更新时间的 Unix 毫秒时间戳
	UpdatedAt int64 `json:"updated_at"`
}

// ListChapterAssignmentArgs 表示列出章节分配请求的参数
type ListChapterAssignmentArgs struct {
	// ChapterID 是目标章节 ID
	ChapterID string `json:"chapter_id" validate:"required"`
	Offset    int    `json:"offset"`
	Limit     int    `json:"limit"`
}

// ListMyAssignmentArgs 表示列出当前用户分配请求的参数
type ListMyAssignmentArgs struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// CreateAssignmentArgs 表示创建分配请求的参数
type CreateAssignmentArgs struct {
	// ChapterID 是目标章节 ID
	ChapterID string `json:"chapter_id" validate:"required"`
	// UserID 是被分配的用户 ID
	UserID string `json:"user_id" validate:"required"`
	// Roles 是分配的角色掩码
	Roles model.RoleMask `json:"roles" validate:"required"`
}

// CreateAssignmentRes 表示创建分配成功后的响应数据
type CreateAssignmentRes struct {
	// ID 是新创建分配记录的标识
	ID string `json:"id"`
}

// UpdateAssignmentArgs 表示更新分配请求的参数
type UpdateAssignmentArgs struct {
	// ID 是要更新的分配记录标识
	ID string `json:"id" validate:"required"`
	// Roles 是更新后的角色掩码（PUT 语义全量替换）
	Roles model.RoleMask `json:"roles" validate:"required"`
}
