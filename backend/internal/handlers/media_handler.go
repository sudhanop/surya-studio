package handlers

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/suryaphotography/backend/internal/config"
	"github.com/suryaphotography/backend/internal/jsonutil"
	"github.com/suryaphotography/backend/internal/models"
	"github.com/suryaphotography/backend/internal/repositories"
	"github.com/suryaphotography/backend/internal/storage"
)

type MediaHandler struct {
	cfg      *config.Config
	repo     *repositories.MediaRepo
	catRepo  *repositories.CategoryRepo
	store    *storage.LocalStorage
}

func NewMediaHandler(cfg *config.Config, repo *repositories.MediaRepo, catRepo *repositories.CategoryRepo, store *storage.LocalStorage) *MediaHandler {
	return &MediaHandler{cfg: cfg, repo: repo, catRepo: catRepo, store: store}
}

func (h *MediaHandler) enrich(m *models.PortfolioMedia) {
	m.URL = h.store.PublicURL(m.FilePath)
	if m.ThumbnailPath != nil {
		u := h.store.PublicURL(*m.ThumbnailPath)
		m.ThumbnailURL = u
	} else if m.MediaType == "photo" {
		m.ThumbnailURL = m.URL
	}
}

func (h *MediaHandler) enrichList(list []models.PortfolioMedia) []models.PortfolioMedia {
	list = jsonutil.Slice(list)
	for i := range list {
		h.enrich(&list[i])
	}
	return list
}

// ListByCategory godoc
// @Summary List media for category
// @Tags portfolio
// @Param slug path string true "Category slug"
// @Param type query string false "photo or video"
// @Success 200 {object} map[string]interface{}
// @Router /api/categories/{slug}/media [get]
func (h *MediaHandler) ListByCategory(c *fiber.Ctx) error {
	cat, err := h.catRepo.GetBySlug(c.Context(), c.Params("slug"))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "category not found"})
	}
	mediaType := c.Query("type")
	list, err := h.repo.ListByCategory(c.Context(), cat.ID, mediaType, true)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to load media"})
	}
	list = h.enrichList(list)
	photos, videos := []models.PortfolioMedia{}, []models.PortfolioMedia{}
	for _, m := range list {
		if m.MediaType == "video" {
			videos = append(videos, m)
		} else {
			photos = append(photos, m)
		}
	}
	return c.JSON(fiber.Map{
		"category": cat,
		"photos":   photos,
		"videos":   videos,
		"has_videos": len(videos) > 0,
	})
}

// Featured godoc
// @Summary Featured portfolio media
// @Tags portfolio
// @Success 200 {array} models.PortfolioMedia
// @Router /api/portfolio/featured [get]
func (h *MediaHandler) Featured(c *fiber.Ctx) error {
	list, err := h.repo.ListFeatured(c.Context(), 12)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed"})
	}
	return c.JSON(h.enrichList(list))
}

// Latest godoc
// @Summary Latest portfolio photos
// @Tags portfolio
// @Success 200 {array} models.PortfolioMedia
// @Router /api/portfolio/latest [get]
func (h *MediaHandler) Latest(c *fiber.Ctx) error {
	list, err := h.repo.ListLatest(c.Context(), 8)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed"})
	}
	return c.JSON(h.enrichList(list))
}

// ListAdmin godoc
// @Summary List media (admin)
// @Tags admin-media
// @Security BearerAuth
// @Param category_id query int false "Filter by category"
// @Router /api/admin/media [get]
func (h *MediaHandler) ListAdmin(c *fiber.Ctx) error {
	catID := c.QueryInt("category_id", 0)
	if catID == 0 {
		// all categories - get featured/latest style full list per category 0 means all - simplified: return latest 100
		list, err := h.repo.ListLatest(c.Context(), 100)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed"})
		}
		return c.JSON(h.enrichList(list))
	}
	list, err := h.repo.ListByCategory(c.Context(), int64(catID), "", false)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed"})
	}
	return c.JSON(h.enrichList(list))
}

