package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"gamification-app/pkg/auth"
	"gamification-app/pkg/common"
	"gamification-app/pkg/constants"
)

func JWTProtected(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == "" {
			return common.ResponseJsonWithCode(c, fiber.StatusUnauthorized, uuid.New(),
				constants.CodeUnauthorized, constants.MessageENUnauthorized, constants.MessageTHUnauthorized, nil)
		}
		claims, err := auth.ValidateToken(tokenStr, jwtSecret)
		if err != nil {
			return common.ResponseJsonWithCode(c, fiber.StatusUnauthorized, uuid.New(),
				constants.CodeUnauthorized, constants.MessageENUnauthorized, constants.MessageTHUnauthorized, nil)
		}

		ctx := auth.SetUserID(c.UserContext(), claims.UserID)
		c.SetUserContext(ctx)
		return c.Next()
	}
}
