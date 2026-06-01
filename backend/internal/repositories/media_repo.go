package repositories

import (
	"context"
	"database/sql"

	"github.com/suryaphotography/backend/internal/models"
)

type MediaRepo struct{ db *sql.DB }

func NewMediaRepo(db *sql.DB) *MediaRepo { return &MediaRepo{db: db} }

const mediaSelect = `
	SELECT pm.id, pm.category_id, pm.title, pm.caption, pm.media_type, pm.file_path, pm.thumbnail_path,
		pm.mime_type, pm.file_size_bytes, pm.duration_sec, pm.is_featured, pm.display_order, pm.is_published,
		pm.created_at, pm.updated_at, c.slug, c.name
	FROM dbo.portfolio_media pm
	INNER JOIN dbo.categories c ON c.id = pm.category_id`

func scanMedia(rows *sql.Rows) ([]models.PortfolioMedia, error) {
	var list []models.PortfolioMedia
	for rows.Next() {
		var m models.PortfolioMedia
		var title, caption, thumb, mime sql.NullString
		var size sql.NullInt64
		var dur sql.NullInt32
		if err := rows.Scan(&m.ID, &m.CategoryID, &title, &caption, &m.MediaType, &m.FilePath, &thumb,
			&mime, &size, &dur, &m.IsFeatured, &m.DisplayOrder, &m.IsPublished,
			&m.CreatedAt, &m.UpdatedAt, &m.CategorySlug, &m.CategoryName); err != nil {
			return nil, err
		}
		if title.Valid {
			m.Title = &title.String
		}
		if caption.Valid {
			m.Caption = &caption.String
		}
		if thumb.Valid {
			m.ThumbnailPath = &thumb.String
		}
		if mime.Valid {
			m.MimeType = &mime.String
		}
		if size.Valid {
			m.FileSizeBytes = &size.Int64
		}
		if dur.Valid {
			d := int(dur.Int32)
			m.DurationSec = &d
		}
		list = append(list, m)
	}
	return list, rows.Err()
}

func (r *MediaRepo) ListByCategory(ctx context.Context, categoryID int64, mediaType string, publishedOnly bool) ([]models.PortfolioMedia, error) {
	q := mediaSelect + ` WHERE pm.category_id = @p1`
	args := []any{categoryID}
	if mediaType != "" {
		q += ` AND pm.media_type = @p2`
		args = append(args, mediaType)
	}
	if publishedOnly {
		q += ` AND pm.is_published = 1`
	}
	q += ` ORDER BY pm.display_order, pm.created_at DESC`

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMedia(rows)
}

func (r *MediaRepo) ListFeatured(ctx context.Context, limit int) ([]models.PortfolioMedia, error) {
	rows, err := r.db.QueryContext(ctx, mediaSelect+`
		WHERE pm.is_featured = 1 AND pm.is_published = 1
		ORDER BY pm.display_order, pm.created_at DESC
		OFFSET 0 ROWS FETCH NEXT @p1 ROWS ONLY`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMedia(rows)
}

func (r *MediaRepo) ListLatest(ctx context.Context, limit int) ([]models.PortfolioMedia, error) {
	rows, err := r.db.QueryContext(ctx, mediaSelect+`
		WHERE pm.is_published = 1 AND pm.media_type = 'photo'
		ORDER BY pm.created_at DESC
		OFFSET 0 ROWS FETCH NEXT @p1 ROWS ONLY`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMedia(rows)
}

func (r *MediaRepo) GetByID(ctx context.Context, id int64) (*models.PortfolioMedia, error) {
	row := r.db.QueryRowContext(ctx, mediaSelect+` WHERE pm.id = @p1`, id)
	var m models.PortfolioMedia
	var title, caption, thumb, mime sql.NullString
	var size sql.NullInt64
	var dur sql.NullInt32
	err := row.Scan(&m.ID, &m.CategoryID, &title, &caption, &m.MediaType, &m.FilePath, &thumb,
		&mime, &size, &dur, &m.IsFeatured, &m.DisplayOrder, &m.IsPublished,
		&m.CreatedAt, &m.UpdatedAt, &m.CategorySlug, &m.CategoryName)
	if err != nil {
		return nil, err
	}
	if title.Valid {
		m.Title = &title.String
	}
	if caption.Valid {
		m.Caption = &caption.String
	}
	if thumb.Valid {
		m.ThumbnailPath = &thumb.String
	}
	if mime.Valid {
		m.MimeType = &mime.String
	}
	if size.Valid {
		m.FileSizeBytes = &size.Int64
	}
	if dur.Valid {
		d := int(dur.Int32)
		m.DurationSec = &d
	}
	return &m, nil
}

func (r *MediaRepo) Create(ctx context.Context, m *models.PortfolioMedia) (int64, error) {
	var id int64
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO dbo.portfolio_media (category_id, title, caption, media_type, file_path, thumbnail_path, mime_type, file_size_bytes, duration_sec, is_featured, display_order, is_published)
		OUTPUT INSERTED.id
		VALUES (@p1,@p2,@p3,@p4,@p5,@p6,@p7,@p8,@p9,@p10,@p11,@p12)`,
		m.CategoryID, m.Title, m.Caption, m.MediaType, m.FilePath, m.ThumbnailPath, m.MimeType, m.FileSizeBytes, m.DurationSec, m.IsFeatured, m.DisplayOrder, m.IsPublished,
	).Scan(&id)
	return id, err
}

func (r *MediaRepo) Update(ctx context.Context, m *models.PortfolioMedia) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE dbo.portfolio_media SET title=@p1, caption=@p2, is_featured=@p3, display_order=@p4, is_published=@p5, updated_at=SYSUTCDATETIME()
		WHERE id=@p6`, m.Title, m.Caption, m.IsFeatured, m.DisplayOrder, m.IsPublished, m.ID)
	return err
}

func (r *MediaRepo) Delete(ctx context.Context, id int64) (*models.PortfolioMedia, error) {
	m, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	_, err = r.db.ExecContext(ctx, `DELETE FROM dbo.portfolio_media WHERE id = @p1`, id)
	return m, err
}

func (r *MediaRepo) CountAll(ctx context.Context) (int, error) {
	var n int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM dbo.portfolio_media`).Scan(&n)
	return n, err
}

func (r *MediaRepo) CountRecent(ctx context.Context, days int) (int, error) {
	var n int
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM dbo.portfolio_media WHERE created_at >= DATEADD(day, -@p1, SYSUTCDATETIME())`, days).Scan(&n)
	return n, err
}
