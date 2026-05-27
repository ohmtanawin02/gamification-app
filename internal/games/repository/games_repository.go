package repository

import (
	"context"
	"gamification-app/internal/games/domain"
	"gamification-app/pkg/common"

	"gorm.io/gorm"
)

type GamesRepository struct {
	db *gorm.DB
}

type GamesRepositoryCfg struct {
	DB *gorm.DB
}

func NewGamesRepository(cfg GamesRepositoryCfg) domain.GamesRepository {
	return &GamesRepository{
		db: cfg.DB,
	}
}

func (r *GamesRepository) Spin(ctx context.Context, input domain.SpinInput) (*domain.SpinResult, error) {
	log := common.NewRepoLogger(ctx, "GamesRepository.Spin")
	var result domain.SpinResult

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var totalPoints int
		if err := tx.Raw("SELECT total_points FROM users WHERE id = ?", input.UserID).Scan(&totalPoints).Error; err != nil {
			log.Error().Err(err).Uint64("user_id", uint64(input.UserID)).Msg("failed to get user points")
			return err
		}

		if totalPoints+input.PointsEarned > 10000 {
			return domain.ErrPointsExceeded
		}

		if err := tx.Exec("INSERT INTO user_game_histories (user_id, points_earned) VALUES (?, ?)", input.UserID, input.PointsEarned).Error; err != nil {
			log.Error().Err(err).Uint64("user_id", uint64(input.UserID)).Msg("failed to insert into user_game_histories")
			return err
		}

		if err := tx.Exec("UPDATE users SET total_points = total_points + ? WHERE id = ?", input.PointsEarned, input.UserID).Error; err != nil {
			log.Error().Err(err).Uint64("user_id", uint64(input.UserID)).Msg("failed to update user points")
			return err
		}

		result.TotalPoints = totalPoints + input.PointsEarned
		return nil
	})

	return &result, err
}
