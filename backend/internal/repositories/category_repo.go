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

type CategoryRepo struct {
	file      string
	mediaFile string
	mu        sync.Mutex
}

func NewCategoryRepo(dataDir string) *CategoryRepo {
	return &CategoryRepo{
		file:      filepath.Join(dataDir, "categories.json"),
		mediaFile: filepath.Join(dataDir, "media.json"),
	}
}

func (r *CategoryRepo) load() ([]models.Category, error) {
	var list []models.Category
	err := readJSONFile(r.file, &list)
	if err != nil {
		if os.IsNotExist(err) {
			return []models.Category{}, nil
		}
		return nil, err
	}
	if list == nil {
		list = []models.Category{}
	}
	return list, nil
}

func (r *CategoryRepo) save(list []models.Category) error {
	return writeJSONFile(r.file, list)
}

func (r *CategoryRepo) withCounts(list []models.Category) ([]models.Category, error) {
	var media []models.PortfolioMedia
	err := readJSONFile(r.mediaFile, &media)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	photos := map[int64]int{}
	videos := map[int64]int{}
	for _, m := range media {
		if !m.IsPublished {
			continue
		}
		if m.MediaType == "video" {
			videos[m.CategoryID]++
		} else {
			photos[m.CategoryID]++
		}
	}
	out := make([]models.Category, 0, len(list))
	for _, c := range list {
		c.PhotoCount = photos[c.ID]
		c.VideoCount = videos[c.ID]
		out = append(out, c)
	}
	return out, nil
}

func (r *CategoryRepo) List(ctx context.Context, activeOnly bool) ([]models.Category, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	list, err := r.load()
	if err != nil {
		return nil, err
	}
	out := make([]models.Category, 0, len(list))
	for _, c := range list {
		if activeOnly && !c.IsActive {
			continue
		}
		out = append(out, c)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].DisplayOrder != out[j].DisplayOrder {
			return out[i].DisplayOrder < out[j].DisplayOrder
		}
		return out[i].Name < out[j].Name
	})
	return r.withCounts(out)
}

func (r *CategoryRepo) GetBySlug(ctx context.Context, slug string) (*models.Category, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	list, err := r.load()
	if err != nil {
		return nil, err
	}
	for _, c := range list {
		if c.Slug == slug {
			enriched, err := r.withCounts([]models.Category{c})
			if err != nil {
				return nil, err
			}
			out := enriched[0]
			return &out, nil
		}
	}
	return nil, ErrNotFound
}

func (r *CategoryRepo) GetByID(ctx context.Context, id int64) (*models.Category, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	list, err := r.load()
	if err != nil {
		return nil, err
	}
	for _, c := range list {
		if c.ID == id {
			enriched, err := r.withCounts([]models.Category{c})
			if err != nil {
				return nil, err
			}
			out := enriched[0]
			return &out, nil
		}
	}
	return nil, ErrNotFound
}

func (r *CategoryRepo) Update(ctx context.Context, c *models.Category) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	list, err := r.load()
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	for i := range list {
		if list[i].ID == c.ID {
			list[i].Name = c.Name
			list[i].Description = c.Description
			list[i].CoverImage = c.CoverImage
			list[i].DisplayOrder = c.DisplayOrder
			list[i].IsActive = c.IsActive
			list[i].UpdatedAt = now
			return r.save(list)
		}
	}
	return ErrNotFound
}
