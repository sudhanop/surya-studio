package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"

	"github.com/suryaphotography/backend/internal/config"
	"github.com/suryaphotography/backend/internal/database"
	"github.com/suryaphotography/backend/internal/email"
	"github.com/suryaphotography/backend/internal/handlers"
	"github.com/suryaphotography/backend/internal/repositories"
	"github.com/suryaphotography/backend/internal/routes"
	"github.com/suryaphotography/backend/internal/services"
	"github.com/suryaphotography/backend/internal/storage"

	_ "github.com/suryaphotography/backend/docs"
)

// @title Surya Photography API
// @version 1.0
// @description Portfolio and studio management API
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	uploadDir, _ := filepath.Abs(cfg.UploadDir)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Fatal(err)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()

	store, err := storage.NewLocalStorage(uploadDir, cfg.PublicMediaURL)
	if err != nil {
		log.Fatal(err)
	}

	adminRepo := repositories.NewAdminRepo(db)
	catRepo := repositories.NewCategoryRepo(db)
	mediaRepo := repositories.NewMediaRepo(db)
	inquiryRepo := repositories.NewInquiryRepo(db)
	funcRepo := repositories.NewFunctionRepo(db)
	settingsRepo := repositories.NewSettingsRepo(db)

	authSvc := services.NewAuthService(cfg, adminRepo)
	dashSvc := services.NewDashboardService(funcRepo, inquiryRepo, mediaRepo)
	mailSvc := email.New(cfg)

	h := routes.Handlers{
		Auth:      handlers.NewAuthHandler(cfg, authSvc),
		Category:  handlers.NewCategoryHandler(cfg, catRepo, store),
		Media:     handlers.NewMediaHandler(cfg, mediaRepo, catRepo, store),
		Inquiry:   handlers.NewInquiryHandler(inquiryRepo, funcRepo, mailSvc),
		Function:  handlers.NewFunctionHandler(funcRepo),
		Dashboard: handlers.NewDashboardHandler(cfg, dashSvc),
		Settings:  handlers.NewSettingsHandler(cfg, settingsRepo, store),
	}

	app := fiber.New(fiber.Config{
		BodyLimit: int(cfg.MaxUploadMB) * 1024 * 1024,
	})
	routes.Setup(app, cfg, h, store)

	log.Printf("Surya Photography API on :%s", cfg.Port)
	log.Printf("Swagger: http://localhost:%s/swagger/index.html", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
