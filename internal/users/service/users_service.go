package service

import (
	"context"
	"sync"

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

func (s *UsersService) FindUserRewardsByUserID(ctx context.Context, userID uint) ([]domain.UserRewardItem, error) {
	log := common.NewAppLogger(ctx, "UsersService.FindUserRewardsByUserID")

	items, err := s.repo.FindUserRewardsByUserID(ctx, userID)
	if err != nil {
		log.Error().Err(err).Uint("user_id", userID).Msg("find user rewards failed")
		return nil, err
	}

	return items, nil
}

func (s *UsersService) GetMe(ctx context.Context, userID uint, nickname string) (*domain.GetMeResult, error) {
	log := common.NewAppLogger(ctx, "UsersService.GetMe")

	var (
		user        *domain.Users
		rewards     []domain.UserRewardItem
		userErr     error
		rewardsErr  error
		wg          sync.WaitGroup
	)

	wg.Add(2)

	go func() {
		defer wg.Done()
		user, userErr = s.repo.FindUserByNickname(ctx, nickname)
	}()

	go func() {
		defer wg.Done()
		rewards, rewardsErr = s.repo.FindUserRewardsByUserID(ctx, userID)
	}()

	wg.Wait()

	if userErr != nil {
		log.Error().Err(userErr).Str("nickname", nickname).Msg("find user failed")
		return nil, userErr
	}
	if rewardsErr != nil {
		log.Error().Err(rewardsErr).Uint("user_id", userID).Msg("find rewards failed")
		return nil, rewardsErr
	}

	return &domain.GetMeResult{User: user, Rewards: rewards}, nil
}
