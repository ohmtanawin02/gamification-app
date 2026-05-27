package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"gamification-app/internal/users/domain"
	"gamification-app/internal/users/handler/dto"
	"gamification-app/pkg/common"
	"gamification-app/pkg/constants"
)

type LoginHandlerCfg struct {
	Service  domain.UsersService
	Validate *validator.Validate
}

func Login(cfg LoginHandlerCfg) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req dto.LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return common.ResponseJsonWithCode(c, fiber.StatusBadRequest, uuid.New(),
				constants.CodeBadRequest, constants.MessageENBadRequest, constants.MessageTHBadRequest, nil)
		}
		if err := req.Validate(cfg.Validate); err != nil {
			return common.ResponseJsonWithCode(c, fiber.StatusBadRequest, uuid.New(),
				constants.CodeBadRequest, err.Error(), constants.MessageTHBadRequest, nil)
		}

		result, err := cfg.Service.Login(c.UserContext(), req.ToDomainInput())
		if err != nil {
			return common.ResponseJsonWithCode(c, fiber.StatusInternalServerError, uuid.New(),
				constants.CodeInternalError, constants.MessageENSomethingWentWrong, constants.MessageTHSomethingWentWrong, nil)
		}

		return common.ResponseJsonWithCode(c, fiber.StatusOK, uuid.Nil,
			constants.CodeOK, constants.MessageENSuccess, constants.MessageTHSuccess,
			dto.ToLoginResponse(*result))
	}
}
