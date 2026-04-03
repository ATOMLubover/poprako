package mock_repo

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"
)

var errNotFound = errors.New("mock: not found")

type mockRepoCtxKey string

const (
	assignmentRepoCtxKey mockRepoCtxKey = "assignment-repo"
	chapterRepoCtxKey    mockRepoCtxKey = "chapter-repo"
	comicRepoCtxKey      mockRepoCtxKey = "comic-repo"
	invitationRepoCtxKey mockRepoCtxKey = "invitation-repo"
	memberRepoCtxKey     mockRepoCtxKey = "member-repo"
	pageRepoCtxKey       mockRepoCtxKey = "page-repo"
	teamRepoCtxKey       mockRepoCtxKey = "team-repo"
	unitRepoCtxKey       mockRepoCtxKey = "unit-repo"
	userRepoCtxKey       mockRepoCtxKey = "user-repo"
	worksetRepoCtxKey    mockRepoCtxKey = "workset-repo"
)

func withMockRepo(cx context.Context, key mockRepoCtxKey, value any) context.Context {
	if cx == nil {
		cx = context.Background()
	}

	return context.WithValue(cx, key, value)
}

func getMockRepoFromCx[T any](cx context.Context, key mockRepoCtxKey, label string) (T, error) {

	var zero T

	if cx == nil {
		return zero, fmt.Errorf("mock: %s repo missing from context", label)
	}

	value, ok := cx.Value(key).(T)
	if !ok {
		return zero, fmt.Errorf("mock: %s repo missing from context", label)
	}

	return value, nil
}

func sortedKeys[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func cloneTimePtr(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}

	cloned := *value
	return &cloned
}

func cloneStringPtr(value *string) *string {
	if value == nil {
		return nil
	}

	cloned := *value
	return &cloned
}

func cloneIntPtr(value *int) *int {
	if value == nil {
		return nil
	}

	cloned := *value
	return &cloned
}

func cloneBoolPtr(value *bool) *bool {
	if value == nil {
		return nil
	}

	cloned := *value
	return &cloned
}
