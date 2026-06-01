package repositories

import (
	"context"
	"database/sql"

	"github.com/suryaphotography/backend/internal/models"
)

type CategoryRepo struct{ db *sql.DB }

func NewCategoryRepo(db *sql.DB) *CategoryRepo { return &CategoryRepo{db: db} }

func scanCategory(row interface{ Scan(...any) error }) (*models.Category, error) {
	var c models.Category
	var desc, cover sql.NullString
	err := row.Scan(&c.ID, &c.Name, &c.Slug, &desc, &cover, &c.DisplayOrder, &c.IsActive, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if desc.Valid {
		c.Description = &desc.String
	}
	if cover.Valid {
		c.CoverImage = &cover.String
	}
	return &c, nil
}

const categorySelect = `
	SELECT c.id, c.name, c.slug, c.description, c.cover_image, c.display_order, c.is_active, c.created_at, c.updated_at,
		(SELECT COUNT(*) FROM dbo.portfolio_media pm WHERE pm.category_id = c.id AND pm.media_type = 'photo' AND pm.is_published = 1),
		(SELECT COUNT(*) FROM dbo.portfolio_media pm WHERE pm.category_id = c.id AND pm.media_type = 'video' AND pm.is_published = 1)
	FROM dbo.categories c`

func (r *CategoryRepo) List(ctx context.Context, activeOnly bool) ([]models.Category, error) {
	q := categorySelect + ` ORDER BY c.display_order, c.name`
	if activeOnly {
		q = categorySelect + ` WHERE c.is_active = 1 ORDER BY c.display_order, c.name`
	}
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Category
	for rows.Next() {
		var c models.Category
		var desc, cover sql.NullString
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug, &desc, &cover, &c.DisplayOrder, &c.IsActive, &c.CreatedAt, &c.UpdatedAt, &c.PhotoCount, &c.VideoCount); err != nil {
			return nil, err
		}
		if desc.Valid {
			c.Description = &desc.String
		}
		if cover.Valid {
			c.CoverImage = &cover.String
		}
		list = append(list, c)
	}
	return list, rows.Err()
}

func (r *CategoryRepo) GetBySlug(ctx context.Context, slug string) (*models.Category, error) {
	row := r.db.QueryRowContext(ctx, categorySelect+` WHERE c.slug = @p1`, slug)
	var c models.Category
	var desc, cover sql.NullString
	err := row.Scan(&c.ID, &c.Name, &c.Slug, &desc, &cover, &c.DisplayOrder, &c.IsActive, &c.CreatedAt, &c.UpdatedAt, &c.PhotoCount, &c.VideoCount)
	if err != nil {
		return nil, err
	}
	if desc.Valid {
		c.Description = &desc.String
	}
	if cover.Valid {
		c.CoverImage = &cover.String
	}
	return &c, nil
}

func (r *CategoryRepo) GetByID(ctx context.Context, id int64) (*models.Category, error) {
	row := r.db.QueryRowContext(ctx, categorySelect+` WHERE c.id = @p1`, id)
	var c models.Category
	var desc, cover sql.NullString
	err := row.Scan(&c.ID, &c.Name, &c.Slug, &desc, &cover, &c.DisplayOrder, &c.IsActive, &c.CreatedAt, &c.UpdatedAt, &c.PhotoCount, &c.VideoCount)
	if err != nil {
		return nil, err
	}
	if desc.Valid {
		c.Description = &desc.String
	}
	if cover.Valid {
		c.CoverImage = &cover.String
	}
	return &c, nil
}

func (r *CategoryRepo) Update(ctx context.Context, c *models.Category) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE dbo.categories SET name=@p1, description=@p2, cover_image=@p3, display_order=@p4, is_active=@p5, updated_at=SYSUTCDATETIME()
		WHERE id=@p6`, c.Name, c.Description, c.CoverImage, c.DisplayOrder, c.IsActive, c.ID)
	return err
}
