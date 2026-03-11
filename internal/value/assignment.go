package value

import (
	"errors"

	"labelplus-next-web-be/internal/domain/model"
)

// AssignmentWithUserInfo 用于 ListChapterAssignments 接口返回，含用户信息。
type AssignmentWithUserInfo struct {
	ID        string         `json:"id"`
	ChapterID string         `json:"chapter_id"`
	User      UserInfo       `json:"user"`
	Roles     model.RoleMask `json:"roles"`
}

func NewAssignmentWithUserInfo(assignment model.AssignmentWithUserInfo, avatarURL string) AssignmentWithUserInfo {
	return AssignmentWithUserInfo{
		ID:        assignment.ID,
		ChapterID: assignment.ChapterID,
		User:      NewUserInfoFromModel(assignment.User, avatarURL),
		Roles:     assignment.RoleMask(),
	}
}

// AssignmentWithChapterInfo 用于 ListUserAssignments 接口返回，含章节+漫画信息。
type AssignmentWithChapterInfo struct {
	ID      string               `json:"id"`
	Chapter ChapterWithComicInfo `json:"chapter"`
	UserID  string               `json:"user_id"`
	Roles   model.RoleMask       `json:"roles"`
}

func NewAssignmentWithChapterInfo(assignment model.AssignmentWithChapterInfo, coverURL string) AssignmentWithChapterInfo {
	return AssignmentWithChapterInfo{
		ID:      assignment.ID,
		Chapter: NewChapterWithComicInfoFromModel(assignment.Chapter, coverURL),
		UserID:  assignment.UserID,
		Roles:   assignment.RoleMask(),
	}
}

// ListChapterAssignmentArgs 用于查询某章节的所有分配。
type ListChapterAssignmentArgs struct {
	ChapterID string `url:"chapter_id"`
}

func (a *ListChapterAssignmentArgs) Validate() error {
	if a == nil {
		return errors.New("参数不能为空")
	}

	if a.ChapterID == "" {
		return errors.New("章节 ID 不能为空")
	}

	return nil
}

// ListUserAssignmentArgs 用于查询当前用户的所有分配，支持分页。
type ListUserAssignmentArgs struct {
	PaginationParams
}

func (a *ListUserAssignmentArgs) Validate() error {
	if a == nil {
		return errors.New("参数不能为空")
	}

	return a.PaginationParams.Validate()
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

func NewCreateChapterAssignmentResult(id string) CreateChapterAssignmentResult {
	return CreateChapterAssignmentResult{ID: id}
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
