package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/suryaphotography/backend/internal/config"
	"github.com/suryaphotography/backend/internal/middleware"
	"github.com/suryaphotography/backend/internal/services"
)

type AuthHandler struct {
	cfg *config.Config
	svc *services.AuthService
}

func NewAuthHandler(cfg *config.Config, svc *services.AuthService) *AuthHandler {
	return &AuthHandler{cfg: cfg, svc: svc}
}

type loginRequest struct {
	Username string `json:"username" example:"surya@admin.com"`
	Password string `json:"password" example:"surya@1995"`
}

type loginResponse struct {
	Token string      `json:"token"`
	Admin interface{} `json:"admin"`
}

// Login godoc
// @Summary Admin login
// @Tags auth
// @Accept json
// @Produce json
// @Param body body loginRequest true "Credentials"
// @Success 200 {object} loginResponse
// @Failure 401 {object} map[string]string
// @Router /api/admin/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req loginRequest
	if err := c.BodyParser(&req); err != nil || req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "username and password required"})
	}
	token, admin, err := h.svc.Login(c.Context(), req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}
	c.Cookie(&fiber.Cookie{
		Name:     "surya_token",
		Value:    token,
		HTTPOnly: true,
		Secure:   h.cfg.AppEnv == "production",
		SameSite: "Lax",
		Path:     "/",
		MaxAge:   h.cfg.JWTExpiryHours * 3600,
	})
	return c.JSON(loginResponse{Token: token, Admin: admin})
}

// Logout godoc
// @Summary Admin logout
// @Tags auth
// @Success 200 {object} map[string]string
// @Router /api/admin/logout [post]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{Name: "surya_token", Value: "", MaxAge: -1, Path: "/"})
	return c.JSON(fiber.Map{"message": "logged out"})
}

// Me godoc
// @Summary Current admin
// @Tags auth
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/admin/me [get]
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	claims := middleware.GetClaims(c)
	if claims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	return c.JSON(fiber.Map{"admin_id": claims.AdminID, "username": claims.Username})
}
