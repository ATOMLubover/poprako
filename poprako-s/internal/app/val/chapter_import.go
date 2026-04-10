package val

// ImportChapterArgs 代表章节导入参数
// Content 为待导入的原始文本
// Format 仅支持 lp 或 prk
// ChapterID 由调用方指定
// 解析和鉴权由 app 层完成
type ImportChapterArgs struct {
	// ChapterID 表示目标章节 ID
	ChapterID string `json:"chapter_id" validate:"required"`
	// Format 表示导入格式
	Format string `json:"format" validate:"required"`
	// Content 表示导入文件内容
	Content string `json:"content" validate:"required"`
}

// ImportChapterRes 代表章节导入结果
// 统计值用于前端展示导入摘要
type ImportChapterRes struct {
	// ImportedPageCount 表示成功处理的页面数量
	ImportedPageCount int `json:"imported_page_count"`
	// ImportedUnitCount 表示成功写入的单元数量
	ImportedUnitCount int `json:"imported_unit_count"`
}