package app

import (
	"context"

	"poprako-s/internal/app/val"
)

// ChapterImportApp 定义章节导入能力
// 支持 Poprako JSON 和 LabelPlus 文本两种格式
// 导入前会校验当前用户在章节中的分配权限
// 导入后返回处理摘要
type ChapterImportApp interface {
	// ImportChapter 导入章节数据
	ImportChapter(
		cx context.Context,
		currUserID string,
		args *val.ImportChapterArgs,
	) (*val.ImportChapterRes, error)
}
