package domain

import (
	"gamification-app/pkg/constants"
	"time"
)

type UserGameHistory struct {
	ID           uint
	UserID       uint
	PointsEarned int
	CreatedAt    time.Time
}

type FindAllGameHistoryRequest struct {
	Nickname  string
	Page      int
	Limit     int
	SortOrder constants.SortOrder
}

type FindAllGameHistoryResult struct {
	Items []UserGameHistory
	Total int64
}
