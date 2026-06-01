package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/suryaphotography/backend/internal/config"
	"github.com/suryaphotography/backend/internal/services"
)

type DashboardHandler struct {
	cfg *config.Config
	svc *services.DashboardService
}

func NewDashboardHandler(cfg *config.Config, svc *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{cfg: cfg, svc: svc}
}

// Stats godoc
// @Summary Dashboard statistics
// @Tags admin-dashboard
// @Security BearerAuth
// @Success 200 {object} models.DashboardStats
// @Router /api/admin/dashboard [get]
func (h *DashboardHandler) Stats(c *fiber.Ctx) error {
	stats, err := h.svc.Stats(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to load dashboard"})
	}
	return c.JSON(stats)
}

// ContactInfo godoc
// @Summary Public contact information
// @Tags public
// @Success 200 {object} models.ContactInfo
// @Router /api/contact-info [get]
func ContactInfo(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"whatsapp":          cfg.WhatsApp,
			"instagram_url":     cfg.InstagramURL,
			"facebook_url":      cfg.FacebookURL,
			"youtube_url":       cfg.YouTubeURL,
			"contact_email":     cfg.ContactEmail,
			"phone_number":      cfg.PhoneNumber,
			"google_maps_embed": cfg.GoogleMapsEmbed,
		})
	}
}
