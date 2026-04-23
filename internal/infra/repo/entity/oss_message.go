package repo_entity

import (
	"time"

	"poprako-s/internal/domain/model"
)

const OSSMessageTable = "oss_message_table"

// OSSMessageRow 是 oss_message_table 的数据库行结构
type OSSMessageRow struct {
	ID           string     `gorm:"column:id"`
	ResourceType string     `gorm:"column:resource_type"`
	ResourceID   string     `gorm:"column:resource_id"`
	Operation    string     `gorm:"column:operation"`
	Status       string     `gorm:"column:status"`
	ObjectKey    string     `gorm:"column:object_key"`
	PayloadJSON  string     `gorm:"column:payload_json"`
	VisibleAt    time.Time  `gorm:"column:visible_at"`
	ExpireAt     *time.Time `gorm:"column:expire_at"`

	ProcessingAt *time.Time `gorm:"column:processing_at"`

	AttemptCount int    `gorm:"column:attempt_count"`
	LastError    string `gorm:"column:last_error"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// ToOSSMessage 将数据库行转换为领域模型
func ToOSSMessage(row OSSMessageRow) model.OSSMessage {
	return model.OSSMessage{
		ID:           row.ID,
		ResourceType: model.OSSResourceType(row.ResourceType),
		ResourceID:   row.ResourceID,
		Operation:    model.OSSOperation(row.Operation),
		Status:       model.OSSMessageStatus(row.Status),
		ObjectKey:    row.ObjectKey,
		PayloadJSON:  row.PayloadJSON,
		VisibleAt:    row.VisibleAt,
		ExpireAt:     row.ExpireAt,
		ProcessingAt: row.ProcessingAt,
		AttemptCount: row.AttemptCount,
		LastError:    row.LastError,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
}
