package domain

import "context"

type UsersService interface {
	Login(ctx context.Context, input LoginInput) (*LoginResult, error)
	FindAll(ctx context.Context, req FindAllUsersRequest) (FindAllUsersResult, error)
	FindUserByNickname(ctx context.Context, nickname string) (*Users, error)
}