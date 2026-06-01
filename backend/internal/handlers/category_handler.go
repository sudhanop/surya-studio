package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/suryaphotography/backend/internal/config"
	"github.com/suryaphotography/backend/internal/jsonutil"
	"github.com/suryaphotography/backend/internal/models"
	"github.com/suryaphotography/backend/internal/repositories"
	"github.com/suryaphotography/backend/internal/storage"
)

type CategoryHandler struct {
	cfg  *config.Config
	repo *repositories.CategoryRepo
	store storage.Storage
}

func NewCategoryHandler(cfg *config.Config, repo *repositories.CategoryRepo, store storage.Storage) *CategoryHandler {
	return &CategoryHandler{cfg: cfg, repo: repo, store: store}
}

func (h *CategoryHandler) enrich(categories []models.Category) {
	for i := range categories {
		if categories[i].CoverImage != nil && *categories[i].CoverImage != "" {
			url := h.store.PublicURL(*categories[i].CoverImage)
			categories[i].CoverImage = &url
		}
	}
}

// ListPublic godoc
// @Summary List active categories
// @Tags categories
// @Produce json
// @Success 200 {array} models.Category
// @Router /api/categories [get]
func (h *CategoryHandler) ListPublic(c *fiber.Ctx) error {
	list, err := h.repo.List(c.Context(), true)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to load categories"})
	}
	h.enrich(list)
	return c.JSON(jsonutil.Slice(list))
}

// GetBySlug godoc
// @Summary Get category by slug
// @Tags categories
// @Param slug path string true "Category slug"
// @Success 200 {object} models.Category
// @Router /api/categories/{slug} [get]
func (h *CategoryHandler) GetBySlug(c *fiber.Ctx) error {
	cat, err := h.repo.GetBySlug(c.Context(), c.Params("slug"))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "category not found"})
	}
	h.enrich([]models.Category{*cat})
	return c.JSON(cat)
}

// ListAdmin godoc
// @Summary List all categories (admin)
// @Tags admin-categories
// @Security BearerAuth
// @Success 200 {array} models.Category
// @Router /api/admin/categories [get]
func (h *CategoryHandler) ListAdmin(c *fiber.Ctx) error {
	list, err := h.repo.List(c.Context(), false)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to load categories"})
	}
	h.enrich(list)
	return c.JSON(jsonutil.Slice(list))
}

type updateCategoryRequest struct {
	Name         string  `json:"name"`
	Description  *string `json:"description"`
	CoverImage   *string `json:"cover_image"`
	DisplayOrder int     `json:"display_order"`
	IsActive     bool    `json:"is_active"`
}

// Update godoc
// @Summary Update category
// @Tags admin-categories
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Param body body updateCategoryRequest true "Category"
// @Success 200 {object} models.MessageResponse
// @Router /api/admin/categories/{id} [put]
func (h *CategoryHandler) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	cat, err := h.repo.GetByID(c.Context(), int64(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	var req updateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}
	if req.Name != "" {
		cat.Name = req.Name
	}
	cat.Description = req.Description
	cat.CoverImage = req.CoverImage
	cat.DisplayOrder = req.DisplayOrder
	cat.IsActive = req.IsActive
	if err := h.repo.Update(c.Context(), cat); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "update failed"})
	}
	return c.JSON(fiber.Map{"message": "updated"})
}