type updateMediaRequest struct {
	Title        *string `json:"title"`
	Caption      *string `json:"caption"`
	IsFeatured   bool    `json:"is_featured"`
	DisplayOrder int     `json:"display_order"`
	IsPublished  bool    `json:"is_published"`
}

// Update godoc
// @Summary Update media metadata
// @Tags admin-media
// @Security BearerAuth
// @Param id path int true "Media ID"
// @Router /api/admin/media/{id} [put]
func (h *MediaHandler) Update(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	m, err := h.repo.GetByID(c.Context(), int64(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	var req updateMediaRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}
	m.Title = req.Title
	m.Caption = req.Caption
	m.IsFeatured = req.IsFeatured
	m.DisplayOrder = req.DisplayOrder
	m.IsPublished = req.IsPublished
	if err := h.repo.Update(c.Context(), m); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "update failed"})
	}
	h.enrich(m)
	return c.JSON(m)
}

// Delete godoc
// @Summary Delete media
// @Tags admin-media
// @Security BearerAuth
// @Param id path int true "Media ID"
// @Router /api/admin/media/{id} [delete]
func (h *MediaHandler) Delete(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	m, err := h.repo.Delete(c.Context(), int64(id))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "delete failed"})
	}
	_ = h.store.Delete(m.FilePath)
	if m.ThumbnailPath != nil {
		_ = h.store.Delete(*m.ThumbnailPath)
	}
	return c.JSON(fiber.Map{"message": "deleted"})
}

// Upload godoc
// @Summary Upload photo or video
// @Tags admin-media
// @Security BearerAuth
// @Accept multipart/form-data
// @Param category_id formData int true "Category ID"
// @Param media_type formData string true "photo or video"
// @Param file formData file true "Media file"
// @Param title formData string false "Title"
// @Param is_featured formData bool false "Featured"
// @Router /api/admin/upload [post]
func (h *MediaHandler) Upload(c *fiber.Ctx) error {
	var categoryID int64
	if v := c.FormValue("category_id"); v != "" {
		fmt.Sscanf(v, "%d", &categoryID)
	}
	if categoryID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "category_id required"})
	}

	cat, err := h.catRepo.GetByID(c.Context(), categoryID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "category not found"})
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "file required"})
	}
	if file.Size > h.cfg.MaxUploadMB*1024*1024 {
		return c.Status(400).JSON(fiber.Map{"error": "file too large"})
	}

	mediaType := c.FormValue("media_type", "photo")
	if mediaType != "photo" && mediaType != "video" {
		mediaType = "photo"
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	filename := uuid.New().String() + ext

	fh, err := file.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "cannot open file"})
	}
	defer fh.Close()

	var relPath string
	if mediaType == "video" {
		relPath, err = h.store.SaveVideo(cat.Slug, filename, fh)
	} else {
		relPath, err = h.store.Save(cat.Slug, filename, fh)
	}
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "upload failed"})
	}

	title := c.FormValue("title")
	var titlePtr *string
	if title != "" {
		titlePtr = &title
	}
	isFeatured := c.FormValue("is_featured") == "true" || c.FormValue("is_featured") == "1"

	m := &models.PortfolioMedia{
		CategoryID:   categoryID,
		Title:        titlePtr,
		MediaType:    mediaType,
		FilePath:     relPath,
		MimeType:     strPtr(file.Header.Get("Content-Type")),
		FileSizeBytes: ptrInt64(file.Size),
		IsFeatured:   isFeatured,
		IsPublished:  true,
	}
	id, err := h.repo.Create(c.Context(), m)
	if err != nil {
		_ = h.store.Delete(relPath)
		return c.Status(500).JSON(fiber.Map{"error": "save failed"})
	}
	m.ID = id
	h.enrich(m)
	return c.Status(201).JSON(m)
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func ptrInt64(n int64) *int64 { return &n }
