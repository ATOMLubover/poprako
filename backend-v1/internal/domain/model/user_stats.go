package model

import "time"

type UserStats struct {
	UserID string

	TotalAssignmentCount    int
	ActiveAssignmentCount   int
	FinishedAssignmentCount int

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUserStats(
	userID string,
	totalAssignmentCount int,
	activeAssignmentCount int,
	finishedAssignmentCount int,
	createdAt time.Time,
	updatedAt time.Time,
) UserStats {
	return UserStats{
		UserID:                  userID,
		TotalAssignmentCount:    totalAssignmentCount,
		ActiveAssignmentCount:   activeAssignmentCount,
		FinishedAssignmentCount: finishedAssignmentCount,
		CreatedAt:               createdAt,
		UpdatedAt:               updatedAt,
	}
}

type UserStatsCreation struct {
	UserID string

	TotalAssignmentCount    int
	ActiveAssignmentCount   int
	FinishedAssignmentCount int
}

func NewUserStatsCreation(
	userID string,
	totalAssignmentCount int,
	activeAssignmentCount int,
	finishedAssignmentCount int,
) UserStatsCreation {
	return UserStatsCreation{
		UserID:                  userID,
		TotalAssignmentCount:    totalAssignmentCount,
		ActiveAssignmentCount:   activeAssignmentCount,
		FinishedAssignmentCount: finishedAssignmentCount,
	}
}

type UserStatsDelta struct {
	UserID string

	TotalAssignmentCountDelta    int
	ActiveAssignmentCountDelta   int
	FinishedAssignmentCountDelta int
}

func NewUserStatsDelta(
	userID string,
	totalAssignmentCountDelta int,
	activeAssignmentCountDelta int,
	finishedAssignmentCountDelta int,
) UserStatsDelta {
	return UserStatsDelta{
		UserID:                       userID,
		TotalAssignmentCountDelta:    totalAssignmentCountDelta,
		ActiveAssignmentCountDelta:   activeAssignmentCountDelta,
		FinishedAssignmentCountDelta: finishedAssignmentCountDelta,
	}
}
