package model

import "time"

// OSSResourceType 标识 OSS 资源的业务类型
type OSSResourceType string

const (
	// OSSResourceUserAvatar 用户头像
	OSSResourceUserAvatar OSSResourceType = "user_avatar"
	// OSSResourceTeamAvatar 汉化组头像
	OSSResourceTeamAvatar OSSResourceType = "team_avatar"
	// OSSResourceComicCover 漫画封面
	OSSResourceComicCover OSSResourceType = "comic_cover"
	// OSSResourcePageImage 章节页面图片
	OSSResourcePageImage OSSResourceType = "page_image"
)

// OSSOperation 标识消息的操作类型
type OSSOperation string

const (
	// OSSOperationCreatePending 上传预留待确认（超时则清理远程对象）
	OSSOperationCreatePending OSSOperation = "create_pending"
	// OSSOperationDeletePending 本地已删除，远程对象待清理
	OSSOperationDeletePending OSSOperation = "delete_pending"
)

// OSSMessageStatus 标识消息的消费状态
type OSSMessageStatus string

const (
	// OSSMessageStatusPending 等待消费
	OSSMessageStatusPending OSSMessageStatus = "pending"
	// OSSMessageStatusProcessing 正在消费
	OSSMessageStatusProcessing OSSMessageStatus = "processing"
	// OSSMessageStatusCompleted 已完成
	OSSMessageStatusCompleted OSSMessageStatus = "completed"
)

// OSSMessage 本地 OSS 消息表的领域模型
type OSSMessage struct {
	ID           string
	ResourceType OSSResourceType
	ResourceID   string
	Operation    OSSOperation
	Status       OSSMessageStatus

	// ObjectKey 用于单对象操作场景
	ObjectKey string
	// PayloadJSON 用于多对象批量操作场景，存储 object keys 列表（JSON 数组字符串）
	PayloadJSON string

	// VisibleAt 消息何时可被 worker 消费，支持延迟重试
	VisibleAt time.Time
	// ExpireAt 主要用于上传预留超时清理（create_pending 场景）
	ExpireAt *time.Time

	// ProcessingAt 用于排查卡死消息
	ProcessingAt *time.Time

	AttemptCount int
	LastError    string

	CreatedAt time.Time
	UpdatedAt time.Time
}
