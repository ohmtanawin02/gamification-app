package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"gamification-app/internal/users/domain"
	"gamification-app/internal/users/repository/models"
	"gamification-app/pkg/auth"
	"gamification-app/pkg/common"
	"gamification-app/pkg/constants"
)

type UsersRepository struct {
	db        *gorm.DB
	jwtSecret string
	jwtTTL    time.Duration
}

type UsersRepositoryCfg struct {
	DB        *gorm.DB
	JWTSecret string
	JWTTTL    time.Duration
}

func NewUsersRepository(cfg UsersRepositoryCfg) domain.UserRepository {
	return &UsersRepository{
		db:        cfg.DB,
		jwtSecret: cfg.JWTSecret,
		jwtTTL:    cfg.JWTTTL,
	}
}

func (r *UsersRepository) Login(ctx context.Context, req domain.LoginInput) (*domain.LoginResult, error) {
	log := common.NewRepoLogger(ctx, "UsersRepository.Login")
	var m models.Users

	err := r.db.WithContext(ctx).Where("nickname = ?", req.Nickname).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		m = models.Users{Nickname: req.Nickname}
		if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
			log.Error().Err(err).Str("nickname", req.Nickname).Msg("failed to create user")
			return nil, err
		}
		log.Info().Uint("user_id", m.ID).Str("nickname", m.Nickname).Msg("new user created")
	} else if err != nil {
		log.Error().Err(err).Str("nickname", req.Nickname).Msg("failed to find user")
		return nil, err
	}

	token, err := auth.GenerateToken(m.ID, m.Nickname, r.jwtSecret, r.jwtTTL)
	if err != nil {
		log.Error().Err(err).Uint("user_id", m.ID).Msg("failed to generate token")
		return nil, err
	}

	user := toUserEntity(m)
	return &domain.LoginResult{Token: token, User: &user}, nil
}

func (r *UsersRepository) FindAll(ctx context.Context, req domain.FindAllUsersRequest) (domain.FindAllUsersResult, error) {
	log := common.NewRepoLogger(ctx, "UsersRepository.FindAll")

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 || req.Limit > 100 {
		req.Limit = 20
	}
	if !req.SortOrder.IsValid() {
		req.SortOrder = constants.SortOrderAsc
	}

	q := r.db.WithContext(ctx).Model(&models.Users{})
	if req.Nickname != "" {
		q = q.Where("nickname ILIKE ?", "%"+req.Nickname+"%")
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		log.Error().Err(err).Msg("count failed")
		return domain.FindAllUsersResult{}, err
	}

	var ms []models.Users
	offset := (req.Page - 1) * req.Limit
	if err := q.Order("id " + req.SortOrder.SQL()).Offset(offset).Limit(req.Limit).Find(&ms).Error; err != nil {
		log.Error().Err(err).Msg("find failed")
		return domain.FindAllUsersResult{}, err
	}

	return domain.FindAllUsersResult{Items: toUserEntities(ms), Total: total}, nil
}

func (r *UsersRepository) FindUserRewardsByUserID(ctx context.Context, userID uint) ([]domain.UserRewardItem, error) {
	log := common.NewRepoLogger(ctx, "UsersRepository.FindUserRewardsByUserID")

	type row struct {
		ID         uint
		Name       string
		CheckPoint int
		Claimed    bool
	}
	var rows []row
	err := r.db.WithContext(ctx).Raw(`
		SELECT r.id, r.name, r.check_point,
		       (ur.user_id IS NOT NULL) AS claimed
		FROM rewards r
		LEFT JOIN user_rewards ur ON ur.reward_id = r.id AND ur.user_id = ?
		ORDER BY r.check_point ASC
	`, userID).Scan(&rows).Error
	if err != nil {
		log.Error().Err(err).Uint("user_id", userID).Msg("failed to find user rewards")
		return nil, err
	}

	items := make([]domain.UserRewardItem, len(rows))
	for i, r := range rows {
		items[i] = domain.UserRewardItem{
			ID:         r.ID,
			Name:       r.Name,
			CheckPoint: r.CheckPoint,
			Claimed:    r.Claimed,
		}
	}
	return items, nil
}

func (r *UsersRepository) FindUserByNickname(ctx context.Context, nickname string) (*domain.Users, error) {
	log := common.NewRepoLogger(ctx, "UsersRepository.FindUserByNickname")

	var m models.Users
	if err := r.db.WithContext(ctx).Where("nickname = ?", nickname).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Error().Err(err).Str("nickname", nickname).Msg("find failed")
		return nil, err
	}
	user := toUserEntity(m)
	return &user, nil
}
