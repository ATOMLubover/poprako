package model

import "time"

// WorksetInfo 表示作品集（Workset）的元信息
type WorksetInfo struct {
	ID string

	TeamID string
	// Team 仅在 includes 指定时填充
	Team *TeamInfo

	// 在汉化组内部的序号
	Index int

	Name        string
	Description string
	ComicCount  int

	CreatedAt time.Time
	UpdatedAt time.Time
}

// WorksetCreation 是创建作品集时的载荷
type WorksetCreation struct {
	// ID 由 domain service 生成，外部不提供
	ID string

	// TeamID 是必须的，因为一个作品集必须属于一个汉化组
	TeamID string
	// Index 是该作品集在汉化组内的序号
	Index int

	Name        string
	Description string
}

// WorksetUpdate 是作品集更新信息的 view
type WorksetUpdate struct {
	ID string

	Name string
	// Description 是可空的，nil 表示不修改
	Description *string
}

// WorksetQueryOpt 指定作品集查询的可选筛选条件
// 所有字段均为可空，nil 表示不参与筛选
type WorksetQueryOpt struct {
	// ID 按作品集 ID 筛选
	ID *string
	// TeamID 按所属汉化组 ID 筛选
	TeamID *string
}
