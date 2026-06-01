package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/suryaphotography/backend/internal/models"
)

type InquiryRepo struct{ db *sql.DB }

func NewInquiryRepo(db *sql.DB) *InquiryRepo { return &InquiryRepo{db: db} }

func scanInquiry(row interface{ Scan(...any) error }) (*models.Inquiry, error) {
	var i models.Inquiry
	var wanted sql.NullTime
	var addr, msg sql.NullString
	var contacted sql.NullTime
	err := row.Scan(&i.ID, &i.CustomerName, &i.PhoneNumber, &i.OccasionType, &wanted, &addr, &msg, &i.Status, &contacted, &i.CreatedAt, &i.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if wanted.Valid {
		i.WantedDate = &wanted.Time
	}
	if addr.Valid {
		i.Address = &addr.String
	}
	if msg.Valid {
		i.Message = &msg.String
	}
	if contacted.Valid {
		i.ContactedAt = &contacted.Time
	}
	return &i, nil
}

func (r *InquiryRepo) Create(ctx context.Context, i *models.Inquiry) (int64, error) {
	var id int64
	var wanted interface{}
	if i.WantedDate != nil {
		wanted = i.WantedDate.Format("2006-01-02")
	}
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO dbo.inquiries (customer_name, phone_number, occasion_type, wanted_date, address, message)
		OUTPUT INSERTED.id
		VALUES (@p1,@p2,@p3,@p4,@p5,@p6)`,
		i.CustomerName, i.PhoneNumber, i.OccasionType, wanted, i.Address, i.Message,
	).Scan(&id)
	return id, err
}

func (r *InquiryRepo) List(ctx context.Context, status string, limit int) ([]models.Inquiry, error) {
	base := `SELECT id, customer_name, phone_number, occasion_type, wanted_date, address, message, status, contacted_at, created_at, updated_at FROM dbo.inquiries`
	var rows *sql.Rows
	var err error
	if status != "" {
		rows, err = r.db.QueryContext(ctx, base+` WHERE status = @p1 ORDER BY created_at DESC OFFSET 0 ROWS FETCH NEXT @p2 ROWS ONLY`, status, limit)
	} else {
		rows, err = r.db.QueryContext(ctx, base+` ORDER BY created_at DESC OFFSET 0 ROWS FETCH NEXT @p1 ROWS ONLY`, limit)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Inquiry
	for rows.Next() {
		i, err := scanInquiry(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, *i)
	}
	return list, rows.Err()
}

func (r *InquiryRepo) GetByID(ctx context.Context, id int64) (*models.Inquiry, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, customer_name, phone_number, occasion_type, wanted_date, address, message, status, contacted_at, created_at, updated_at
		FROM dbo.inquiries WHERE id = @p1`, id)
	return scanInquiry(row)
}

func (r *InquiryRepo) UpdateStatus(ctx context.Context, id int64, status string) error {
	var contacted interface{}
	if status == "contacted" {
		contacted = time.Now().UTC()
	}
	_, err := r.db.ExecContext(ctx, `
		UPDATE dbo.inquiries SET status=@p1, contacted_at=COALESCE(@p2, contacted_at), updated_at=SYSUTCDATETIME() WHERE id=@p3`,
		status, contacted, id)
	return err
}

func (r *InquiryRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM dbo.inquiries WHERE id = @p1`, id)
	return err
}

func (r *InquiryRepo) CountRecent(ctx context.Context, days int) (int, error) {
	var n int
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM dbo.inquiries WHERE created_at >= DATEADD(day, -@p1, SYSUTCDATETIME()) AND status = 'new'`, days).Scan(&n)
	return n, err
}
