package repo_iface

import (
	"time"

	"poprako-s/internal/domain/model/aggr"
)

type UserRepo interface {
	// `GetById` gets a user by his id.
	// A error will be returns if not found.
	GetById(id string) (*aggr.User, RepoErr)
	GetByQid(qid string) (*aggr.User, RepoErr)
	GetCredsByQid(qid string) (*aggr.UserCreds, RepoErr)

	// TODO: DO NOT add List API now.

	// `Register` creates a new user,
	Register(reg *aggr.UserReg) (*aggr.User, RepoErr)

	// Update(upd *aggr.UserUpd) RepoErr

	// // `Remove` executes a **soft** delete on given user id.
	// Remove(id string) RepoErr

	// `Refresh` update user last_active_at.
	Refresh(id string, activeAt time.Time) RepoErr

	// `PrefillAvatarKey` prefill avatar key, as client-side
	// uploading is not monitored by server directly.
	PrefillAvatarKey(id string, key string) RepoErr
	// `MarkAvatarUploaded` set avatar_uploaded true.
	// NOTE: PrefillAvatarKey **must** be called before it.
	MarkAvatarUploaded(id string) RepoErr
}
