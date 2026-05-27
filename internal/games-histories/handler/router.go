package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"gamification-app/internal/games-histories/repository"
	"gamification-app/internal/games-histories/service"
)

type NewGamesHistoryRouterCfg struct {
	PublicApp    fiber.Router
	ProtectedApp fiber.Router
	WriteDB      *gorm.DB
	Logger       zerolog.Logger
	Validate     *validator.Validate
}

func (cfg NewGamesHistoryRouterCfg) NewGamesHistoryRouter() {
	repo := repository.NewGamesHistoryRepository(repository.GamesHistoryRepositoryCfg{
		DB: cfg.WriteDB,
	})
	svc := service.NewGamesHistoryService(service.GamesHistoryServiceCfg{Repo: repo})

	cfg.ProtectedApp.Get("/games-histories", FindAllGamesHistory(FindAllGamesHistoryHandlerCfg{Service: svc}))
}
