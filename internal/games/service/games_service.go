package service

import (
	"context"
	"errors"
	"gamification-app/internal/games/domain"
	"gamification-app/pkg/common"
)

type GamesService struct {
	repo domain.GamesRepository
}

type GamesServiceCfg struct {
	Repo domain.GamesRepository
}

func NewGamesService(cfg GamesServiceCfg) domain.GamesService {
	return &GamesService{repo: cfg.Repo}
}

var validPoints = map[int]bool{300: true, 500: true, 1000: true, 3000: true}

func (s *GamesService) Spin(ctx context.Context, req domain.SpinInput) (*domain.SpinResult, error) {
	log := common.NewAppLogger(ctx, "GamesService.Spin")

	if !validPoints[req.PointsEarned] {
		return &domain.SpinResult{}, domain.ErrInvalidPoints
	}

	result, err := s.repo.Spin(ctx, req)
	if err != nil {
		if !errors.Is(err, domain.ErrPointsExceeded) && !errors.Is(err, domain.ErrInvalidPoints) {
			log.Error().Err(err).Msg("spin failed")
		}
		return &domain.SpinResult{}, err
	}

	return result, nil
}
