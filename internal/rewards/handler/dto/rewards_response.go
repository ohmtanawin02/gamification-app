package dto

import (
	"time"

	"gamification-app/internal/rewards/domain"
	"gamification-app/pkg/common"
)

type RewardsResponse struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	CheckPoint int       `json:"check_point"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type RewardsListResponse struct {
	Items []RewardsResponse     `json:"items"`
	Meta  common.PaginationMeta `json:"meta"`
}

func ToRewardsResponse(r domain.Rewards) RewardsResponse {
	return RewardsResponse{
		ID:         r.ID,
		Name:       r.Name,
		CheckPoint: r.CheckPoint,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}
}

type UserRewardResponse struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	CheckPoint int       `json:"check_point"`
	ClaimedAt  time.Time `json:"claimed_at"`
}

type UserRewardsListResponse struct {
	Items []UserRewardResponse  `json:"items"`
	Meta  common.PaginationMeta `json:"meta"`
}

func ToUserRewardResponse(r domain.UserReward) UserRewardResponse {
	return UserRewardResponse{
		ID:         r.ID,
		Name:       r.Name,
		CheckPoint: r.CheckPoint,
		ClaimedAt:  r.ClaimedAt,
	}
}

func ToUserRewardsListResponse(result domain.FindUserRewardsResult, page, limit int) UserRewardsListResponse {
	items := make([]UserRewardResponse, 0, len(result.Items))
	for _, r := range result.Items {
		items = append(items, ToUserRewardResponse(r))
	}
	totalPages := int64(0)
	if limit > 0 {
		totalPages = (result.Total + int64(limit) - 1) / int64(limit)
	}
	return UserRewardsListResponse{
		Items: items,
		Meta:  common.PaginationMeta{Page: page, Limit: limit, Total: result.Total, TotalPages: totalPages},
	}
}

func ToRewardsListResponse(result domain.FindAllRewardsResult, page, limit int) RewardsListResponse {
	items := make([]RewardsResponse, 0, len(result.Items))
	for _, r := range result.Items {
		items = append(items, ToRewardsResponse(r))
	}
	totalPages := int64(0)
	if limit > 0 {
		totalPages = (result.Total + int64(limit) - 1) / int64(limit)
	}
	return RewardsListResponse{
		Items: items,
		Meta:  common.PaginationMeta{Page: page, Limit: limit, Total: result.Total, TotalPages: totalPages},
	}
}
