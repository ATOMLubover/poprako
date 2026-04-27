package aggr

import (
	"time"

	"poprako-s/internal/domain/model/enum"
)

// `OssMsg` is the generic queue message used by background processors.
type OssMsg struct {
	Id string

	ResTyp enum.OssResTyp
	ResId  string
	Op     enum.OssOp
	Status enum.OssMsgStatus

	ObjKeys []string

	VisibleAt time.Time
	ExpireAt  time.Time
	ProcAt    *time.Time

	AttemptCnt int
	LastErr    string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type OssCreMsg struct {
	Id string

	// `ResTyp` and `ResId` together identify the logical resource associated with this message.
	ResTyp enum.OssResTyp
	ResId  string
	Status enum.OssMsgStatus

	// `ObjKeys` is a list of object keys in OSS that are involved in this message.
	// NOTE: it is typically only supported by PostgreSQL-based repositories.
	ObjKeys []string

	VisibleAt time.Time
	ExpireAt  time.Time
	// `ProcAt` is set when the message is claimed for processing,
	// and is used to detect stuck messages.
	ProcAt *time.Time

	AttemptCnt int
	LastErr    string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type OssDelMsg struct {
	Id string

	ResTyp enum.OssResTyp
	ResId  string
	Status enum.OssMsgStatus

	// `ObjKeys` is a list of object keys in OSS that are involved in this message.
	// NOTE: it is typically only supported by PostgreSQL-based repositories.
	ObjKeys []string

	VisibleAt time.Time
	ExpireAt  time.Time
	// `ProcAt` is set when the message is claimed for processing,
	// and is used to detect stuck messages.
	ProcAt *time.Time

	AttemptCnt int
	LastErr    string

	CreatedAt time.Time
	UpdatedAt time.Time
}
