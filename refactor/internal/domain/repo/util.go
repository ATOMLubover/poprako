package repo_iface

// `ErrClsf` classifies repository error
// into categories like not found, timeout, and unavailable.`
type ErrClsf interface {
	// `IsNotFound` returns true if the error indicates a not found condition.
	IsNotFound(err RepoErr) bool

	// `IsTimeout` returns true if the error indicates a timeout condition.
	IsTimeout(err RepoErr) bool

	// `IsUnavailable` returns true if the error indicates an unavailable condition.
	IsUnavailable(err RepoErr) bool
}
