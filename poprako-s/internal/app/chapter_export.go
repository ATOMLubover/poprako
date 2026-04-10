package app

import (
	"context"

	"poprako-s/internal/app/val"
)

// ChapterExportApp 定义了章节数据导出相关的应用服务接口
type ChapterExportApp interface {
	// ExportChapter 导出章节数据，包含章节信息、页面信息和翻译单元信息
	// 调用方须在该章节中拥有分配记录，否则返回权限不足错误
	ExportChapter(
		cx context.Context,
		currUserID string,
		chapterID string,
	) (*val.ChapterExport, error)

	// ExportChapterLp 以 LabelPlus 格式导出章节数据
	// 调用方须在该章节中拥有分配记录，否则返回权限不足错误
	// 返回 LabelPlus 格式的纯文本内容
	ExportChapterLp(
		cx context.Context,
		currUserID string,
		chapterID string,
	) (string, error)
}
