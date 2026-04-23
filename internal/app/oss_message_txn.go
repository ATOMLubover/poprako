package app

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"poprako-s/internal/domain/model"
	"poprako-s/internal/domain/repo"
)

const ossCreatePendingExpiry = 24 * time.Hour

func genOSSMessageID() string {
	return fmt.Sprintf("ossmsg-%d", time.Now().UnixNano())
}

func upsertCreatePendingMessage(
	msgRepo repo.OSSMessageRepo,
	txCx context.Context,
	resourceType model.OSSResourceType,
	resourceID string,
	objectKey string,
) error {
	msgRepoTxn, err := msgRepo.FromTxnCx(txCx)
	if err != nil {
		return err
	}

	expireAt := time.Now().Add(ossCreatePendingExpiry)

	return msgRepoTxn.UpsertCreatePending(&model.OSSMessage{
		ID:           genOSSMessageID(),
		ResourceType: resourceType,
		ResourceID:   resourceID,
		ObjectKey:    objectKey,
		VisibleAt:    time.Now(),
		ExpireAt:     &expireAt,
	})
}

func completeCreatePendingMessage(
	msgRepo repo.OSSMessageRepo,
	txCx context.Context,
	resourceType model.OSSResourceType,
	resourceID string,
) error {
	msgRepoTxn, err := msgRepo.FromTxnCx(txCx)
	if err != nil {
		return err
	}

	return msgRepoTxn.CompleteCreatePendingByResource(resourceType, resourceID)
}

func enqueueDeleteSingleMessage(
	msgRepo repo.OSSMessageRepo,
	txCx context.Context,
	resourceType model.OSSResourceType,
	resourceID string,
	objectKey string,
) error {
	if objectKey == "" {
		return nil
	}

	msgRepoTxn, err := msgRepo.FromTxnCx(txCx)
	if err != nil {
		return err
	}

	return msgRepoTxn.InsertDeletePending(&model.OSSMessage{
		ID:           genOSSMessageID(),
		ResourceType: resourceType,
		ResourceID:   resourceID,
		ObjectKey:    objectKey,
		VisibleAt:    time.Now(),
	})
}

func enqueueDeleteBatchMessage(
	msgRepo repo.OSSMessageRepo,
	txCx context.Context,
	resourceType model.OSSResourceType,
	resourceID string,
	objectKeys []string,
) error {
	keys := make([]string, 0, len(objectKeys))
	for _, key := range objectKeys {
		if key != "" {
			keys = append(keys, key)
		}
	}

	if len(keys) == 0 {
		return nil
	}

	payloadJSON, err := json.Marshal(map[string][]string{"object_keys": keys})
	if err != nil {
		return err
	}

	msgRepoTxn, err := msgRepo.FromTxnCx(txCx)
	if err != nil {
		return err
	}

	return msgRepoTxn.InsertDeletePending(&model.OSSMessage{
		ID:           genOSSMessageID(),
		ResourceType: resourceType,
		ResourceID:   resourceID,
		PayloadJSON:  string(payloadJSON),
		VisibleAt:    time.Now(),
	})
}
