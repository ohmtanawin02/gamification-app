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

type FindAllRewardsHandlerCfg struct {
	Service domain.RewardsService
}

// FindAllRewards godoc
// @Summary      Get all rewards
// @Tags         rewards
// @Produce      json
// @Param        page        query  int     false  "Page"
// @Param        limit       query  int     false  "Limit"
// @Param        search      query  string  false  "Search by name"
// @Param        sort_order  query  string  false  "asc or desc"
// @Success      200
// @Failure      500
// @Security     BearerAuth
// @Router       /api/v1/rewards [get]
func FindAllRewards(cfg FindAllRewardsHandlerCfg) fiber.Handler {
	return func(c *fiber.Ctx) error {
		page, _ := strconv.Atoi(c.Query("page", "1"))
		limit, _ := strconv.Atoi(c.Query("limit", "20"))
		search := c.Query("search", "")
		sortOrder := constants.SortOrder(c.Query("sort_order", string(constants.SortOrderAsc)))

		req := domain.FindAllRewardsRequest{
			Search:    search,
			Page:      page,
			Limit:     limit,
			SortOrder: sortOrder,
		}

		result, err := cfg.Service.FindAll(c.UserContext(), req)
		if err != nil {
			return common.ResponseJsonWithCode(c, fiber.StatusInternalServerError, uuid.New(),
				constants.CodeInternalError, constants.MessageENSomethingWentWrong, constants.MessageTHSomethingWentWrong, nil)
		}

		return common.ResponseJsonWithCode(c, fiber.StatusOK, uuid.Nil,
			constants.CodeOK, constants.MessageENSuccess, constants.MessageTHSuccess,
			dto.ToRewardsListResponse(result, page, limit))
	}
}
