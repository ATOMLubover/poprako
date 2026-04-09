package app

import (
	"errors"
	"sync"
	"testing"
	"time"
)

func TestOSSDeleteExecutorDeleteOneRetriesThreeTimes(t *testing.T) {
	ossClient := newMockOSSClient()
	ossClient.SetDeleteErr(errors.New("boom"))

	slept := make([]time.Duration, 0)
	executor := newOSSDeleteExecutor(ossClient)
	executor.sleepFn = func(d time.Duration) {
		slept = append(slept, d)
	}

	err := executor.deleteOne("k-1")
	if err == nil {
		t.Fatal("expected deleteOne to fail")
	}

	deleted := ossClient.Deleted()
	if len(deleted) != 3 {
		t.Fatalf("expected 3 delete attempts, got %#v", deleted)
	}

	if len(slept) != 2 || slept[0] != 100*time.Millisecond || slept[1] != 200*time.Millisecond {
		t.Fatalf("unexpected retry sleep durations: %#v", slept)
	}
}

func TestOSSDeleteExecutorDeleteManyRunsConcurrently(t *testing.T) {
	ossClient := newMockOSSClient()

	var mu sync.Mutex
	current := 0
	maxParallel := 0

	entered := make(chan struct{}, 3)
	block := make(chan struct{})

	ossClient.SetDeleteFunc(func(_ string) error {
		mu.Lock()
		current++
		if current > maxParallel {
			maxParallel = current
		}
		mu.Unlock()

		entered <- struct{}{}
		<-block

		mu.Lock()
		current--
		mu.Unlock()

		return nil
	})

	executor := newOSSDeleteExecutor(ossClient)
	executor.sleepFn = func(time.Duration) {}

	done := make(chan error, 1)
	go func() {
		done <- executor.deleteMany([]string{"k-1", "k-2", "k-3"})
	}()

	for i := 0; i < 3; i++ {
		select {
		case <-entered:
		case <-time.After(1 * time.Second):
			close(block)
			t.Fatal("expected concurrent delete workers to start")
		}
	}

	close(block)

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("unexpected deleteMany error: %v", err)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("deleteMany did not finish in time")
	}

	if maxParallel < 2 {
		t.Fatalf("expected concurrent deletes, max parallel=%d", maxParallel)
	}
}
