package event_handler

import (
	"context"
	"fmt"

	"poprako-s/internal/domain/repo"
)

type txnRepoCtxKey string

const (
	assignmentRepoTxnKey txnRepoCtxKey = "assignment-repo-txn"
	chapterRepoTxnKey    txnRepoCtxKey = "chapter-repo-txn"
	comicRepoTxnKey      txnRepoCtxKey = "comic-repo-txn"
	pageRepoTxnKey       txnRepoCtxKey = "page-repo-txn"
	userRepoTxnKey       txnRepoCtxKey = "user-repo-txn"
	worksetRepoTxnKey    txnRepoCtxKey = "workset-repo-txn"
)

func withTxnRepo(cx context.Context, key txnRepoCtxKey, value any) context.Context {
	if cx == nil {
		cx = context.Background()
	}

	return context.WithValue(cx, key, value)
}

func getTxnRepo[T any](cx context.Context, key txnRepoCtxKey, label string) (T, error) {
	var zero T

	if cx == nil {
		return zero, fmt.Errorf("[%s] 事务上下文为空", label)
	}

	value, ok := cx.Value(key).(T)
	if !ok {
		return zero, fmt.Errorf("[%s] 无法从上下文中获取事务版 Repo", label)
	}

	return value, nil
}

func WithAssignmentRepoTxn(cx context.Context, assignmentRepo repo.AssignmentRepo) context.Context {
	return withTxnRepo(cx, assignmentRepoTxnKey, assignmentRepo)
}

func assignmentRepoFromCx(cx context.Context) (repo.AssignmentRepo, error) {
	return getTxnRepo[repo.AssignmentRepo](cx, assignmentRepoTxnKey, "event_handler.assignmentRepoFromCx")
}

func WithChapterRepoTxn(cx context.Context, chapterRepo repo.ChapterRepo) context.Context {
	return withTxnRepo(cx, chapterRepoTxnKey, chapterRepo)
}

func chapterRepoFromCx(cx context.Context) (repo.ChapterRepo, error) {
	return getTxnRepo[repo.ChapterRepo](cx, chapterRepoTxnKey, "event_handler.chapterRepoFromCx")
}

func WithComicRepoTxn(cx context.Context, comicRepo repo.ComicRepo) context.Context {
	return withTxnRepo(cx, comicRepoTxnKey, comicRepo)
}

func comicRepoFromCx(cx context.Context) (repo.ComicRepo, error) {
	return getTxnRepo[repo.ComicRepo](cx, comicRepoTxnKey, "event_handler.comicRepoFromCx")
}

func WithPageRepoTxn(cx context.Context, pageRepo repo.PageRepo) context.Context {
	return withTxnRepo(cx, pageRepoTxnKey, pageRepo)
}

func pageRepoFromCx(cx context.Context) (repo.PageRepo, error) {
	return getTxnRepo[repo.PageRepo](cx, pageRepoTxnKey, "event_handler.pageRepoFromCx")
}

func WithUserRepoTxn(cx context.Context, userRepo repo.UserRepo) context.Context {
	return withTxnRepo(cx, userRepoTxnKey, userRepo)
}

func userRepoFromCx(cx context.Context) (repo.UserRepo, error) {
	return getTxnRepo[repo.UserRepo](cx, userRepoTxnKey, "event_handler.userRepoFromCx")
}

func WithWorksetRepoTxn(cx context.Context, worksetRepo repo.WorksetRepo) context.Context {
	return withTxnRepo(cx, worksetRepoTxnKey, worksetRepo)
}

func worksetRepoFromCx(cx context.Context) (repo.WorksetRepo, error) {
	return getTxnRepo[repo.WorksetRepo](cx, worksetRepoTxnKey, "event_handler.worksetRepoFromCx")
}
