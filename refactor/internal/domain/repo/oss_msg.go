package repo_iface

import (
	"time"

	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
)

type OssMsgRepo interface {
	SavePendingCre(msg *aggr.OssCreMsg) RepoErr
	SavePendingDel(msg *aggr.OssDelMsg) RepoErr

	// `MarkCompleted` simply marks a message as completed by its ID,
	// without checking its current status or associated resource.
	MarkCompleted(id string) RepoErr
	// `MarkCompletedByRes` marks a pending create message as completed based on the resource type and ID.
	// A resId may be linked to more than one message, but at most one of them can be in pending state.
	// This method will mark that one pending message as completed.
	MarkCompletedByRes(ty enum.OssResTyp, resId string) RepoErr

	// `ClaimPending` claims one pending message for the specified operation.
	// It returns nil when no message is available.
	ClaimPending(op enum.OssOp) (*aggr.OssMsg, RepoErr)

	// `MarkPending` resets one processing message to pending and records retry metadata.
	MarkPending(id string, errMsg string, nextVisibleAt time.Time) RepoErr

	// `ResetStuck` resets messages that are in a processing state
	// and were created before the specified time.
	ResetStuck(bef time.Time) RepoErr

	// `CleanCompleted` cleans up completed messages that were created before the specified time.
	CleanCompleted(bef time.Time) RepoErr
}
