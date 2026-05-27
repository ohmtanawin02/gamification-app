package domain

import "context"

type UserRepository interface {
	Login(ctx context.Context, req LoginInput) (*LoginResult, error)
	FindAll(ctx context.Context, req FindAllUsersRequest) (FindAllUsersResult, error)
	FindUserByNickname(ctx context.Context, nickname string) (*Users, error)
}
