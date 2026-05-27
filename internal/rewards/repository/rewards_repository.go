package repository

import (
	"context"
	"time"

	"gamification-app/internal/rewards/domain"
	"gamification-app/internal/rewards/repository/models"
	"gamification-app/pkg/common"
	"gamification-app/pkg/constants"

	"gorm.io/gorm"
)

type RewardsRepository struct {
	db *gorm.DB
}

type RewardsRepositoryCfg struct {
	DB *gorm.DB
}

func NewRewardsRepository(cfg RewardsRepositoryCfg) domain.RewardsRepository {
	return &RewardsRepository{
		db: cfg.DB,
	}
}

func (r *RewardsRepository) FindAll(ctx context.Context, req domain.FindAllRewardsRequest) (domain.FindAllRewardsResult, error) {
	log := common.NewRepoLogger(ctx, "RewardsRepository.FindAll")

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 || req.Limit > 100 {
		req.Limit = 20
	}
	if !req.SortOrder.IsValid() {
		req.SortOrder = constants.SortOrderAsc
	}

	query := r.db.WithContext(ctx).Model(&models.Rewards{})
	if req.Search != "" {
		query = query.Where("name ILIKE ?", "%"+req.Search+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		log.Error().Err(err).Msg("count failed")
		return domain.FindAllRewardsResult{}, err
	}

	var ms []models.Rewards
	offset := (req.Page - 1) * req.Limit
	if err := query.Order("id " + req.SortOrder.SQL()).Offset(offset).Limit(req.Limit).Find(&ms).Error; err != nil {
		log.Error().Err(err).Msg("find failed")
		return domain.FindAllRewardsResult{}, err
	}

	return domain.FindAllRewardsResult{Items: toRewardEntities(ms), Total: total}, nil
}

func (r *RewardsRepository) FindUserRewards(ctx context.Context, req domain.FindUserRewardsRequest) (domain.FindUserRewardsResult, error) {
	log := common.NewRepoLogger(ctx, "RewardsRepository.FindUserRewards")

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 || req.Limit > 100 {
		req.Limit = 20
	}

	type row struct {
		ID         uint
		Name       string
		CheckPoint int
		ClaimedAt  time.Time
	}

	baseQuery := r.db.WithContext(ctx).
		Table("user_rewards ur").
		Joins("JOIN rewards r ON r.id = ur.reward_id").
		Joins("JOIN users u ON u.id = ur.user_id")

	if req.Nickname != "" {
		baseQuery = baseQuery.Where("u.nickname ILIKE ?", "%"+req.Nickname+"%")
	} else if req.UserID > 0 {
		baseQuery = baseQuery.Where("ur.user_id = ?", req.UserID)
	}

	var total int64
	if err := baseQuery.Count(&total).Error; err != nil {
		log.Error().Err(err).Msg("count failed")
		return domain.FindUserRewardsResult{}, err
	}

	var rows []row
	offset := (req.Page - 1) * req.Limit
	if err := baseQuery.
		Select("r.id, r.name, r.check_point, ur.claimed_at").
		Offset(offset).Limit(req.Limit).
		Order("ur.claimed_at DESC").
		Scan(&rows).Error; err != nil {
		log.Error().Err(err).Msg("find failed")
		return domain.FindUserRewardsResult{}, err
	}

	items := make([]domain.UserReward, len(rows))
	for i, r := range rows {
		items[i] = domain.UserReward{
			ID:         r.ID,
			Name:       r.Name,
			CheckPoint: r.CheckPoint,
			ClaimedAt:  r.ClaimedAt,
		}
	}

	return domain.FindUserRewardsResult{Items: items, Total: total}, nil
}

func (r *RewardsRepository) ClaimReward(ctx context.Context, userID uint, rewardID uint) error {
	log := common.NewRepoLogger(ctx, "RewardsRepository.ClaimReward")
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var reward models.Rewards
		if err := tx.Where("id = ?", rewardID).First(&reward).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return domain.ErrRewardNotFound
			}
			log.Error().Err(err).Uint64("reward_id", uint64(rewardID)).Msg("failed to find reward")
			return err
		}

		var userPoints int64
		if err := tx.Raw("SELECT total_points FROM users WHERE id = ?", userID).Scan(&userPoints).Error; err != nil {
			log.Error().Err(err).Uint64("user_id", uint64(userID)).Msg("failed to get user points")
			return err
		}
		if userPoints < int64(reward.CheckPoint) {
			return domain.ErrInsufficientPoints
		}

		var count int64
		if err := tx.Raw("SELECT COUNT(*) FROM user_rewards WHERE user_id = ? AND reward_id = ?", userID, rewardID).Scan(&count).Error; err != nil {
			log.Error().Err(err).Uint64("user_id", uint64(userID)).Uint64("reward_id", uint64(rewardID)).Msg("failed to check if reward is already claimed")
			return err
		}
		if count > 0 {
			return domain.ErrRewardAlreadyClaimed
		}

		if err := tx.Exec("INSERT INTO user_rewards (user_id, reward_id) VALUES (?, ?)", userID, rewardID).Error; err != nil {
			log.Error().Err(err).Uint64("user_id", uint64(userID)).Uint64("reward_id", uint64(rewardID)).Msg("failed to insert into user_rewards")
			return err
		}

		return nil
	})
}
