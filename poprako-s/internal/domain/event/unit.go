package event

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

	// 统计变化量，由 domain service 在事务内计算完成后封装进事件
	TotalDelta      int
	TranslatedDelta int
	ProofreadDelta  int
}

func (e *UnitSaveEvent) EventType() EventType {
	return EventTypeUnitSave
}

func (e *UnitSaveEvent) Payload() any {
	return e
}
