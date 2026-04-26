package val

type UserStats struct {
	UserId string `json:"user_id"`

	TotalAssignmentCnt    int `json:"total_assignment_count"`
	ActiveAssignmentCnt   int `json:"active_assignment_count"`
	FinishedAssignmentCnt int `json:"finished_assignment_count"`
}
