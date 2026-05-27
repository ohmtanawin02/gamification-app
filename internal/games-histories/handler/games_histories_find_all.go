package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"gamification-app/internal/games-histories/domain"
	"gamification-app/internal/games-histories/handler/dto"
	"gamification-app/pkg/common"
	"gamification-app/pkg/constants"
)

type FindAllGamesHistoryHandlerCfg struct {
	Service domain.GamesHistoryService
}

// FindAllGamesHistory godoc
// @Summary      Get game histories
// @Tags         histories
// @Produce      json
// @Param        page        query  int     false  "Page"
// @Param        limit       query  int     false  "Limit"
// @Param        nickname    query  string  false  "Filter by nickname"
// @Param        sort_order  query  string  false  "asc or desc"
// @Success      200
// @Failure      500
// @Security     BearerAuth
// @Router       /api/v1/games-histories [get]
func FindAllGamesHistory(cfg FindAllGamesHistoryHandlerCfg) fiber.Handler {
	return func(c *fiber.Ctx) error {
		page, _ := strconv.Atoi(c.Query("page", "1"))
		limit, _ := strconv.Atoi(c.Query("limit", "20"))
		nickname := c.Query("nickname", "")
		sortOrder := constants.SortOrder(c.Query("sort_order", string(constants.SortOrderAsc)))

		req := domain.FindAllGameHistoryRequest{
			Nickname:  nickname,
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
			dto.ToGamesHistoryListResponse(result, page, limit))
	}
}
