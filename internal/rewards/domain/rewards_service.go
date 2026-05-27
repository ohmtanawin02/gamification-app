package domain

import "context"

type RewardsService interface {
	ClaimReward(ctx context.Context, userID uint, rewardID uint) error
	FindAll(ctx context.Context, req FindAllRewardsRequest) (FindAllRewardsResult, error)
}
