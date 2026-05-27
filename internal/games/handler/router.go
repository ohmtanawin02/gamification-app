package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"gamification-app/internal/games/repository"
	"gamification-app/internal/games/service"
)

type NewGamesRouterCfg struct {
	PublicApp    fiber.Router
	ProtectedApp fiber.Router
	WriteDB      *gorm.DB
	Logger       zerolog.Logger
	Validate     *validator.Validate
}

func (cfg NewGamesRouterCfg) NewGamesRouter() {
	repo := repository.NewGamesRepository(repository.GamesRepositoryCfg{
		DB: cfg.WriteDB,
	})
	svc := service.NewGamesService(service.GamesServiceCfg{Repo: repo})

	cfg.ProtectedApp.Post("/games/spin", Spin(GamesSpinHandlerCfg{Service: svc, Validate: cfg.Validate}))
}
