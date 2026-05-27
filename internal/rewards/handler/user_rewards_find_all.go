package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"gamification-app/internal/rewards/domain"
	"gamification-app/internal/rewards/handler/dto"
	"gamification-app/pkg/auth"
	"gamification-app/pkg/common"
	"gamification-app/pkg/constants"
)

type FindAllUserRewardsHandlerCfg struct {
	Service domain.RewardsService
}

// FindAllUserRewards godoc
// @Summary      Get claimed rewards by user
// @Tags         rewards
// @Produce      json
// @Param        page      query  int     false  "Page"
// @Param        limit     query  int     false  "Limit"
// @Param        nickname  query  string  false  "Filter by nickname"
// @Success      200
// @Failure      500
// @Security     BearerAuth
// @Router       /api/v1/user-rewards [get]
func FindAllUserRewards(cfg FindAllUserRewardsHandlerCfg) fiber.Handler {
	return func(c *fiber.Ctx) error {
		page, _ := strconv.Atoi(c.Query("page", "1"))
		limit, _ := strconv.Atoi(c.Query("limit", "20"))
		nickname := c.Query("nickname", "")
		userID, _ := auth.GetUserID(c.UserContext())

		req := domain.FindUserRewardsRequest{
			UserID:   userID,
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
