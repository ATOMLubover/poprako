package repository

import "labelplus-next-web-be/internal/infrastructure/repository/entity"

// 导出给上层使用的表名常量，禁止在 app 层写裸字符串。
const (
	UserTable       = entity.UserTable
	MemberTable     = entity.MemberTable
	TeamTable       = entity.TeamTable
	WorksetTable    = entity.WorksetTable
	ComicTable      = entity.ComicTable
	ChapterTable    = entity.ChapterTable
	PageTable       = entity.PageTable
	InvitationTable = entity.InvitationTable
	AssignmentTable = entity.AssignmentTable
)
