package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"gamification-app/internal/rewards/domain"
	"gamification-app/internal/rewards/handler/dto"
	"gamification-app/pkg/common"
	"gamification-app/pkg/constants"
)

type FindAllUserRewardsHandlerCfg struct {
	Service domain.RewardsService
}

func FindAllUserRewards(cfg FindAllUserRewardsHandlerCfg) fiber.Handler {
	return func(c *fiber.Ctx) error {
		page, _ := strconv.Atoi(c.Query("page", "1"))
		limit, _ := strconv.Atoi(c.Query("limit", "20"))
		nickname := c.Query("nickname", "")

		req := domain.FindUserRewardsRequest{
			Nickname: nickname,
			Page:     page,
			Limit:    limit,
		}

		result, err := cfg.Service.FindUserRewards(c.UserContext(), req)
		if err != nil {
			return common.ResponseJsonWithCode(c, fiber.StatusInternalServerError, uuid.New(),
				constants.CodeInternalError, constants.MessageENSomethingWentWrong, constants.MessageTHSomethingWentWrong, nil)
		}

		return common.ResponseJsonWithCode(c, fiber.StatusOK, uuid.Nil,
			constants.CodeOK, constants.MessageENSuccess, constants.MessageTHSuccess,
			dto.ToUserRewardsListResponse(result, page, limit))
	}
}
