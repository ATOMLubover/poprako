package repo_infra

import (
	"context"
	"errors"
	"net"
	"strings"

	repo_iface "poprako-s/internal/domain/repo"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

const (
	pgErrDupKey      = "23505"
	pgErrQueryCancel = "57014"
	pgErrSerFail     = "40001"
	pgErrDeadlock    = "40P01"
)

var errConditionalUpdateFailed = errors.New("conditional update failed")

// `IsDupKey` reports whether `err` is a duplicate-key conflict.
func IsDupKey(err repo_iface.RepoErr) bool {
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}

	if pgErr := takePgErr(err); pgErr != nil {
		return pgErr.Code == pgErrDupKey
	}

	return false
}

// `IsNotFound` reports whether `err` means no record was found.
func IsNotFound(err repo_iface.RepoErr) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

// `IsConflict` reports whether `err` should be treated as a write conflict.
func IsConflict(err repo_iface.RepoErr) bool {
	if IsDupKey(err) {
		return true
	}

	if pgErr := takePgErr(err); pgErr != nil {
		return pgErr.Code == pgErrSerFail || pgErr.Code == pgErrDeadlock
	}

	return false
}

// `IsConditionalUpdateFailed` reports whether a guarded update precondition was not met.
func IsConditionalUpdateFailed(err repo_iface.RepoErr) bool {
	return errors.Is(err, errConditionalUpdateFailed)
}

// `IsCanceled` reports whether `err` is caused by explicit cancellation.
func IsCanceled(err repo_iface.RepoErr) bool {
	if errors.Is(err, context.Canceled) {
		return true
	}

	if pgErr := takePgErr(err); pgErr != nil {
		return pgErr.Code == pgErrQueryCancel
	}

	return false
}

// `IsTimeout` reports whether `err` is caused by timeout.
func IsTimeout(err repo_iface.RepoErr) bool {
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		return netErr.Timeout()
	}

	if pgErr := takePgErr(err); pgErr != nil {
		return pgErr.Code == pgErrQueryCancel
	}

	return false
}

// `IsUnavailable` reports whether `err` indicates temporary backend unavailability.
func IsUnavailable(err repo_iface.RepoErr) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, gorm.ErrInvalidDB) {
		return true
	}

	if pgErr := takePgErr(err); pgErr != nil {
		if strings.HasPrefix(pgErr.Code, "08") {
			return true
		}

		switch pgErr.Code {
		case "57P01", "57P02", "57P03":
			return true
		}
	}

	return false
}

// `takePgErr` tries to extract typed Postgres error.
func takePgErr(err error) *pgconn.PgError {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr
	}

	return nil
}
