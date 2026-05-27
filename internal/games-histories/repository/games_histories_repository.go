package repository

import (
	"context"
	"gamification-app/internal/games-histories/domain"
	"gamification-app/internal/games-histories/repository/models"
	"gamification-app/pkg/common"
	"gamification-app/pkg/constants"

	"gorm.io/gorm"
)

type GamesHistoryRepository struct {
	db *gorm.DB
}

type GamesHistoryRepositoryCfg struct {
	DB *gorm.DB
}

func NewGamesHistoryRepository(cfg GamesHistoryRepositoryCfg) domain.GamesHistoryRepository {
	return &GamesHistoryRepository{
		db: cfg.DB,
	}
}

func (r *GamesHistoryRepository) FindAll(ctx context.Context, req domain.FindAllGameHistoryRequest) (domain.FindAllGameHistoryResult, error) {
	log := common.NewRepoLogger(ctx, "GamesHistoryRepository.FindAll")

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 || req.Limit > 100 {
		req.Limit = 20
	}
	if !req.SortOrder.IsValid() {
		req.SortOrder = constants.SortOrderAsc
	}

	query := r.db.WithContext(ctx).Model(&models.UserGameHistory{}).
		Joins("JOIN users ON users.id = user_game_histories.user_id")
	if req.Nickname != "" {
		query = query.Where("users.nickname ILIKE ?", "%"+req.Nickname+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		log.Error().Err(err).Msg("count failed")
		return domain.FindAllGameHistoryResult{}, err
	}

	var ms []models.UserGameHistory
	offset := (req.Page - 1) * req.Limit
	if err := query.Order("id " + req.SortOrder.SQL()).Offset(offset).Limit(req.Limit).Find(&ms).Error; err != nil {
		log.Error().Err(err).Msg("find failed")
		return domain.FindAllGameHistoryResult{}, err
	}

	return domain.FindAllGameHistoryResult{Items: toUserGameHistoryEntities(ms), Total: total}, nil
}
