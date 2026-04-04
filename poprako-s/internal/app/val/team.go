package val

// TeamInfo 表示用于在应用层和接口层传递的汉化组信息 VO
type TeamInfo struct {
	// ID 是汉化组的唯一标识
	ID string `json:"id"`

	// Name 是汉化组的名称
	Name string `json:"name"`
	// Description 是汉化组的描述
	Description string `json:"description"`

	// AvatarURL 是汉化组头像的可访问地址
	AvatarURL string `json:"avatar_url"`
	// IsAvatarUploaded 表示汉化组是否已上传头像
	IsAvatarUploaded bool `json:"is_avatar_uploaded"`

	// CreatedAt 是记录创建时间的 Unix 毫秒时间戳
	CreatedAt int64 `json:"created_at"`
	// UpdatedAt 是记录最近一次更新时间的 Unix 毫秒时间戳
	UpdatedAt int64 `json:"updated_at"`
}

// CreateTeamArgs 表示创建汉化组请求的参数
type CreateTeamArgs struct {
	// Name 是汉化组名称
	Name string `json:"name" validate:"required"`
	// Description 是汉化组描述
	Description string `json:"description"`
}

// CreateTeamRes 表示创建汉化组成功后的响应数据
type CreateTeamRes struct {
	// ID 是新创建汉化组的标识
	ID string `json:"id"`
}

// UpdateTeamArgs 表示更新汉化组信息的参数
type UpdateTeamArgs struct {
	// ID 是要更新的汉化组标识
	ID string `json:"id" validate:"required"`
	// Name 是更新后的汉化组名称
	Name string `json:"name" validate:"required"`
	// Description 是更新后的汉化组描述
	Description string `json:"description"`
}

// ReserveTeamAvatarRes 表示预留汉化组头像上传接口的响应数据
type ReserveTeamAvatarRes struct {
	// PutURL 是用于上传头像的预签名 URL
	PutURL string `json:"put_url"`
}

// ListTeamArgs 表示列出所有汉化组请求的参数
type ListTeamArgs struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// ListMyTeamArgs 表示列出当前用户所属汉化组请求的参数
type ListMyTeamArgs struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}
