package value

import "labelplus-next-web-be/internal/domain/model"

type AssignmentWithUserInfo struct {
	ID        string         `json:"id"`
	ChapterID string         `json:"chapter_id"`
	User      UserInfo       `json:"user"`
	Roles     model.RoleMask `json:"roles"`
}

type AssignmentWithChapterInfo struct {
	ID      string         `json:"id"`
	Chapter ChapterDetail  `json:"chapter"`
	UserID  string         `json:"user_id"`
	Roles   model.RoleMask `json:"roles"`
}

type CreateChapterAssignmentArgs struct {
	ChapterID string         `json:"chapter_id"`
	UserID    string         `json:"user_id"`
	Role      model.RoleMask `json:"role"`
}

type CreateChapterAssignmentResult struct {
	ID string `json:"id"`
}

type UpdateAssignmentArgs struct {
	ID   string         `json:"id"`
	Role model.RoleMask `json:"role"`
}
