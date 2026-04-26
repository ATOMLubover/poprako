package repo_infra

import (
	"errors"

	repo_iface "poprako-s/internal/domain/repo"

	"gorm.io/gorm"
)

func IsDupKey(err repo_iface.RepoErr) bool {
	return errors.Is(err, gorm.ErrDuplicatedKey)
}
