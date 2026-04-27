package repo_iface

import "poprako-s/internal/domain/model/aggr"

type SysMailRepo interface {
	// `Send` creates a new system mail record for the given mail creation aggregate.
	Send(cre *aggr.SysMailCre) RepoErr

	// `MarkRead` marks the system mail with the given ID as read.
	MarkRead(id string) RepoErr
}
