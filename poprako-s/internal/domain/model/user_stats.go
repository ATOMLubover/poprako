package model

type UserStats struct {
	UserID string

	TotalAssignmentCount    int
	ActiveAssignmentCount   int
	FinishedAssignmentCount int
}

// UserStatsPatch 是用于更新用户统计信息的结构体
// 其中只记录变化量 delta
type UserStatsPatch struct {
	// 这里的 UserID 是必须的，因为我们需要知道要更新哪个用户的统计信息
	UserID string

	// 以下这些变量保持对负值的兼容性，以便于在需要减少统计计数时使用
	// 如果不需要修改，则对应的 delta 应该设置为 0（默认值）
	TotalAssignmentCountDelta    int
	ActiveAssignmentCountDelta   int
	FinishedAssignmentCountDelta int
}
