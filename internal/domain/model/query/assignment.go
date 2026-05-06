package query

import "poprako-s/internal/domain/model/enum"

// `ListAssignmentOpt` defines list filters for assignment.
type ListAssignmentOpt struct {
	ChapterId *string
	UserId    *string

	// `Includes` controls relation preloads supported by the repo.
	Includes []enum.AssignmentIncl

	Pagi PagiOpt
}
