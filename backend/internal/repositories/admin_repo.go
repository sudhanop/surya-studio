package repositories

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/suryaphotography/backend/internal/models"
)

type AdminRepo struct {
	file string
	mu   sync.Mutex
}

type adminRecord struct {
	ID           int64      `json:"id"`
	Username     string     `json:"username"`
	PasswordHash string     `json:"password_hash"`
	Email        string     `json:"email"`
	DisplayName  *string    `json:"display_name,omitempty"`
	IsActive     bool       `json:"is_active"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
}

func NewAdminRepo(dataDir string) *AdminRepo {
	return &AdminRepo{file: filepath.Join(dataDir, "admins.json")}
}

func (r *AdminRepo) load() ([]adminRecord, error) {
	var list []adminRecord
	err := readJSONFile(r.file, &list)
	if err != nil {
		if os.IsNotExist(err) {
			return []adminRecord{}, nil
		}
		return nil, err
	}
	if list == nil {
		list = []adminRecord{}
	}
	return list, nil
}

func (r *AdminRepo) save(list []adminRecord) error {
	return writeJSONFile(r.file, list)
}

func (r *AdminRepo) GetByUsername(ctx context.Context, username string) (*models.Admin, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	list, err := r.load()
	if err != nil {
		return nil, err
	}
	for _, rec := range list {
		if rec.IsActive && rec.Username == username {
			return &models.Admin{
				ID:           rec.ID,
				Username:     rec.Username,
				PasswordHash: rec.PasswordHash,
				Email:        rec.Email,
				DisplayName:  rec.DisplayName,
				IsActive:     rec.IsActive,
				LastLoginAt:  rec.LastLoginAt,
			}, nil
		}
	}
	return nil, ErrNotFound
}

func (r *AdminRepo) UpdateLastLogin(ctx context.Context, id int64) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	list, err := r.load()
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	for i := range list {
		if list[i].ID == id {
			list[i].LastLoginAt = &now
			return r.save(list)
		}
	}
	return ErrNotFound
}
