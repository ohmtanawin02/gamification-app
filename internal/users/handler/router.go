package handler

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"gamification-app/internal/users/repository"
	"gamification-app/internal/users/service"
)

type NewUsersRouterCfg struct {
	PublicApp    fiber.Router
	ProtectedApp fiber.Router
	WriteDB      *gorm.DB
	Logger       zerolog.Logger
	Validate     *validator.Validate
	JWTSecret    string
	JWTTTL       time.Duration
}

func (cfg NewUsersRouterCfg) NewUsersRouter() {
	repo := repository.NewUsersRepository(repository.UsersRepositoryCfg{
		DB:        cfg.WriteDB,
		JWTSecret: cfg.JWTSecret,
		JWTTTL:    cfg.JWTTTL,
	})
	svc := service.NewUsersService(service.UsersServiceCfg{Repo: repo})

	cfg.PublicApp.Post("/auth/login", Login(LoginHandlerCfg{Service: svc, Validate: cfg.Validate}))

	cfg.ProtectedApp.Get("/users", FindAllUsers(FindAllUsersHandlerCfg{Service: svc}))
	cfg.ProtectedApp.Get("/users/me", GetMe(GetMeHandlerCfg{Service: svc}))
}
