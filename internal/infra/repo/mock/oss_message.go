package mock_repo

import (
	"context"
	"time"

	"poprako-s/internal/domain/model"
	iface "poprako-s/internal/domain/repo"
)

// OSSMessageRepo 是 repo.OSSMessageRepo 的内存 mock，用于单元测试。
type OSSMessageRepo struct {
	UpsertErr   error
	InsertErr   error
	CompleteErr error
	ResetErr    error
	DeleteErr   error

	Messages []*model.OSSMessage
}

func NewMockOSSMessageRepo() *OSSMessageRepo {
	return &OSSMessageRepo{}
}

func WithMockOSSMessageRepo(cx context.Context, repo *OSSMessageRepo) context.Context {
	return withMockRepo(cx, ossMessageRepoCtxKey, repo)
}

func NewMockOSSMessageRepoFromCx(cx context.Context) (iface.OSSMessageRepo, error) {
	return getMockRepoFromCx[*OSSMessageRepo](cx, ossMessageRepoCtxKey, "oss_message")
}

func (r *OSSMessageRepo) FromTxnCx(cx context.Context) (iface.OSSMessageRepo, error) {
	return NewMockOSSMessageRepoFromCx(cx)
}

func (r *OSSMessageRepo) UpsertCreatePending(msg *model.OSSMessage) error {
	if r.UpsertErr != nil {
		return r.UpsertErr
	}

	r.Messages = append(r.Messages, msg)

	return nil
}

func (r *OSSMessageRepo) InsertDeletePending(msg *model.OSSMessage) error {
	if r.InsertErr != nil {
		return r.InsertErr
	}

	r.Messages = append(r.Messages, msg)

	return nil
}

func (r *OSSMessageRepo) MarkCompleted(id string) error {
	return nil
}

func (r *OSSMessageRepo) CompleteCreatePendingByResource(resourceType model.OSSResourceType, resourceID string) error {
	return r.CompleteErr
}

func (r *OSSMessageRepo) ClaimPending(operation model.OSSOperation) (*model.OSSMessage, error) {
	return nil, nil
}

func (r *OSSMessageRepo) ResetToRetry(id string, lastError string, nextVisibleAt interface{}) error {
	return nil
}

func (r *OSSMessageRepo) ResetStuckProcessing(before time.Time) error {
	return r.ResetErr
}

func (r *OSSMessageRepo) DeleteCompletedBefore(before time.Time) error {
	return r.DeleteErr
}
