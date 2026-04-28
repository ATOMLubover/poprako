package entity

import (
	"time"

	"poprako-s/internal/domain/model/aggr"
)

// `SYS_MAIL_TABLE` is the table name for system mail records.
const SYS_MAIL_TABLE = "t_system_mail"

// `SysMailRow` maps one system mail record for read queries.
type SysMailRow struct {
	Id string `gorm:"column:id;primaryKey"`

	RcvId   string `gorm:"column:receiver_id"`
	Title   string `gorm:"column:title"`
	Content string `gorm:"column:content"`

	Read bool `gorm:"column:read"`

	CreatedAt time.Time `gorm:"column:created_at"`
}

// `TableName` returns the table name of `SysMailRow`.
func (*SysMailRow) TableName() string {
	return SYS_MAIL_TABLE
}

// `ToSysMailAggr` converts a row into `SysMail` aggregate.
func (r *SysMailRow) ToSysMailAggr() *aggr.SysMail {
	if r == nil {
		return nil
	}

	return &aggr.SysMail{
		Id:        r.Id,
		RcvId:     r.RcvId,
		Read:      r.Read,
		Title:     r.Title,
		Content:   r.Content,
		CreatedAt: r.CreatedAt,
	}
}

// `SysMailCreRow` maps one system mail record for create queries.
type SysMailCreRow struct {
	Id string `gorm:"column:id;primaryKey"`

	RcvId   string `gorm:"column:receiver_id"`
	Title   string `gorm:"column:title"`
	Content string `gorm:"column:content"`

	CreatedAt time.Time `gorm:"column:created_at"`
}

// `TableName` returns the table name of `SysMailCreRow`.
func (*SysMailCreRow) TableName() string {
	return SYS_MAIL_TABLE
}

// `NewSysMailCreRowFromAggr` builds a create row from `SysMailCre`.
func NewSysMailCreRowFromAggr(cre *aggr.SysMailCre) *SysMailCreRow {
	return &SysMailCreRow{
		Id:        cre.Id,
		RcvId:     cre.RcvId,
		Title:     cre.Title,
		Content:   cre.Content,
		CreatedAt: time.Now(),
	}
}

// `sysMailMarkReadUpdRow` maps update columns for mark-read operation.
type sysMailMarkReadUpdRow struct {
	Read bool `gorm:"column:read"`
}

// `NewSysMailMarkReadUpdRow` builds update row for mark-read operation.
func NewSysMailMarkReadUpdRow() *sysMailMarkReadUpdRow {
	return &sysMailMarkReadUpdRow{Read: true}
}
