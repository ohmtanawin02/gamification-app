package domain

import (
	"gamification-app/pkg/constants"
	"time"
)

type Users struct {
	ID uint
	Nickname string
	TotalPoints int
	CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt *time.Time
}

type LoginInput struct {
	Nickname string
}

type LoginResult struct {
	Token string
	User  *Users
}

type FindAllUsersRequest struct {
    Nickname    string
    Page      int
    Limit     int
    SortOrder constants.SortOrder
}

type FindAllUsersResult struct {
    Items []Users
    Total int64
}