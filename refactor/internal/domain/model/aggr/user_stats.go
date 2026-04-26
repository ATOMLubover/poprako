package aggr

type UserStats struct {
	UserId string

	TotalAssignmentCnt    int
	ActiveAssignmentCnt   int
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
