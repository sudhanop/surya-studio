package handlers

import (
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/suryaphotography/backend/internal/config"
	"github.com/suryaphotography/backend/internal/models"
	"github.com/suryaphotography/backend/internal/repositories"
	"github.com/suryaphotography/backend/internal/storage"
)

type SettingsHandler struct {
	cfg   *config.Config
	repo  *repositories.SettingsRepo
	store *storage.LocalStorage
}

func NewSettingsHandler(cfg *config.Config, repo *repositories.SettingsRepo, store *storage.LocalStorage) *SettingsHandler {
	return &SettingsHandler{cfg: cfg, repo: repo, store: store}
}

func defaultTestimonials() []models.Testimonial {
	return []models.Testimonial{
		{Name: "Priya & Arun", Text: "Every frame felt like a movie. Surya Photography captured our wedding with pure magic."},
		{Name: "Lakshmi Family", Text: "Professional, warm, and incredibly talented. Our puberty function album is breathtaking."},
		{Name: "Divya", Text: "The maternity shoot was dreamy. Highly recommend for anyone who wants cinematic quality."},
	}
}

func (h *SettingsHandler) loadMap(c *fiber.Ctx) map[string]string {
	m, err := h.repo.GetAll(c.Context())
	if err != nil || m == nil {
		return map[string]string{}
	}
	return m
}

func (h *SettingsHandler) parseTestimonials(raw string) []models.Testimonial {
	if strings.TrimSpace(raw) == "" {
		return defaultTestimonials()
	}
	var list []models.Testimonial
	if err := json.Unmarshal([]byte(raw), &list); err != nil || len(list) == 0 {
		return defaultTestimonials()
	}
	return list
}

func (h *SettingsHandler) buildPublic(m map[string]string) models.PublicSiteData {
	get := func(k, fallback string) string {
		if v := strings.TrimSpace(m[k]); v != "" {
			return v
		}
		return fallback
	}
	portraitURL := ""
	if p := m["owner_portrait_path"]; p != "" {
		portraitURL = h.store.PublicURL(p)
	}
	logoURL := ""
	if p := m["site_logo_path"]; p != "" {
		logoURL = h.store.PublicURL(p)
	}
	return models.PublicSiteData{
		EventsCovered:    get("events_covered", "500+"),
		YearsOfCraft:     get("years_of_craft", "10+"),
		HappyFamilies:    get("happy_families", "1000+"),
		OwnerPortraitURL: portraitURL,
		LogoURL:          logoURL,
		Testimonials:     h.parseTestimonials(m["testimonials_json"]),
		Contact: models.ContactInfo{
			WhatsApp:        get("whatsapp", h.cfg.WhatsApp),
			InstagramURL:    get("instagram_url", h.cfg.InstagramURL),
			FacebookURL:     get("facebook_url", h.cfg.FacebookURL),
			YouTubeURL:      get("youtube_url", h.cfg.YouTubeURL),
			ContactEmail:    get("contact_email", h.cfg.ContactEmail),
			PhoneNumber:     get("phone_primary", h.cfg.PhoneNumber),
			PhoneSecondary:  get("phone_secondary", h.cfg.PhoneSecondary),
			Address:         get("address", h.cfg.Address),
			Pincode:         get("pincode", h.cfg.Pincode),
			GoogleMapsEmbed: h.cfg.GoogleMapsEmbed,
		},
	}
}

func (h *SettingsHandler) PublicContact(c *fiber.Ctx) error {
	return c.JSON(h.buildPublic(h.loadMap(c)).Contact)
}

func (h *SettingsHandler) PublicSite(c *fiber.Ctx) error {
	return c.JSON(h.buildPublic(h.loadMap(c)))
}

func (h *SettingsHandler) AdminGet(c *fiber.Ctx) error {
	return c.JSON(h.buildPublic(h.loadMap(c)))
}

type settingsUpdateRequest struct {
	EventsCovered  string               `json:"events_covered"`
	YearsOfCraft   string               `json:"years_of_craft"`
	HappyFamilies  string               `json:"happy_families"`
	ContactEmail   string               `json:"contact_email"`
	PhonePrimary   string               `json:"phone_primary"`
	PhoneSecondary string               `json:"phone_secondary"`
	InstagramURL   string               `json:"instagram_url"`
	YouTubeURL     string               `json:"youtube_url"`
	Address        string               `json:"address"`
	Pincode        string               `json:"pincode"`
	WhatsApp       string               `json:"whatsapp"`
	Testimonials   []models.Testimonial `json:"testimonials"`
}

func (h *SettingsHandler) AdminUpdate(c *fiber.Ctx) error {
	var req settingsUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}
	pairs := map[string]string{
		"events_covered":  req.EventsCovered,
		"years_of_craft":  req.YearsOfCraft,
		"happy_families":  req.HappyFamilies,
		"contact_email":   req.ContactEmail,
		"phone_primary":   req.PhonePrimary,
		"phone_secondary": req.PhoneSecondary,
		"instagram_url":   req.InstagramURL,
		"youtube_url":     req.YouTubeURL,
		"address":         req.Address,
		"pincode":         req.Pincode,
		"whatsapp":        req.WhatsApp,
	}
	if len(req.Testimonials) > 0 {
		b, _ := json.Marshal(req.Testimonials)
		pairs["testimonials_json"] = string(b)
	}
	if err := h.repo.SetMany(c.Context(), pairs); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "save failed"})
	}
	return c.JSON(h.buildPublic(h.loadMap(c)))
}

func (h *SettingsHandler) uploadSiteImage(c *fiber.Ctx, settingKey, prefix string) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "file required"})
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" && ext != ".svg" {
		return c.Status(400).JSON(fiber.Map{"error": "image only"})
	}
	rel := "site/" + prefix + "-" + uuid.New().String() + ext
	fh, err := file.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "cannot open file"})
	}
	defer fh.Close()
	if err := h.store.SaveRelative(rel, fh); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "upload failed"})
	}
	m := h.loadMap(c)
	if old := m[settingKey]; old != "" {
		_ = h.store.Delete(old)
	}
	if err := h.repo.Set(c.Context(), settingKey, rel); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "save failed"})
	}
	return c.JSON(fiber.Map{"url": h.store.PublicURL(rel)})
}

func (h *SettingsHandler) UploadOwnerPortrait(c *fiber.Ctx) error {
	return h.uploadSiteImage(c, "owner_portrait_path", "owner")
}

func (h *SettingsHandler) DeleteOwnerPortrait(c *fiber.Ctx) error {
	return h.deleteSiteImage(c, "owner_portrait_path")
}

func (h *SettingsHandler) UploadLogo(c *fiber.Ctx) error {
	if err := h.uploadSiteImage(c, "site_logo_path", "logo"); err != nil {
		return err
	}
	m := h.loadMap(c)
	return c.JSON(fiber.Map{"logo_url": h.store.PublicURL(m["site_logo_path"])})
}

func (h *SettingsHandler) DeleteLogo(c *fiber.Ctx) error {
	if err := h.deleteSiteImage(c, "site_logo_path"); err != nil {
		return err
	}
	return c.JSON(fiber.Map{"message": "removed"})
}

func (h *SettingsHandler) deleteSiteImage(c *fiber.Ctx, settingKey string) error {
	m := h.loadMap(c)
	if p := m[settingKey]; p != "" {
		_ = h.store.Delete(p)
	}
	_ = h.repo.Set(c.Context(), settingKey, "")
	return c.JSON(fiber.Map{"message": "removed"})
}
