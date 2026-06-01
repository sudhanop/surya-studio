package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	fiberSwagger "github.com/gofiber/swagger"

	"github.com/suryaphotography/backend/internal/config"
	"github.com/suryaphotography/backend/internal/handlers"
	"github.com/suryaphotography/backend/internal/middleware"
	"github.com/suryaphotography/backend/internal/storage"
)

type Handlers struct {
	Auth      *handlers.AuthHandler
	Category  *handlers.CategoryHandler
	Media     *handlers.MediaHandler
	Inquiry   *handlers.InquiryHandler
	Function  *handlers.FunctionHandler
	Dashboard *handlers.DashboardHandler
	Settings  *handlers.SettingsHandler
}

func Setup(app *fiber.App, cfg *config.Config, h Handlers, store *storage.LocalStorage) {
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(middleware.SecurityHeaders())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSOrigins,
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
	}))

	app.Static("/uploads", store.BaseDir())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "service": "surya-photography-api"})
	})

	app.Get("/api/contact-info", h.Settings.PublicContact)
	app.Get("/api/site", h.Settings.PublicSite)

	// Public API
	api := app.Group("/api")
	api.Get("/categories", h.Category.ListPublic)
	api.Get("/categories/:slug", h.Category.GetBySlug)
	api.Get("/categories/:slug/media", h.Media.ListByCategory)
	api.Get("/portfolio/featured", h.Media.Featured)
	api.Get("/portfolio/latest", h.Media.Latest)
	api.Post("/inquiries", h.Inquiry.Create)

	// Admin auth
	admin := api.Group("/admin")
	admin.Post("/login", middleware.LoginRateLimit(), h.Auth.Login)
	admin.Post("/logout", h.Auth.Logout)

	protected := admin.Group("", middleware.Protected(cfg))
	protected.Get("/me", h.Auth.Me)
	protected.Get("/dashboard", h.Dashboard.Stats)

	protected.Get("/categories", h.Category.ListAdmin)
	protected.Put("/categories/:id", h.Category.Update)

	protected.Get("/media", h.Media.ListAdmin)
	protected.Put("/media/:id", h.Media.Update)
	protected.Delete("/media/:id", h.Media.Delete)
	protected.Post("/upload", h.Media.Upload)

	protected.Get("/inquiries", h.Inquiry.List)
	protected.Put("/inquiries/:id/status", h.Inquiry.UpdateStatus)
	protected.Delete("/inquiries/:id", h.Inquiry.Delete)
	protected.Post("/inquiries/:id/convert", h.Inquiry.Convert)

	protected.Get("/functions", h.Function.List)
	protected.Get("/functions/upcoming", h.Function.Upcoming)
	protected.Post("/functions", h.Function.Create)
	protected.Put("/functions/:id", h.Function.Update)
	protected.Delete("/functions/:id", h.Function.Delete)

	protected.Get("/settings", h.Settings.AdminGet)
	protected.Put("/settings", h.Settings.AdminUpdate)
	protected.Post("/settings/owner-portrait", h.Settings.UploadOwnerPortrait)
	protected.Delete("/settings/owner-portrait", h.Settings.DeleteOwnerPortrait)
	protected.Post("/settings/logo", h.Settings.UploadLogo)
	protected.Delete("/settings/logo", h.Settings.DeleteLogo)

	app.Get("/swagger/*", fiberSwagger.HandlerDefault)
}
