package domain

import (
	"errors"
	"gamification-app/pkg/constants"
	"time"
)

type ClaimRewardInput struct {
	UserID   uint
	RewardID uint
}

type Rewards struct {
	ID         uint
	Name       string
	CheckPoint int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type FindAllRewardsRequest struct {
	Search    string
	Page      int
	Limit     int
	SortOrder constants.SortOrder
}

type FindAllRewardsResult struct {
	Items []Rewards
	Total int64
}

type FindUserRewardsRequest struct {
	Nickname string
	Page     int
	Limit    int
}

type UserReward struct {
	ID         uint
	Name       string
	CheckPoint int
	ClaimedAt  time.Time
}

type FindUserRewardsResult struct {
	Items []UserReward
	Total int64
}

var (
	ErrRewardNotFound       = errors.New("reward not found")
	ErrInsufficientPoints   = errors.New("insufficient points")
	ErrRewardAlreadyClaimed = errors.New("reward already claimed")
)
