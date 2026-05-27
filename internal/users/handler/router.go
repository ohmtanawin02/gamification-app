package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"gamification-app/config"
	"gamification-app/internal/users/repository"
	"gamification-app/internal/users/service"
	"gamification-app/pkg/cache"
)

type NewUsersRouterCfg struct {
	PublicApp    fiber.Router 
	ProtectedApp fiber.Router 
	WriteDB      *gorm.DB
	Redis        *redis.Client
	Logger       zerolog.Logger
	Validate     *validator.Validate
	Cfg          *config.Config
}

func (cfg NewUsersRouterCfg) NewUsersRouter() {
	repo := repository.NewUsersRepository(repository.UsersRepositoryCfg{
		DB:             cfg.WriteDB,
		JWTSecret:      cfg.Cfg.JWTSecret,
		JWTExpireHours: cfg.Cfg.JWTExpireHours,
	})
	_ = cache.NewRedisCache(cfg.Redis)
	svc := service.NewUsersService(service.UsersServiceCfg{Repo: repo})

	cfg.PublicApp.Post("/auth/login", Login(LoginHandlerCfg{Service: svc, Validate: cfg.Validate}))

	cfg.ProtectedApp.Get("/users", FindAllUsers(FindAllUsersHandlerCfg{Service: svc}))
	cfg.ProtectedApp.Get("/users/me", GetMe(GetMeHandlerCfg{Service: svc}))
}
