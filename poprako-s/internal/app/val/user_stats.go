package val

type UserStatsInfo struct {
	UserID string `json:"user_id"`

	TotalAssignmentCount    int `json:"total_assignment_count"`
	ActiveAssignmentCount   int `json:"active_assignment_count"`
	FinishedAssignmentCount int `json:"finished_assignment_count"`
}
