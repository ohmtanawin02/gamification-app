package dto

import (
	"time"

	"gamification-app/internal/users/domain"
	"gamification-app/pkg/common"
)

type UserResponse struct {
	ID          uint      `json:"id"`
	Nickname    string    `json:"nickname"`
	TotalPoints int       `json:"total_points"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type UserListResponse struct {
	Items []UserResponse        `json:"items"`
	Meta  common.PaginationMeta `json:"meta"`
}

type RewardItem struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	CheckPoint int    `json:"checkpoint"`
	Claimed    bool   `json:"claimed"`
}

type UserMeResponse struct {
	Nickname    string       `json:"nickname"`
	TotalPoints int          `json:"total_points"`
	Rewards     []RewardItem `json:"rewards"`
}

func ToUserResponse(u domain.Users) UserResponse {
	return UserResponse{
		ID:          u.ID,
		Nickname:    u.Nickname,
		TotalPoints: u.TotalPoints,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}

func ToLoginResponse(result domain.LoginResult) LoginResponse {
	return LoginResponse{
		Token: result.Token,
		User:  ToUserResponse(*result.User),
	}
}

func ToUserListResponse(result domain.FindAllUsersResult, page, limit int) UserListResponse {
	items := make([]UserResponse, 0, len(result.Items))
	for _, u := range result.Items {
		items = append(items, ToUserResponse(u))
	}
	totalPages := int64(0)
	if limit > 0 {
		totalPages = (result.Total + int64(limit) - 1) / int64(limit)
	}
	return UserListResponse{
		Items: items,
		Meta:  common.PaginationMeta{Page: page, Limit: limit, Total: result.Total, TotalPages: totalPages},
	}
}
