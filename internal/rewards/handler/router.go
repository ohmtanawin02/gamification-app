package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"gamification-app/internal/rewards/repository"
	"gamification-app/internal/rewards/service"
)

type NewRewardsRouterCfg struct {
	PublicApp    fiber.Router
	ProtectedApp fiber.Router
	WriteDB      *gorm.DB
	Logger       zerolog.Logger
	Validate     *validator.Validate
}

func (cfg NewRewardsRouterCfg) NewRewardsRouter() {
	repo := repository.NewRewardsRepository(repository.RewardsRepositoryCfg{
		DB: cfg.WriteDB,
	})
	svc := service.NewRewardsService(service.RewardsServiceCfg{Repo: repo})

	cfg.ProtectedApp.Get("/rewards", FindAllRewards(FindAllRewardsHandlerCfg{Service: svc}))
	cfg.ProtectedApp.Post("/rewards/claim", ClaimReward(ClaimRewardHandlerCfg{Service: svc, Validate: cfg.Validate}))
}
