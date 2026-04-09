package app

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"poprako-s/internal/domain/ext/oss"
)

const (
	ossDeleteRetryAttempts = 3
	ossDeleteBackoffStep   = 100 * time.Millisecond
	ossDeleteParallelism   = 8
)

type ossDeleteExecutor struct {
	client  oss.Client
	sleepFn func(time.Duration)
}

func newOSSDeleteExecutor(client oss.Client) *ossDeleteExecutor {
	return &ossDeleteExecutor{
		client:  client,
		sleepFn: time.Sleep,
	}
}

func (e *ossDeleteExecutor) deleteOne(objectKey string) error {
	if objectKey == "" {
		return nil
	}

	var lastErr error

	for attempt := 1; attempt <= ossDeleteRetryAttempts; attempt++ {
		err := e.client.Delete(objectKey)
		if err == nil {
			return nil
		}

		lastErr = err

		if attempt < ossDeleteRetryAttempts {
			e.sleepFn(time.Duration(attempt) * ossDeleteBackoffStep)
		}
	}

	return fmt.Errorf("delete oss object failed: %s: %w", objectKey, lastErr)
}

func (e *ossDeleteExecutor) deleteMany(objectKeys []string) error {
	keys := make([]string, 0, len(objectKeys))
	for _, key := range objectKeys {
		if key == "" {
			continue
		}

		keys = append(keys, key)
	}

	if len(keys) == 0 {
		return nil
	}

	sem := make(chan struct{}, ossDeleteParallelism)
	errCh := make(chan error, len(keys))
	wg := sync.WaitGroup{}

	for _, key := range keys {
		wg.Add(1)

		go func(objectKey string) {
			defer wg.Done()

			sem <- struct{}{}
			err := e.deleteOne(objectKey)
			<-sem

			if err != nil {
				errCh <- err
			}
		}(key)
	}

	wg.Wait()
	close(errCh)

	errs := make([]error, 0)
	for err := range errCh {
		errs = append(errs, err)
	}

	if len(errs) == 0 {
		return nil
	}

	return errors.Join(errs...)
}
