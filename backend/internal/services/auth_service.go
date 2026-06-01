package services

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/suryaphotography/backend/internal/config"
	"github.com/suryaphotography/backend/internal/middleware"
	"github.com/suryaphotography/backend/internal/models"
	"github.com/suryaphotography/backend/internal/repositories"
)

type AuthService struct {
	cfg  *config.Config
	repo *repositories.AdminRepo
}

func NewAuthService(cfg *config.Config, repo *repositories.AdminRepo) *AuthService {
	return &AuthService{cfg: cfg, repo: repo}
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, *models.Admin, error) {
	admin, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}
	_ = s.repo.UpdateLastLogin(ctx, admin.ID)
	token, err := s.issueToken(admin)
	if err != nil {
		return "", nil, err
	}
	admin.PasswordHash = ""
	return token, admin, nil
}

func (s *AuthService) issueToken(admin *models.Admin) (string, error) {
	exp := time.Now().Add(time.Duration(s.cfg.JWTExpiryHours) * time.Hour)
	claims := &middleware.AdminClaims{
		AdminID:  admin.ID,
		Username: admin.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(s.cfg.JWTSecret))
}
