package repository

import (
	"gamification-app/internal/rewards/domain"
	"gamification-app/internal/rewards/repository/models"
)

func toRewardEntity(reward models.Rewards) domain.Rewards {
	return domain.Rewards{
		ID:         reward.ID,
		Name:       reward.Name,
		CheckPoint: reward.CheckPoint,
		CreatedAt:  reward.CreatedAt,
		UpdatedAt:  reward.UpdatedAt,
	}
}

func toRewardEntities(rewards []models.Rewards) []domain.Rewards {
	result := make([]domain.Rewards, len(rewards))
	for i, reward := range rewards {
		result[i] = toRewardEntity(reward)
	}
	return result
}
