package server

import (
	"errors"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"

	"gamification-app/config"
	rewardsHandler "gamification-app/internal/rewards/handler"
	usersHandler "gamification-app/internal/users/handler"
	"gamification-app/pkg/common"
	"gamification-app/pkg/constants"
	"gamification-app/pkg/middleware"
)

func NewServer(cfg *config.Config) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}
			return common.ResponseJsonWithCode(c, code, uuid.New(),
				constants.CodeInternalError,
				constants.MessageENSomethingWentWrong,
				constants.MessageTHSomethingWentWrong,
				nil)
		},
		ReadBufferSize: 4 * 1024 * 1024,
		BodyLimit:      10 * 1024 * 1024,
	})

	app.Use(cors.New())

	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	app.Use(fiberzerolog.New(fiberzerolog.Config{Logger: &logger}))
	app.Use(func(c *fiber.Ctx) error {
		return common.SetLogCtx(c, logger)
	})

	readDB, err := cfg.DBRead.Connect(cfg.Debug)
	if err != nil {
		logger.Panic().Err(err).Msg("failed to connect read db")
	}
	writeDB, err := cfg.DBWrite.Connect(cfg.Debug)
	if err != nil {
		logger.Panic().Err(err).Msg("failed to connect write db")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		readSQLDB, _ := readDB.DB()
		if err := readSQLDB.Ping(); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "unhealthy", "read_db": err.Error(),
			})
		}
		writeSQLDB, _ := writeDB.DB()
		if err := writeSQLDB.Ping(); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "unhealthy", "write_db": err.Error(),
			})
		}
		if err := rdb.Ping(c.UserContext()).Err(); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "unhealthy", "redis": err.Error(),
			})
		}
		return c.JSON(fiber.Map{"status": "ok"})
	})

	validate := validator.New()
	jwtTTL := time.Duration(cfg.JWTExpireHours) * time.Hour

	api := app.Group("/api/v1")
	api.Use(middleware.JWTProtected(cfg.JWTSecret))

	usersHandler.NewUsersRouterCfg{
		PublicApp:    app,
		ProtectedApp: api,
		WriteDB:      writeDB,
		Logger:       logger,
		Validate:     validate,
		JWTSecret:    cfg.JWTSecret,
		JWTTTL:       jwtTTL,
	}.NewUsersRouter()

	rewardsHandler.NewRewardsRouterCfg{
		PublicApp:    app,
		ProtectedApp: api,
		WriteDB:      writeDB,
		Logger:       logger,
		Validate:     validate,
	}.NewRewardsRouter()

	app.Use(func(c *fiber.Ctx) error {
		return common.ResponseJsonWithCode(c, fiber.StatusNotFound, uuid.New(),
			constants.CodeNotFound,
			constants.MessageENNotFound,
			constants.MessageTHNotFound,
			nil)
	})

	return app
}
