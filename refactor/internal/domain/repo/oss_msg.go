package repo_iface

import (
	"time"

	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
)

type OssMsgRepo interface {
	SavePendingCre(msg *aggr.OssCreMsg) RepoErr
	SavePendingDel(msg *aggr.OssDelMsg) RepoErr

	// `MarkCmpl` simply marks a message as completed by its ID,
	// without checking its current status or associated resource.
	MarkCmpl(id string) RepoErr
	// `MarkCmplByRes` marks a pending create message as completed based on the resource type and ID.
	// A resId may be linked to more than one message, but at most one of them can be in pending state.
	// This method will mark that one pending message as completed.
	MarkCmplByRes(ty enum.OssResTyp, resId string) RepoErr

	ClaimPending()

	MarkPending(id string) RepoErr

	// `ResetStuck` resets messages that are in a processing state
	// and were created before the specified time.
	ResetStuck(bef time.Time) RepoErr

	// `CleanCmpl` cleans up completed messages that were created before the specified time.
	CleanCmpl(bef time.Time) RepoErr
}
