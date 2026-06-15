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

type MediaRepo struct {
	file           string
	categoriesFile string
	mu             sync.Mutex
}

func NewMediaRepo(dataDir string) *MediaRepo {
	return &MediaRepo{
		file:           filepath.Join(dataDir, "media.json"),
		categoriesFile: filepath.Join(dataDir, "categories.json"),
	}
}

func (r *MediaRepo) load() ([]models.PortfolioMedia, error) {
	var list []models.PortfolioMedia
	err := readJSONFile(r.file, &list)
	if err != nil {
		if os.IsNotExist(err) {
			return []models.PortfolioMedia{}, nil
		}
		return nil, err
	}
	if list == nil {
		list = []models.PortfolioMedia{}
	}
	return list, nil
}

func (r *MediaRepo) save(list []models.PortfolioMedia) error {
	return writeJSONFile(r.file, list)
}

func (r *MediaRepo) ListByCategory(ctx context.Context, categoryID int64, mediaType string, publishedOnly bool) ([]models.PortfolioMedia, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	list, err := r.load()
	if err != nil {
		return nil, err
	}
	out := make([]models.PortfolioMedia, 0, len(list))
	for _, it := range list {
		if it.CategoryID != categoryID {
			continue
		}
		if mediaType != "" && it.MediaType != mediaType {
			continue
		}
		if publishedOnly && !it.IsPublished {
			continue
		}
		out = append(out, it)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].DisplayOrder != out[j].DisplayOrder {
			return out[i].DisplayOrder < out[j].DisplayOrder
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return r.enrichCategory(out)
}

func (r *MediaRepo) ListFeatured(ctx context.Context, limit int) ([]models.PortfolioMedia, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	list, err := r.load()
	if err != nil {
		return nil, err
	}
	out := make([]models.PortfolioMedia, 0, len(list))
	for _, it := range list {
		if it.IsFeatured && it.IsPublished {
			out = append(out, it)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].DisplayOrder != out[j].DisplayOrder {
			return out[i].DisplayOrder < out[j].DisplayOrder
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return r.enrichCategory(out)
}

func (r *MediaRepo) ListLatest(ctx context.Context, limit int) ([]models.PortfolioMedia, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	list, err := r.load()
	if err != nil {
		return nil, err
	}
	out := make([]models.PortfolioMedia, 0, len(list))
	for _, it := range list {
		if it.IsPublished && it.MediaType == "photo" {
			out = append(out, it)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return r.enrichCategory(out)
}

func (r *MediaRepo) GetByID(ctx context.Context, id int64) (*models.PortfolioMedia, error) {
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
			enriched, err := r.enrichCategory([]models.PortfolioMedia{cp})
			if err != nil {
				return nil, err
			}
			out := enriched[0]
			return &out, nil
		}
	}
	return nil, ErrNotFound
}

func (r *MediaRepo) Create(ctx context.Context, m *models.PortfolioMedia) (int64, error) {
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
	n := *m
	n.ID = maxID + 1
	n.CreatedAt = now
	n.UpdatedAt = now
	list = append(list, n)
	if err := r.save(list); err != nil {
		return 0, err
	}
	return n.ID, nil
}

func (r *MediaRepo) Update(ctx context.Context, m *models.PortfolioMedia) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	list, err := r.load()
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	for i := range list {
		if list[i].ID == m.ID {
			list[i].Title = m.Title
			list[i].Caption = m.Caption
			list[i].IsFeatured = m.IsFeatured
			list[i].DisplayOrder = m.DisplayOrder
			list[i].IsPublished = m.IsPublished
			list[i].UpdatedAt = now
			return r.save(list)
		}
	}
	return ErrNotFound
}

func (r *MediaRepo) Delete(ctx context.Context, id int64) (*models.PortfolioMedia, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	list, err := r.load()
	if err != nil {
		return nil, err
	}
	var deleted *models.PortfolioMedia
	out := list[:0]
	for _, it := range list {
		if it.ID == id {
			cp := it
			deleted = &cp
			continue
		}
		out = append(out, it)
	}
	if deleted == nil {
		return nil, ErrNotFound
	}
	if err := r.save(out); err != nil {
		return nil, err
	}
	enriched, err := r.enrichCategory([]models.PortfolioMedia{*deleted})
	if err != nil {
		return nil, err
	}
	ret := enriched[0]
	return &ret, nil
}

func (r *MediaRepo) CountAll(ctx context.Context) (int, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	list, err := r.load()
	if err != nil {
		return 0, err
	}
	return len(list), nil
}

func (r *MediaRepo) CountRecent(ctx context.Context, days int) (int, error) {
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
		if it.CreatedAt.After(since) || it.CreatedAt.Equal(since) {
			n++
		}
	}
	return n, nil
}

func (r *MediaRepo) enrichCategory(list []models.PortfolioMedia) ([]models.PortfolioMedia, error) {
	var cats []models.Category
	err := readJSONFile(r.categoriesFile, &cats)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	cmap := map[int64]models.Category{}
	for _, c := range cats {
		cmap[c.ID] = c
	}
	out := make([]models.PortfolioMedia, 0, len(list))
	for _, m := range list {
		if c, ok := cmap[m.CategoryID]; ok {
			m.CategorySlug = c.Slug
			m.CategoryName = c.Name
		}
		out = append(out, m)
	}
	return out, nil
}
