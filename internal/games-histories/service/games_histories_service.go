package service

import (
	"context"
	"gamification-app/internal/games-histories/domain"
	"gamification-app/pkg/common"
)

type GamesHistoryService struct {
	repo domain.GamesHistoryRepository
}

type GamesHistoryServiceCfg struct {
	Repo domain.GamesHistoryRepository
}

func NewGamesHistoryService(cfg GamesHistoryServiceCfg) domain.GamesHistoryService {
	return &GamesHistoryService{repo: cfg.Repo}
}

func (s *GamesHistoryService) FindAll(ctx context.Context, req domain.FindAllGameHistoryRequest) (domain.FindAllGameHistoryResult, error) {
	log := common.NewAppLogger(ctx, "GamesHistoryService.FindAll")

	result, err := s.repo.FindAll(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("find all failed")
		return domain.FindAllGameHistoryResult{}, err
	}

	return result, nil
}
