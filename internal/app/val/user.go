package val

// UserInfo 表示用于在应用层和接口层传递的用户信息 VO
type UserInfo struct {
	// ID 是用户的唯一标识（例如 UUID 或数据库主键）
	ID string `json:"id"`

	// Name 是用户的显示名称或昵称
	Name string `json:"name"`
	// QQ 是用户的 QQ 号，用于登录或联系方式
	QQ string `json:"qq"`

	// AvatarURL 是用户头像的可访问地址
	AvatarURL string `json:"avatar_url"`
	// IsAvatarUploaded 表示用户是否已上传头像
	IsAvatarUploaded bool `json:"is_avatar_uploaded"`

	// IsSuperAdmin 表示用户是否具有超级管理员权限
	IsSuperAdmin bool `json:"is_super_admin"`
	// LastLoginAt 是一个 Unix 毫秒时间戳，表示用户最后一次登录的时间
	LastLoginAt int64 `json:"last_login_at"`

	// CreatedAt 是记录创建时间的 Unix 毫秒时间戳
	CreatedAt int64 `json:"created_at"`
	// UpdatedAt 是记录最近一次更新时间的 Unix 毫秒时间戳
	UpdatedAt int64 `json:"updated_at"`
}

// LoginUserArgs 表示用户登录请求的参数
type LoginUserArgs struct {
	// QQ 是登录使用的 QQ 号
	QQ string `json:"qq" validate:"required"`
	// Pwd 是登录使用的密码
	Pwd string `json:"password" validate:"required"`
}

// LoginUserRes 表示登录成功后的响应数据
type LoginUserRes struct {
	// UserID 是已登录用户的标识
	UserID string `json:"user_id"`
	// AccessToken 是用于后续鉴权的访问令牌
	AccessToken string `json:"access_token"`
}

// RegUserArgs 表示用户注册请求的参数
type RegUserArgs struct {
	// QQ 是注册使用的 QQ 号
	QQ string `json:"qq" validate:"required"`
	// Pwd 是注册使用的密码
	Pwd string `json:"password" validate:"required"`
	// Name 是用户注册时设置的显示名称
	Name string `json:"name" validate:"required"`
	// InvCode 是用于注册的邀请代码
	InvCode string `json:"invitation_code" validate:"required"`
}

// RegUserRes 表示注册成功后的响应数据
type RegUserRes struct {
	// UserID 是新创建用户的标识
	UserID string `json:"user_id"`
	// AccessToken 是注册后返回的访问令牌
	AccessToken string `json:"access_token"`
}

// UpdateUserArgs 是用户更新信息的 VO，用于 PUT 语义的请求载荷
type UpdateUserArgs struct {
	// ID 是要更新的用户标识
	ID string `json:"id" validate:"required"`

	// Name 表示更新后的显示名称，PUT 语义下为必填
	Name string `json:"name"`
	// QQ 表示更新后的 QQ 号，PUT 语义下为必填
	QQ string `json:"qq"`
}

// ReserveUserAvatarArgs 表示预留头像上传接口的请求参数
type ReserveUserAvatarArgs struct {
	// FileName 是用户头像文件的原始名称，主要用于 OSS 存储时保留扩展名
	// 需要携带文件扩展名以便 OSS 正确识别文件类型，例如 "avatar.png"
	FileName string `json:"file_name" validate:"required"`
}

// ReserveUserAvatarRes 表示预留头像上传接口的响应数据
type ReserveUserAvatarRes struct {
	// PutURL 是用于上传头像的预签名 URL，客户端可以直接使用该 URL 上传头像文件
	PutURL string `json:"put_url"`
}
