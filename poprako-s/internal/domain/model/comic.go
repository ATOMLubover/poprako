package model

import (
	"fmt"
	"time"
)

// ComicInfo 包含作品集/漫画的元信息，用于聚合其章节和统计数据
type ComicInfo struct {
	ID string

	WorksetID string
	// Workset 仅在 includes 指定时填充
	Workset *WorksetInfo

	// 在作品集内部的序号
	Index       int
	Title       string
	Author      string
	Description string

	ChapterCount int

	CreatorID string
	// Creator 仅在 includes 指定时填充
	Creator *UserInfo

	LastActiveAt time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

// ComposeComicTitle 根据序号、作者和标题组合生成展示用的漫画标题
func (c *ComicInfo) ComposeComicTitle() string {
	return fmt.Sprintf("【%d】[%s] %s", c.Index, c.Author, c.Title)
}

// ComicCreation 是创建漫画时的载荷
type ComicCreation struct {
	// ID 由 domain service 生成，外部不提供
	ID string

	// WorksetID 是必须的，因为一个漫画必须属于一个作品集
	WorksetID string
	// Index 是该漫画在作品集内的序号
	Index int

	Title       string
	Author      string
	Description string
	// CreatorID 是必须的，记录创建者
	CreatorID string
}

// ComicUpdate 是漫画更新信息的 view，是 PUT 语义的载荷
type ComicUpdate struct {
	ID string

	Title       string
	Author      string
	Description string
}

// ComicQueryOpt 指定漫画查询的可选筛选条件
// 所有字段均为可空，nil 表示不参与筛选
type ComicQueryOpt struct {
	// ID 按漫画 ID 筛选
	ID *string

	// WorksetID 按所属作品集 ID 筛选
	// 这是必选的，因为漫画必须属于一个作品集
	WorksetID string
	// FuzzyTitle 按标题模糊匹配筛选
	// 需要注意的是，其匹配的是 composed title，即包含序号和作者的标题
	FuzzyTitle *string

	// 进度筛选，根据 pinned 字段值筛选
	UploadStatus    *WorkflowPhase
	TranslateStatus *WorkflowPhase
	ProofreadStatus *WorkflowPhase
	TypesetStatus   *WorkflowPhase
	ReviewStatus    *WorkflowPhase
	PublishStatus   *WorkflowPhase

	// Includes 指定查询时要包含的反向单射数据
	Includes []ComicInclude

	Pagination
}
