package models

import "time"

type UserGameHistory struct {
	ID           uint      `gorm:"primaryKey"`
	UserID       uint      `gorm:"index;not null"`
	PointsEarned int       `gorm:"not null"`
	CreatedAt    time.Time `gorm:"index"`
}

func (UserGameHistory) TableName() string {
	return "user_game_histories"
}
