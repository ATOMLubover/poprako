package val

// TeamInfo 表示用于在应用层和接口层传递的团队信息 VO
type TeamInfo struct {
	// ID 是团队的唯一标识
	ID string `json:"id"`

	// Name 是团队的名称
	Name string `json:"name"`
	// Description 是团队的描述
	Description string `json:"description"`

	// AvatarURL 是团队头像的可访问地址
	AvatarURL string `json:"avatar_url"`
	// IsAvatarUploaded 表示团队是否已上传头像
	IsAvatarUploaded bool `json:"is_avatar_uploaded"`

	// CreatedAt 是记录创建时间的 Unix 毫秒时间戳
	CreatedAt int64 `json:"created_at"`
	// UpdatedAt 是记录最近一次更新时间的 Unix 毫秒时间戳
	UpdatedAt int64 `json:"updated_at"`
}

// CreateTeamArgs 表示创建团队请求的参数
type CreateTeamArgs struct {
	// Name 是团队名称
	Name string `json:"name" validate:"required"`
	// Description 是团队描述
	Description string `json:"description"`
}

// CreateTeamRes 表示创建团队成功后的响应数据
type CreateTeamRes struct {
	// ID 是新创建团队的标识
	ID string `json:"id"`
}

// UpdateTeamArgs 表示更新团队信息的参数
type UpdateTeamArgs struct {
	// ID 是要更新的团队标识
	ID string `json:"id" validate:"required"`
	// Name 是更新后的团队名称
	Name string `json:"name" validate:"required"`
	// Description 是更新后的团队描述
	Description string `json:"description"`
}

// ReserveTeamAvatarRes 表示预留团队头像上传接口的响应数据
type ReserveTeamAvatarRes struct {
	// PutURL 是用于上传头像的预签名 URL
	PutURL string `json:"put_url"`
}

// ListTeamArgs 表示列出所有团队请求的参数
type ListTeamArgs struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// ListMyTeamArgs 表示列出当前用户所属团队请求的参数
type ListMyTeamArgs struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}
