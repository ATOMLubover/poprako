package value

import (
	"errors"

	"labelplus-next-web-be/internal/domain/model"
)

// 标准的章节分配信息，可以递归地添加 include 字段
type AssignmentInfo struct {
	ID string `json:"id"`

	// 关联的用户信息
	UserID string    `json:"user_id"`
	User   *UserInfo `json:"user,omitempty"`

	// 关联的章节信息
	ChapterID string       `json:"chapter_id"`
	Chapter   *ChapterInfo `json:"chapter,omitempty"`

	AssignedRawProviderAt *int64 `json:"assigned_raw_provider_at,omitempty"`
	AssignedTranslatorAt  *int64 `json:"assigned_translator_at,omitempty"`
	AssignedProofreaderAt *int64 `json:"assigned_proofreader_at,omitempty"`
	AssignedTypesetterAt  *int64 `json:"assigned_typesetter_at,omitempty"`
	AssignedRedrawerAt    *int64 `json:"assigned_redrawer_at,omitempty"`
	AssignedReviewerAt    *int64 `json:"assigned_reviewer_at,omitempty"`
	AssignedPublisherAt   *int64 `json:"assigned_publisher_at,omitempty"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

// ListAssignmentArgs 用于查询当前用户的所有分配，支持分页。
type ListAssignmentArgs struct {
	// 嵌套关联字段
	Includes []string `url:"includes"`

	// 查询条件
	ChapterID string `url:"chapter_id,omitempty"`
	UserID    string `url:"user_id,omitempty"`

	PaginationParams
}

func (a *ListAssignmentArgs) Validate() error {
	if a == nil {
		return errors.New("参数不能为空")
	}

	if err := a.PaginationParams.Validate(); err != nil {
		return err
	}

	if len(a.Includes) > 0 {
	}

	return nil
}

type CreateChapterAssignmentArgs struct {
	ChapterID string         `json:"chapter_id"`
	UserID    string         `json:"user_id"`
	Role      model.RoleMask `json:"role"`
}

func (a *CreateChapterAssignmentArgs) Validate() error {
	if a == nil {
		return errors.New("参数不能为空")
	}

	if a.ChapterID == "" {
		return errors.New("章节 ID 不能为空")
	}

	if a.UserID == "" {
		return errors.New("用户 ID 不能为空")
	}

	if a.Role == 0 {
		return errors.New("分工角色不能为空")
	}

	return nil
}

type CreateChapterAssignmentResult struct {
	ID string `json:"id"`
}

// UpdateAssignmentArgs 采用 PUT 语义的全量替换。
type UpdateAssignmentArgs struct {
	ID   string         `json:"id"`
	Role model.RoleMask `json:"role"`
}

func (a *UpdateAssignmentArgs) Validate() error {
	if a == nil {
		return errors.New("参数不能为空")
	}

	if a.ID == "" {
		return errors.New("分配 ID 不能为空")
	}

	if a.Role == 0 {
		return errors.New("分工角色不能为空")
	}

	return nil
}

// ListChapterAssignmentArgs 查询某章节的所有分配
type ListChapterAssignmentArgs struct {
	ChapterID string   `url:"chapter_id"`
	Includes  []string `url:"includes[]"`
	PaginationParams
}

func (a *ListChapterAssignmentArgs) Validate() error {
	if a == nil {
		return errors.New("参数不能为空")
	}

	if a.ChapterID == "" {
		return errors.New("章节 ID 不能为空")
	}

	if err := a.PaginationParams.Validate(); err != nil {
		return errors.New("分页参数无效: " + err.Error())
	}

	return nil
}
