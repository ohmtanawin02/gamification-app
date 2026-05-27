package handler

import (
	"errors"
	"gamification-app/internal/games/domain"
	"gamification-app/internal/games/handler/dto"
	"gamification-app/pkg/auth"
	"gamification-app/pkg/common"
	"gamification-app/pkg/constants"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type GamesSpinHandlerCfg struct {
	Service  domain.GamesService
	Validate *validator.Validate
}

// Spin godoc
// @Summary      Spin the wheel to earn points
// @Tags         games
// @Accept       json
// @Produce      json
// @Param        body  body  dto.SpinRequest  true  "Points earned (300, 500, 1000, 3000)"
// @Success      200
// @Failure      400
// @Failure      422
// @Failure      500
// @Security     BearerAuth
// @Router       /api/v1/games/spin [post]
func Spin(cfg GamesSpinHandlerCfg) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req dto.SpinRequest
		if err := c.BodyParser(&req); err != nil {
			return common.ResponseJsonWithCode(c, fiber.StatusBadRequest, uuid.New(),
				constants.CodeBadRequest, constants.MessageENBadRequest, constants.MessageTHBadRequest, nil)
		}
		if err := req.Validate(cfg.Validate); err != nil {
			return common.ResponseJsonWithCode(c, fiber.StatusBadRequest, uuid.New(),
				constants.CodeBadRequest, err.Error(), constants.MessageTHBadRequest, nil)
		}

		userID, _ := auth.GetUserID(c.UserContext())
		result, err := cfg.Service.Spin(c.UserContext(), req.ToDomainInput(userID))
		if err != nil {
			switch {
			case errors.Is(err, domain.ErrPointsExceeded):
				return common.ResponseJsonWithCode(c, fiber.StatusUnprocessableEntity, uuid.New(),
					constants.CodeUnprocessable, constants.MessageENPointsExceeded, constants.MessageTHPointsExceeded, nil)
			case errors.Is(err, domain.ErrInvalidPoints):
				return common.ResponseJsonWithCode(c, fiber.StatusUnprocessableEntity, uuid.New(),
					constants.CodeUnprocessable, constants.MessageENInvalidPoints, constants.MessageTHInvalidPoints, nil)
			}
			return common.ResponseJsonWithCode(c, fiber.StatusInternalServerError, uuid.New(),
				constants.CodeInternalError, constants.MessageENSomethingWentWrong, constants.MessageTHSomethingWentWrong, nil)
		}

		return common.ResponseJsonWithCode(c, fiber.StatusOK, uuid.Nil,
			constants.CodeOK, constants.MessageENSuccess, constants.MessageTHSuccess, result)
	}
}
