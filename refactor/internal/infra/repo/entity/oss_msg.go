package entity

import (
	"time"

	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
)

// `OSS_MSG_TABLE` is the table name for oss messages
const OSS_MSG_TABLE = "t_oss_message"

// `OssMsgRow` maps one oss message record for read and update queries
type OssMsgRow struct {
	Id string `gorm:"column:id;primaryKey"`

	ResTyp string `gorm:"column:resource_type"`
	ResId  string `gorm:"column:resource_id"`
	Op     string `gorm:"column:operation"`
	Status string `gorm:"column:status"`

	ObjKeys []string `gorm:"column:object_keys;type:text[]"`

	VisibleAt time.Time  `gorm:"column:visible_at"`
	ExpireAt  time.Time  `gorm:"column:expire_at"`
	ProcAt    *time.Time `gorm:"column:processing_at"`

	AttemptCnt int    `gorm:"column:attempt_count"`
	LastErr    string `gorm:"column:last_error"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// `TableName` returns the table name of `OssMsgRow`
func (*OssMsgRow) TableName() string {
	return OSS_MSG_TABLE
}

// `ToOssMsgAggr` converts one row to the generic queue aggregate.
func (r *OssMsgRow) ToOssMsgAggr() *aggr.OssMsg {
	if r == nil {
		// Keep nil-safe conversion behavior for optional query results.
		return nil
	}

	return &aggr.OssMsg{
		Id:         r.Id,
		ResTyp:     enum.OssResTyp(r.ResTyp),
		ResId:      r.ResId,
		Op:         enum.OssOp(r.Op),
		Status:     enum.OssMsgStatus(r.Status),
		ObjKeys:    r.ObjKeys,
		VisibleAt:  r.VisibleAt,
		ExpireAt:   r.ExpireAt,
		ProcAt:     r.ProcAt,
		AttemptCnt: r.AttemptCnt,
		LastErr:    r.LastErr,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}
}

// `OssCreMsgCreRow` maps one create-operation message record for insert
type OssCreMsgCreRow struct {
	Id string `gorm:"column:id;primaryKey"`

	ResTyp string `gorm:"column:resource_type"`
	ResId  string `gorm:"column:resource_id"`
	Op     string `gorm:"column:operation"`
	Status string `gorm:"column:status"`

	ObjKeys []string `gorm:"column:object_keys;type:text[]"`

	VisibleAt time.Time  `gorm:"column:visible_at"`
	ExpireAt  time.Time  `gorm:"column:expire_at"`
	ProcAt    *time.Time `gorm:"column:processing_at"`

	AttemptCnt int    `gorm:"column:attempt_count"`
	LastErr    string `gorm:"column:last_error"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// `NewOssCreMsgCreRowFromAggr` builds a create row from `OssCreMsg`
func NewOssCreMsgCreRowFromAggr(msg *aggr.OssCreMsg) *OssCreMsgCreRow {
	now := time.Now()

	return &OssCreMsgCreRow{
		Id:         msg.Id,
		ResTyp:     string(msg.ResTyp),
		ResId:      msg.ResId,
		Op:         string(enum.OssOpCre),
		Status:     string(msg.Status),
		ObjKeys:    msg.ObjKeys,
		VisibleAt:  msg.VisibleAt,
		ExpireAt:   msg.ExpireAt,
		ProcAt:     msg.ProcAt,
		AttemptCnt: msg.AttemptCnt,
		LastErr:    msg.LastErr,
		CreatedAt:  msg.CreatedAt,
		UpdatedAt:  now,
	}
}

// `TableName` returns the table name of `OssCreMsgCreRow`
func (*OssCreMsgCreRow) TableName() string {
	return OSS_MSG_TABLE
}

// `OssDelMsgCreRow` maps one delete-operation message record for insert
type OssDelMsgCreRow struct {
	Id string `gorm:"column:id;primaryKey"`

	ResTyp string `gorm:"column:resource_type"`
	ResId  string `gorm:"column:resource_id"`
	Op     string `gorm:"column:operation"`
	Status string `gorm:"column:status"`

	ObjKeys []string `gorm:"column:object_keys;type:text[]"`

	VisibleAt time.Time  `gorm:"column:visible_at"`
	ExpireAt  time.Time  `gorm:"column:expire_at"`
	ProcAt    *time.Time `gorm:"column:processing_at"`

	AttemptCnt int    `gorm:"column:attempt_count"`
	LastErr    string `gorm:"column:last_error"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// `NewOssDelMsgCreRowFromAggr` builds a create row from `OssDelMsg`
func NewOssDelMsgCreRowFromAggr(msg *aggr.OssDelMsg) *OssDelMsgCreRow {
	now := time.Now()

	return &OssDelMsgCreRow{
		Id:         msg.Id,
		ResTyp:     string(msg.ResTyp),
		ResId:      msg.ResId,
		Op:         string(enum.OssOpDel),
		Status:     string(msg.Status),
		ObjKeys:    msg.ObjKeys,
		VisibleAt:  msg.VisibleAt,
		ExpireAt:   msg.ExpireAt,
		ProcAt:     msg.ProcAt,
		AttemptCnt: msg.AttemptCnt,
		LastErr:    msg.LastErr,
		CreatedAt:  msg.CreatedAt,
		UpdatedAt:  now,
	}
}

// `TableName` returns the table name of `OssDelMsgCreRow`
func (*OssDelMsgCreRow) TableName() string {
	return OSS_MSG_TABLE
}

// `ossMsgMarkCmplUpdRow` maps update columns for completion mark
type ossMsgMarkCmplUpdRow struct {
	Status    string    `gorm:"column:status"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// `NewOssMsgMarkCmplUpdRow` builds update row for completion mark
func NewOssMsgMarkCmplUpdRow() *ossMsgMarkCmplUpdRow {
	return &ossMsgMarkCmplUpdRow{
		Status:    string(enum.OssMsgStateCmpl),
		UpdatedAt: time.Now(),
	}
}

// `ossMsgMarkProcUpdRow` maps update columns for process claim
type ossMsgMarkProcUpdRow struct {
	Status    string     `gorm:"column:status"`
	ProcAt    *time.Time `gorm:"column:processing_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
}

// `NewOssMsgMarkProcUpdRow` builds update row for process claim
func NewOssMsgMarkProcUpdRow(now time.Time) *ossMsgMarkProcUpdRow {
	return &ossMsgMarkProcUpdRow{
		Status:    string(enum.OssMsgStateProc),
		ProcAt:    &now,
		UpdatedAt: now,
	}
}

// `ossMsgResetStuckUpdRow` maps update columns for stuck reset
type ossMsgResetStuckUpdRow struct {
	Status    string     `gorm:"column:status"`
	ProcAt    *time.Time `gorm:"column:processing_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
}

// `NewOssMsgResetStuckUpdRow` builds update row for stuck reset
func NewOssMsgResetStuckUpdRow() *ossMsgResetStuckUpdRow {
	return &ossMsgResetStuckUpdRow{
		Status:    string(enum.OssMsgStatePend),
		ProcAt:    nil,
		UpdatedAt: time.Now(),
	}
}

// `ossMsgMarkPendingUpdRow` maps update columns for pending reset,
// including status, processing_at, visible_at, last_error, and updated_at
type ossMsgMarkPendingUpdRow struct {
	Status    string     `gorm:"column:status"`
	ProcAt    *time.Time `gorm:"column:processing_at"`
	VisibleAt time.Time  `gorm:"column:visible_at"`
	LastErr   string     `gorm:"column:last_error"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
}

// `NewOssMsgMarkPendingUpdRow` builds an update row for pending reset with the given retry metadata
func NewOssMsgMarkPendingUpdRow(errMsg string, nextVisibleAt time.Time, now time.Time) *ossMsgMarkPendingUpdRow {
	return &ossMsgMarkPendingUpdRow{
		Status:    string(enum.OssMsgStatePend),
		ProcAt:    nil,
		VisibleAt: nextVisibleAt,
		LastErr:   errMsg,
		UpdatedAt: now,
	}
}
