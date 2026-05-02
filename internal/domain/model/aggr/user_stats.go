package aggr

// `UserStats` holds aggregated assignment statistics for a single user
type UserStats struct {
	// `UserId` is the identifier of the user these stats belong to
	UserId string

	// `TotalAssignmentCnt` is the total number of assignments ever created for this user
	TotalAssignmentCnt int
	// `ActiveAssignmentCnt` is the number of currently active assignments
	ActiveAssignmentCnt int
	// `FinishedAssignmentCnt` is the number of completed assignments
	FinishedAssignmentCnt int
}

// `UserStatsPatch` contains deltas to
// optimistically update user stats.
type UserStatsPatch struct {
	UserId string

	TotalAssignmentDelta    int
	ActiveAssignmentDelta   int
	FinishedAssignmentDelta int
}
