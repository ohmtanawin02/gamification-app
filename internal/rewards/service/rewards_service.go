package service

import (
	"context"
	"gamification-app/internal/rewards/domain"
	"gamification-app/pkg/common"
)

type RewardsService struct {
	repo domain.RewardsRepository
}

type RewardsServiceCfg struct {
	Repo domain.RewardsRepository
}

func NewRewardsService(cfg RewardsServiceCfg) domain.RewardsService {
	return &RewardsService{repo: cfg.Repo}
}

func (s *RewardsService) FindAll(ctx context.Context, req domain.FindAllRewardsRequest) (domain.FindAllRewardsResult, error) {
	log := common.NewAppLogger(ctx, "RewardsService.FindAll")

	result, err := s.repo.FindAll(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("find all failed")
		return domain.FindAllRewardsResult{}, err
	}

	return result, nil
}

func (s *RewardsService) ClaimReward(ctx context.Context, userID uint, rewardID uint) error {
	log := common.NewAppLogger(ctx, "RewardsService.ClaimReward")

	err := s.repo.ClaimReward(ctx, userID, rewardID)
	if err != nil {
		log.Error().Err(err).Uint64("user_id", uint64(userID)).Uint64("reward_id", uint64(rewardID)).Msg("claim reward failed")
		return err
	}

	return nil
}
