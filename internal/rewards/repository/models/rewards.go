package models

import (
	"time"
)

type Rewards struct {
	ID         uint   `gorm:"primaryKey"`
	Name       string `gorm:"not null"`
	CheckPoint int    `gorm:"not null;default:0"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (Rewards) TableName() string {
	return "rewards"
}
