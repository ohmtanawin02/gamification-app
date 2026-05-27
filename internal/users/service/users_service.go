package service

import (
	"context"

	"gamification-app/internal/users/domain"
	"gamification-app/pkg/common"
)

type UsersService struct {
	repo domain.UserRepository
}

type UsersServiceCfg struct {
	Repo domain.UserRepository
}

func NewUsersService(cfg UsersServiceCfg) domain.UsersService {
	return &UsersService{repo: cfg.Repo}
}

func (s *UsersService) Login(ctx context.Context, input domain.LoginInput) (*domain.LoginResult, error) {
	log := common.NewAppLogger(ctx, "UsersService.Login")

	result, err := s.repo.Login(ctx, input)
	if err != nil {
		log.Error().Err(err).Str("nickname", input.Nickname).Msg("login failed")
		return nil, err
	}

	return result, nil
}

func (s *UsersService) FindAll(ctx context.Context, req domain.FindAllUsersRequest) (domain.FindAllUsersResult, error) {
	log := common.NewAppLogger(ctx, "UsersService.FindAll")

	result, err := s.repo.FindAll(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("find all failed")
		return domain.FindAllUsersResult{}, err
	}

	return result, nil
}

func (s *UsersService) FindUserByNickname(ctx context.Context, nickname string) (*domain.Users, error) {
	log := common.NewAppLogger(ctx, "UsersService.FindUserByNickname")

	user, err := s.repo.FindUserByNickname(ctx, nickname)
	if err != nil {
		log.Error().Err(err).Str("nickname", nickname).Msg("find by nickname failed")
		return nil, err
	}

	return user, nil
}
