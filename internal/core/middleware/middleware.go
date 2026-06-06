package core_middleware

import (
	core_jwt "api/internal/core/jwt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

var userID = "user_id"

func JWTCheck() fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(http.StatusUnauthorized).JSON("unauthorized")
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.ParseWithClaims(tokenStr, &core_jwt.Claims{}, func(token *jwt.Token) (any, error) {
			return core_jwt.JwtSecret, nil
		})

		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{
				"message": "invalid token",
			})
		}

		claims := token.Claims.(*core_jwt.Claims)

		c.Locals(userID, claims.UserID)

		return c.Next()
	}
}

