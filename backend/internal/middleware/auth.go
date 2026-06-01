package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"github.com/suryaphotography/backend/internal/config"
)

const AdminClaimsKey = "admin"

type AdminClaims struct {
	AdminID  int64  `json:"admin_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func Protected(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenStr := extractToken(c, cfg)
		if tokenStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}
		claims := &AdminClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}
		c.Locals(AdminClaimsKey, claims)
		return c.Next()
	}
}

func extractToken(c *fiber.Ctx, cfg *config.Config) string {
	if auth := c.Get("Authorization"); strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}
	return c.Cookies("surya_token")
}

func GetClaims(c *fiber.Ctx) *AdminClaims {
	if v := c.Locals(AdminClaimsKey); v != nil {
		if claims, ok := v.(*AdminClaims); ok {
			return claims
		}
	}
	return nil
}
