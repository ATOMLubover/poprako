package aggr

import (
	"time"
)

// `Assignment` represents one chapter assignment of one user.
type Assignment struct {
	TimedRoles

	Id string

	ChapterId string
	UserId    string

	CreatedAt time.Time
	UpdatedAt time.Time
}

// `AssignmentCre` is the create payload for assignment.
type AssignmentCre struct {
	Id string

	ChapterId string
	UserId    string

	TimedRoles TimedRoles
}

// `AssignmentPut` is the full overwrite payload of assignment roles.
type AssignmentPut struct {
	Id string

	TimedRoles TimedRoles
}
