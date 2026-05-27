package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"gamification-app/internal/users/domain"
	"gamification-app/internal/users/handler/dto"
	"gamification-app/pkg/auth"
	"gamification-app/pkg/common"
	"gamification-app/pkg/constants"
)

type GetMeHandlerCfg struct {
	Service domain.UsersService
}

// GetMe godoc
// @Summary      Get current user profile with rewards
// @Tags         users
// @Produce      json
// @Success      200
// @Failure      401
// @Failure      404
// @Failure      500
// @Security     BearerAuth
// @Router       /api/v1/users/me [get]
func GetMe(cfg GetMeHandlerCfg) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := auth.GetUserID(c.UserContext())
		if !ok {
			return common.ResponseJsonWithCode(c, fiber.StatusUnauthorized, uuid.New(),
				constants.CodeUnauthorized, constants.MessageENUnauthorized, constants.MessageTHUnauthorized, nil)
		}
		nickname, _ := auth.GetNickname(c.UserContext())

		result, err := cfg.Service.GetMe(c.UserContext(), userID, nickname)
		if err != nil {
			return common.ResponseJsonWithCode(c, fiber.StatusInternalServerError, uuid.New(),
				constants.CodeInternalError, constants.MessageENSomethingWentWrong, constants.MessageTHSomethingWentWrong, nil)
		}
		if result.User == nil {
			return common.ResponseJsonWithCode(c, fiber.StatusNotFound, uuid.New(),
				constants.CodeNotFound, constants.MessageENNotFound, constants.MessageTHNotFound, nil)
		}

		rewards := make([]dto.RewardItem, len(result.Rewards))
		for i, r := range result.Rewards {
			rewards[i] = dto.RewardItem{
				ID:         r.ID,
				Name:       r.Name,
				CheckPoint: r.CheckPoint,
				Claimed:    r.Claimed,
			}
		}

		return common.ResponseJsonWithCode(c, fiber.StatusOK, uuid.Nil,
			constants.CodeOK, constants.MessageENSuccess, constants.MessageTHSuccess,
			dto.UserMeResponse{
				Nickname:    result.User.Nickname,
				TotalPoints: result.User.TotalPoints,
				Rewards:     rewards,
			})
	}
}
