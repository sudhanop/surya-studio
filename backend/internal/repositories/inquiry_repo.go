package repositories

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/suryaphotography/backend/internal/models"
)

type InquiryRepo struct {
	file string
	mu   sync.Mutex
}

func NewInquiryRepo(dataDir string) *InquiryRepo {
	return &InquiryRepo{file: filepath.Join(dataDir, "inquiries.json")}
}

func (r *InquiryRepo) load() ([]models.Inquiry, error) {
	var list []models.Inquiry
	err := readJSONFile(r.file, &list)
	if err != nil {
		if os.IsNotExist(err) {
			return []models.Inquiry{}, nil
		}
		return nil, err
	}
	if list == nil {
		list = []models.Inquiry{}
	}
	return list, nil
}

func (r *InquiryRepo) save(list []models.Inquiry) error {
	return writeJSONFile(r.file, list)
}

func (r *InquiryRepo) Create(ctx context.Context, i *models.Inquiry) (int64, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	list, err := r.load()
	if err != nil {
		return 0, err
	}
	var maxID int64
	for _, it := range list {
		if it.ID > maxID {
			maxID = it.ID
		}
	}
	now := time.Now().UTC()
	n := *i
	n.ID = maxID + 1
	if n.Status == "" {
		n.Status = "new"
	}
	n.CreatedAt = now
	n.UpdatedAt = now
	list = append(list, n)
	if err := r.save(list); err != nil {
		return 0, err
	}
	return n.ID, nil
}

func (r *InquiryRepo) List(ctx context.Context, status string, limit int) ([]models.Inquiry, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	list, err := r.load()
	if err != nil {
		return nil, err
	}
	out := make([]models.Inquiry, 0, len(list))
	for _, it := range list {
		if status == "" || it.Status == status {
			out = append(out, it)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out, nil
}

func (r *InquiryRepo) GetByID(ctx context.Context, id int64) (*models.Inquiry, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	list, err := r.load()
	if err != nil {
		return nil, err
	}
	for _, it := range list {
		if it.ID == id {
			cp := it
			return &cp, nil
		}
	}
	return nil, ErrNotFound
}

func (r *InquiryRepo) UpdateStatus(ctx context.Context, id int64, status string) error {
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
			list[i].Status = status
			if status == "contacted" && list[i].ContactedAt == nil {
				list[i].ContactedAt = &now
			}
			list[i].UpdatedAt = now
			return r.save(list)
		}
	}
	return ErrNotFound
}

func (r *InquiryRepo) Delete(ctx context.Context, id int64) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	list, err := r.load()
	if err != nil {
		return err
	}
	out := list[:0]
	found := false
	for _, it := range list {
		if it.ID == id {
			found = true
			continue
		}
		out = append(out, it)
	}
	if !found {
		return ErrNotFound
	}
	return r.save(out)
}

func (r *InquiryRepo) CountRecent(ctx context.Context, days int) (int, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	list, err := r.load()
	if err != nil {
		return 0, err
	}
	since := time.Now().UTC().AddDate(0, 0, -days)
	n := 0
	for _, it := range list {
		if it.Status == "new" && (it.CreatedAt.After(since) || it.CreatedAt.Equal(since)) {
			n++
		}
	}
	return n, nil
}
