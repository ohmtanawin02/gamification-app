package models

import (
	"time"

	"gorm.io/gorm"
)

type Users struct {
	ID uint `gorm:"primaryKey"`
	Nickname string `gorm:"index;not null"`
	TotalPoints int `gorm:"not null;default:0"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt 
}

func (Users) TableName() string {
	return "users"
}