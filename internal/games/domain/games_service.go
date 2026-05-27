package domain

import "context"

type GamesService interface {
	Spin(ctx context.Context, input SpinInput) (*SpinResult, error)
}
