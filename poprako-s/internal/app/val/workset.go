package val

// WorksetInfo 表示用于在应用层和接口层传递的作品集信息 VO
type WorksetInfo struct {
	// ID 是作品集的唯一标识
	ID string `json:"id"`

	// TeamID 是所属汉化组 ID
	TeamID string `json:"team_id"`
	// Team 是可选的汉化组信息（仅在 includes 时填充）
	Team *TeamInfo `json:"team,omitempty"`

	// Index 是作品集在汉化组内的序号
	Index int `json:"index"`

	// Name 是作品集名称
	Name string `json:"name"`
	// Description 是作品集描述
	Description string `json:"description"`
	// ComicCount 是作品集中漫画数量
	ComicCount int `json:"comic_count"`

	// CreatedAt 是记录创建时间的 Unix 毫秒时间戳
	CreatedAt int64 `json:"created_at"`
	// UpdatedAt 是记录最近一次更新时间的 Unix 毫秒时间戳
	UpdatedAt int64 `json:"updated_at"`
}

// ListWorksetArgs 表示列出作品集列表请求的参数
type ListWorksetArgs struct {
	// TeamID 是目标汉化组 ID
	TeamID string `json:"team_id" validate:"required"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
}

// CreateWorksetArgs 表示创建作品集请求的参数
type CreateWorksetArgs struct {
	// TeamID 是目标汉化组 ID
	TeamID string `json:"team_id" validate:"required"`
	// Name 是作品集名称
	Name string `json:"name" validate:"required"`
	// Description 是作品集描述（可选）
	Description *string `json:"description"`
}

// CreateWorksetRes 表示创建作品集成功后的响应数据
type CreateWorksetRes struct {
	// ID 是新创建作品集的标识
	ID string `json:"id"`
}

// UpdateWorksetArgs 表示更新作品集请求的参数
type UpdateWorksetArgs struct {
	// ID 是要更新的作品集标识
	ID string `json:"id" validate:"required"`
	// Name 是更新后的名称
	Name string `json:"name" validate:"required"`
	// Description 是更新后的描述（可选）
	Description *string `json:"description"`
}
