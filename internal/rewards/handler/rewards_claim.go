package handler

import (
	"errors"
	"gamification-app/internal/rewards/domain"
	"gamification-app/internal/rewards/handler/dto"
	"gamification-app/pkg/auth"
	"gamification-app/pkg/common"
	"gamification-app/pkg/constants"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ClaimRewardHandlerCfg struct {
	Service  domain.RewardsService
	Validate *validator.Validate
}

func ClaimReward(cfg ClaimRewardHandlerCfg) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req dto.ClaimRewardRequest
		if err := c.BodyParser(&req); err != nil {
			return common.ResponseJsonWithCode(c, fiber.StatusBadRequest, uuid.New(),
				constants.CodeBadRequest, constants.MessageENBadRequest, constants.MessageTHBadRequest, nil)
		}
		if err := req.Validate(cfg.Validate); err != nil {
			return common.ResponseJsonWithCode(c, fiber.StatusBadRequest, uuid.New(),
				constants.CodeBadRequest, err.Error(), constants.MessageTHBadRequest, nil)
		}

		userID, _ := auth.GetUserID(c.UserContext())
		input := req.ToDomainInput(userID)
		err := cfg.Service.ClaimReward(c.UserContext(), input.UserID, input.RewardID)
		if err != nil {
			switch {
			case errors.Is(err, domain.ErrRewardNotFound):
				return common.ResponseJsonWithCode(c, fiber.StatusNotFound, uuid.New(),
					constants.CodeNotFound, constants.MessageENNotFound, constants.MessageTHNotFound, nil)
			case errors.Is(err, domain.ErrInsufficientPoints):
				return common.ResponseJsonWithCode(c, fiber.StatusUnprocessableEntity, uuid.New(),
					constants.CodeUnprocessable, constants.MessageENInsufficientPoints, constants.MessageTHInsufficientPoints, nil)
			case errors.Is(err, domain.ErrRewardAlreadyClaimed):
				return common.ResponseJsonWithCode(c, fiber.StatusConflict, uuid.New(),
					constants.CodeConflict, constants.MessageENRewardAlreadyClaimed, constants.MessageTHRewardAlreadyClaimed, nil)
			}
			return common.ResponseJsonWithCode(c, fiber.StatusInternalServerError, uuid.New(),
				constants.CodeInternalError, constants.MessageENSomethingWentWrong, constants.MessageTHSomethingWentWrong, nil)
		}

		return common.ResponseJsonWithCode(c, fiber.StatusOK, uuid.Nil,
			constants.CodeOK, constants.MessageENSuccess, constants.MessageTHSuccess, nil)
	}
}
