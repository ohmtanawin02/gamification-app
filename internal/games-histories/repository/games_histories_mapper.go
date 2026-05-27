package repository

import (
	"gamification-app/internal/games-histories/domain"
	"gamification-app/internal/games-histories/repository/models"
)

func toUserGameHistoryEntity(u models.UserGameHistory) domain.UserGameHistory {
	return domain.UserGameHistory{
		ID:           u.ID,
		UserID:       u.UserID,
		PointsEarned: u.PointsEarned,
		CreatedAt:    u.CreatedAt,
	}
}

func toUserGameHistoryEntities(histories []models.UserGameHistory) []domain.UserGameHistory {
	result := make([]domain.UserGameHistory, len(histories))
	for i, u := range histories {
		result[i] = toUserGameHistoryEntity(u)
	}
	return result
}
