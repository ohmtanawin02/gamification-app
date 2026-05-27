package domain

import (
	"errors"
	"time"
)

type UserGameHistory struct {
	ID           uint
	UserID       uint
	PointsEarned int
	CreatedAt    time.Time
}

type SpinInput struct {
	UserID       uint
	PointsEarned int
}

type SpinResult struct {
	TotalPoints int
}

var (
	ErrPointsExceeded = errors.New("total points cannot exceed 10000")
	ErrInvalidPoints  = errors.New("invalid points value")
)
