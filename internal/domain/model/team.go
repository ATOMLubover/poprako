package model

import "time"

// TeamInfo 表示汉化组（如汉化组）的元信息
type TeamInfo struct {
	ID string

	Name             string
	Description      string
	AvatarOSSKey     string
	IsAvatarUploaded bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

// TeamCreation 是创建汉化组时的载荷
type TeamCreation struct {
	// ID 由 domain service 生成，外部不提供
	ID string

	Name        string
	Description string
}

// TeamUpdate 是汉化组更新信息的 view，是 PUT 语义的载荷
type TeamUpdate struct {
	ID string

	Name        string
	Description string
}

// TeamQueryOpt 指定汉化组查询的可选筛选条件
// 所有字段均为可空，nil 表示不参与筛选
type TeamQueryOpt struct {
	// ID 按汉化组 ID 筛选
	ID *string
}
