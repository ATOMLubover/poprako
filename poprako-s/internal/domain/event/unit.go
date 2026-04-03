package event

import "context"

const (
	EventTypeUnitSave EventType = "UnitSaveEvent"
)

// UnitSaveEvent 是一个事件，在一个 page 批保存 units 后触发
type UnitSaveEvent struct {
	PageID    string
	ChapterID string

	InsertCount int
	PatchCount  int
	DeleteCount int

	// 统计变化量，由 Save 在事务内计算完成后随事件传递给 Handler
	TotalDelta      int
	TranslatedDelta int
	ProofreadDelta  int

	// 包含事务的上下文，可以使用 FromCx 方法来创建
	// 带事务的 repo
	Cx context.Context
}

func (e *UnitSaveEvent) EventType() EventType {
	return EventTypeUnitSave
}

func (e *UnitSaveEvent) PubType() PubType {
	// UnitSaveEvent 必须在事务中被同步处理，否则将会导致
	// chapter 的统计数据出错
	return PubTypeSync
}

func (e *UnitSaveEvent) Payload() any {
	return e
}
