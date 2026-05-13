package repo_iface

import (
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/query"
)

// `SysMailRepo` defines persistence contract for system mail messages.
type SysMailRepo interface {
	// `Send` creates a new system mail record for the given mail creation aggregate.
	Send(cre *aggr.SysMailCre) RepoErr

	// `SendBatch` creates multiple system mail records in one batch.
	SendBatch(cres []*aggr.SysMailCre) RepoErr

	// `ListUnreadByRcvId` returns unread system mails by receiver and pagination options.
	ListUnreadByRcvId(rcvId string, pagi query.PagiOpt) ([]aggr.SysMail, RepoErr)

	// `MarkReadByRcvId` marks one system mail as read by id and receiver id.
	MarkReadByRcvId(id string, rcvId string) RepoErr
}
