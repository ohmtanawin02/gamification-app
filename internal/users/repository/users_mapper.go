package repository

import (
	"gamification-app/internal/users/domain"
	"gamification-app/internal/users/repository/models"
	"time"
)

func toUserEntity(u models.Users) domain.Users {
	var deletedAt *time.Time
	if u.DeletedAt.Valid {
		time := u.DeletedAt.Time
		deletedAt = &time
	}
	return domain.Users{
		ID:          u.ID,
		Nickname:    u.Nickname,
		TotalPoints: u.TotalPoints,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		DeletedAt:   deletedAt,
	}
}

func toUserEntities(users []models.Users) []domain.Users {
	result := make([]domain.Users, len(users))
	for i, u := range users {
		result[i] = toUserEntity(u)
	}
	return result
}

func toUserModel(u domain.Users) models.Users {
	return models.Users{
		ID:          u.ID,
		Nickname:    u.Nickname,
		TotalPoints: u.TotalPoints,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}	
}