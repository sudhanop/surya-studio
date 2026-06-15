package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port              string
	AppEnv            string
	FrontendURL       string
	CORSOrigins       string
	JWTSecret         string
	JWTExpiryHours    int
	DataDir           string
	UploadDir         string
	MaxUploadMB       int64
	PublicMediaURL    string
	SMTPHost          string
	SMTPPort          int
	SMTPUser          string
	SMTPPassword      string
	SMTPFrom          string
	AdminEmail        string
	WhatsApp          string
	InstagramURL      string
	FacebookURL       string
	YouTubeURL        string
	ContactEmail    string
	PhoneNumber     string
	PhoneSecondary  string
	Address         string
	Pincode         string
	GoogleMapsEmbed string
}

func Load() (*Config, error) {
	// Overload so backend/.env wins over stale machine-level DB_* variables
	_ = godotenv.Overload()

	expiry, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))
	smtpPort, _ := strconv.Atoi(getEnv("SMTP_PORT", "587"))
	maxMB, _ := strconv.ParseInt(getEnv("MAX_UPLOAD_MB", "100"), 10, 64)

	return &Config{
		Port:            getEnv("PORT", "8080"),
		AppEnv:          getEnv("APP_ENV", "development"),
		FrontendURL:     getEnv("FRONTEND_URL", "http://localhost:3000"),
		CORSOrigins:     getEnv("CORS_ORIGINS", getEnv("FRONTEND_URL", "http://localhost:3000")),
		JWTSecret:       getEnv("JWT_SECRET", "dev-secret-change-in-production"),
		JWTExpiryHours:  expiry,
		DataDir:         getEnv("DATA_DIR", "data"),
		UploadDir:       getEnv("UPLOAD_DIR", "../uploads"),
		MaxUploadMB:     maxMB,
		PublicMediaURL:  strings.TrimRight(getEnv("PUBLIC_MEDIA_URL", "http://localhost:8080/uploads"), "/"),
		SMTPHost:        os.Getenv("SMTP_HOST"),
		SMTPPort:        smtpPort,
		SMTPUser:        os.Getenv("SMTP_USER"),
		SMTPPassword:    os.Getenv("SMTP_PASSWORD"),
		SMTPFrom:        getEnv("SMTP_FROM", "noreply@suryaphotography.com"),
		AdminEmail:      getEnv("ADMIN_EMAIL", "admin@suryaphotography.com"),
		WhatsApp:        getEnv("WHATSAPP_NUMBER", "919715241568"),
		InstagramURL:    getEnv("INSTAGRAM_URL", "https://www.instagram.com/surya_photography_nkl"),
		FacebookURL:     getEnv("FACEBOOK_URL", ""),
		YouTubeURL:      getEnv("YOUTUBE_URL", "https://www.youtube.com/@suryaphotography4303"),
		ContactEmail:    getEnv("CONTACT_EMAIL", "suryaphotographyrsp@gmail.com"),
		PhoneNumber:     getEnv("PHONE_NUMBER", "9715241568"),
		PhoneSecondary:  getEnv("PHONE_SECONDARY", "8884897499"),
		Address:         getEnv("STUDIO_ADDRESS", "Surya Photography, near DNC (Chamundi) Theater, opposite to Adhisindha Thirumana Mandabam, Pattanam road, Rasipuram"),
		Pincode:         getEnv("STUDIO_PINCODE", "637408"),
		GoogleMapsEmbed: getEnv("GOOGLE_MAPS_EMBED", ""),
	}, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
