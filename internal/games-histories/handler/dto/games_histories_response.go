package dto

import (
	"time"

	"gamification-app/internal/games-histories/domain"
	"gamification-app/pkg/common"
)

type GamesHistoryResponse struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	PointsEarned int       `json:"points_earned"`
	CreatedAt    time.Time `json:"created_at"`
}

type GamesHistoryListResponse struct {
	Items []GamesHistoryResponse `json:"items"`
	Meta  common.PaginationMeta  `json:"meta"`
}

func ToGamesHistoryResponse(r domain.UserGameHistory) GamesHistoryResponse {
	return GamesHistoryResponse{
		ID:           r.ID,
		UserID:       r.UserID,
		PointsEarned: r.PointsEarned,
		CreatedAt:    r.CreatedAt,
	}
}

func ToGamesHistoryListResponse(result domain.FindAllGameHistoryResult, page, limit int) GamesHistoryListResponse {
	items := make([]GamesHistoryResponse, 0, len(result.Items))
	for _, r := range result.Items {
		items = append(items, ToGamesHistoryResponse(r))
	}
	totalPages := int64(0)
	if limit > 0 {
		totalPages = (result.Total + int64(limit) - 1) / int64(limit)
	}
	return GamesHistoryListResponse{
		Items: items,
		Meta:  common.PaginationMeta{Page: page, Limit: limit, Total: result.Total, TotalPages: totalPages},
	}
}
