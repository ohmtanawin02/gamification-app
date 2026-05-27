package domain

import "context"

type GamesHistoryRepository interface {
	FindAll(ctx context.Context, req FindAllGameHistoryRequest) (FindAllGameHistoryResult, error)
}
