package model

import "time"

// PageInfo 表示一个漫画章节中的页面信息
type PageInfo struct {
	ID string

	ChapterID string
	// 页面在章节中的顺序（0-based）
	Index int

	OSSKey     string
	IsUploaded bool

	CreatorID string
	// Creator 仅在 includes 指定时填充
	Creator *UserInfo

	TotalUnitCount      int
	TranslatedUnitCount int
	ProofreadUnitCount  int

	CreatedAt time.Time
	UpdatedAt time.Time
}

// PageCreation 是创建页面时的载荷
type PageCreation struct {
	// ID 由 domain service 生成，外部不提供
	ID string

	ChapterID string
	// Index 是 0-based 的页面序号
	Index int

	OSSKey    string
	CreatorID string
}

// PageUpdate 是页面更新信息的 view
type PageUpdate struct {
	ID string

	Index      int
	OSSKey     string
	IsUploaded bool

	// 以下字段仅用于 unit save 时同步更新页面的统计数据
	TotalUnitCount      int
	TranslatedUnitCount int
	ProofreadUnitCount  int
}

// PageStats 聚合一个页面的统计数据
type PageStats struct {
	PageID string

	TotalUnitCount      int
	TranslatedUnitCount int
	ProofreadUnitCount  int
}

// PageQueryOpt 指定页面查询的可选筛选条件
// 所有字段均为可空，nil 表示不参与筛选
type PageQueryOpt struct {
	// ID 按页面 ID 筛选
	ID *string
	// ChapterID 按所属章节 ID 筛选
	ChapterID *string
}
