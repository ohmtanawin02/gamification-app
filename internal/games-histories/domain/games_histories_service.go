package domain

import "context"

type GamesHistoryService interface {
	FindAll(ctx context.Context, req FindAllGameHistoryRequest) (FindAllGameHistoryResult, error)
}
