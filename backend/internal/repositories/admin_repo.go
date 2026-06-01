package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/suryaphotography/backend/internal/models"
)

type AdminRepo struct{ db *sql.DB }

func NewAdminRepo(db *sql.DB) *AdminRepo { return &AdminRepo{db: db} }

func (r *AdminRepo) GetByUsername(ctx context.Context, username string) (*models.Admin, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, username, password_hash, email, display_name, is_active, last_login_at
		FROM dbo.admin WHERE username = @p1 AND is_active = 1`, username)

	var a models.Admin
	var display sql.NullString
	var lastLogin sql.NullTime
	err := row.Scan(&a.ID, &a.Username, &a.PasswordHash, &a.Email, &display, &a.IsActive, &lastLogin)
	if err != nil {
		return nil, err
	}
	if display.Valid {
		a.DisplayName = &display.String
	}
	if lastLogin.Valid {
		a.LastLoginAt = &lastLogin.Time
	}
	return &a, nil
}

func (r *AdminRepo) UpdateLastLogin(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `UPDATE dbo.admin SET last_login_at = @p1, updated_at = @p1 WHERE id = @p2`, time.Now().UTC(), id)
	return err
}
