package repo_infra

import (
	repo_iface "poprako-s/internal/domain/repo"
)

// `defaultErrClassifier` adapts standalone classification functions
// into the `repo_iface.ErrClassifier` interface.
type defaultErrClassifier struct{}

func (defaultErrClassifier) IsNotFound(err repo_iface.RepoErr) bool    { return IsNotFound(err) }
func (defaultErrClassifier) IsTimeout(err repo_iface.RepoErr) bool     { return IsTimeout(err) }
func (defaultErrClassifier) IsUnavailable(err repo_iface.RepoErr) bool { return IsUnavailable(err) }

// `NewErrClassifier` returns a concrete `repo_iface.ErrClassifier`.
func NewErrClassifier() repo_iface.ErrClsf {
	return defaultErrClassifier{}
}
