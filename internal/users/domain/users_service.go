package domain

import "context"

type UsersService interface {
	Login(ctx context.Context, input LoginInput) (*LoginResult, error)
	FindAll(ctx context.Context, req FindAllUsersRequest) (FindAllUsersResult, error)
	FindUserByNickname(ctx context.Context, nickname string) (*Users, error)
	FindUserRewardsByUserID(ctx context.Context, userID uint) ([]UserRewardItem, error)
	GetMe(ctx context.Context, userID uint, nickname string) (*GetMeResult, error)
}