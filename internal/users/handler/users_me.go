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

func GetMe(cfg GetMeHandlerCfg) fiber.Handler {
	return func(c *fiber.Ctx) error {
		nickname, ok := auth.GetNickname(c.UserContext())
		if !ok {
			return common.ResponseJsonWithCode(c, fiber.StatusUnauthorized, uuid.New(),
				constants.CodeUnauthorized, constants.MessageENUnauthorized, constants.MessageTHUnauthorized, nil)
		}

		user, err := cfg.Service.FindUserByNickname(c.UserContext(), nickname)
		if err != nil {
			return common.ResponseJsonWithCode(c, fiber.StatusInternalServerError, uuid.New(),
				constants.CodeInternalError, constants.MessageENSomethingWentWrong, constants.MessageTHSomethingWentWrong, nil)
		}
		if user == nil {
			return common.ResponseJsonWithCode(c, fiber.StatusNotFound, uuid.New(),
				constants.CodeNotFound, constants.MessageENNotFound, constants.MessageTHNotFound, nil)
		}

		rewardItems, err := cfg.Service.FindUserRewardsByUserID(c.UserContext(), user.ID)
		if err != nil {
			return common.ResponseJsonWithCode(c, fiber.StatusInternalServerError, uuid.New(),
				constants.CodeInternalError, constants.MessageENSomethingWentWrong, constants.MessageTHSomethingWentWrong, nil)
		}

		rewards := make([]dto.RewardItem, len(rewardItems))
		for i, r := range rewardItems {
			rewards[i] = dto.RewardItem{
				ID:         r.ID,
				Name:       r.Name,
				CheckPoint: r.CheckPoint,
				Claimed:    r.Claimed,
			}
		}

		resp := dto.UserMeResponse{
			Nickname:    user.Nickname,
			TotalPoints: user.TotalPoints,
			Rewards:     rewards,
		}

		return common.ResponseJsonWithCode(c, fiber.StatusOK, uuid.Nil,
			constants.CodeOK, constants.MessageENSuccess, constants.MessageTHSuccess, resp)
	}
}
