package domain

import "context"

type GamesRepository interface {
	Spin(ctx context.Context, input SpinInput) (*SpinResult, error)
}
