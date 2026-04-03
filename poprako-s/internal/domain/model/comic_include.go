package model

// ComicInclude 定义了在查询漫画列表时可以包含的 **反向** 单射数据
type ComicInclude string

const (
	ComicIncludeWorkset ComicInclude = "workset"
	ComicIncludeTeam    ComicInclude = "workset.team"
	ComicIncludeCreator ComicInclude = "creator"
)
