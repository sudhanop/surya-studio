package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"sort"
	"time"

	"github.com/suryaphotography/backend/internal/models"
)

type FunctionRepo struct{ db *sql.DB }

func NewFunctionRepo(db *sql.DB) *FunctionRepo { return &FunctionRepo{db: db} }

func scanFunction(row interface{ Scan(...any) error }) (*models.Function, error) {
	var f models.Function
	var inquiry sql.NullInt64
	var addr, editor, notes, links, bookingNotes, servicesJSON, complimentaryJSON sql.NullString
	var assigned sql.NullTime
	err := row.Scan(&f.ID, &inquiry, &f.CustomerName, &f.PhoneNumber, &addr, &f.FunctionType, &f.FunctionDate,
		&f.TotalAmount, &f.AdvancePaid, &f.BalanceAmount,
		&editor, &assigned, &f.AlbumStatus, &f.VideoStatus, &f.DeliveryStatus, &f.OverallStatus,
		&bookingNotes, &servicesJSON, &complimentaryJSON, &notes, &links, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if inquiry.Valid {
		f.InquiryID = &inquiry.Int64
	}
	if addr.Valid {
		f.Address = &addr.String
	}
	if editor.Valid {
		f.AssignedEditor = &editor.String
	}
	if assigned.Valid {
		f.AssignedDate = &assigned.Time
	}
	if bookingNotes.Valid {
		f.CustomerBookingNotes = &bookingNotes.String
	}
	if notes.Valid {
		f.AdminNotes = &notes.String
	}
	if links.Valid {
		f.DriveLinks = &links.String
	}
	if servicesJSON.Valid && servicesJSON.String != "" {
		_ = json.Unmarshal([]byte(servicesJSON.String), &f.Services)
	}
	if complimentaryJSON.Valid && complimentaryJSON.String != "" {
		_ = json.Unmarshal([]byte(complimentaryJSON.String), &f.Complimentary)
	}
	if f.Services == nil {
		f.Services = []string{}
	}
	if f.Complimentary == nil {
		f.Complimentary = []string{}
	}
	return &f, nil
}

const functionSelect = `
	SELECT id, inquiry_id, customer_name, phone_number, address, function_type, function_date,
		total_amount, advance_paid, balance_amount, assigned_editor, assigned_date,
		album_status, video_status, delivery_status, overall_status,
		customer_booking_notes, services_json, complimentary_json, admin_notes, drive_links, created_at, updated_at
	FROM dbo.functions`

func (r *FunctionRepo) enrich(ctx context.Context, f *models.Function) error {
	dates, err := r.ListEventDates(ctx, f.ID)
	if err != nil {
		return err
	}
	f.EventDates = dates
	return nil
}

func (r *FunctionRepo) enrichAll(ctx context.Context, list []models.Function) ([]models.Function, error) {
	for i := range list {
		if err := r.enrich(ctx, &list[i]); err != nil {
			return nil, err
		}
	}
	return list, nil
}

func (r *FunctionRepo) List(ctx context.Context, status string, limit int) ([]models.Function, error) {
	var rows *sql.Rows
	var err error
	if status != "" {
		rows, err = r.db.QueryContext(ctx, functionSelect+` WHERE overall_status = @p1 ORDER BY function_date ASC OFFSET 0 ROWS FETCH NEXT @p2 ROWS ONLY`, status, limit)
	} else {
		rows, err = r.db.QueryContext(ctx, functionSelect+` ORDER BY function_date ASC OFFSET 0 ROWS FETCH NEXT @p1 ROWS ONLY`, limit)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list, err := r.scanAll(rows)
	if err != nil {
		return nil, err
	}
	return r.enrichAll(ctx, list)
}

func (r *FunctionRepo) scanAll(rows *sql.Rows) ([]models.Function, error) {
	var list []models.Function
	for rows.Next() {
		f, err := scanFunction(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, *f)
	}
	return list, rows.Err()
}

func (r *FunctionRepo) Upcoming(ctx context.Context, limit int) ([]models.Function, error) {
	rows, err := r.db.QueryContext(ctx, functionSelect+`
		WHERE overall_status IN ('upcoming', 'editing') AND function_date >= CAST(GETUTCDATE() AS DATE)
		ORDER BY function_date ASC OFFSET 0 ROWS FETCH NEXT @p1 ROWS ONLY`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list, err := r.scanAll(rows)
	if err != nil {
		return nil, err
	}
	return r.enrichAll(ctx, list)
}

func (r *FunctionRepo) GetByID(ctx context.Context, id int64) (*models.Function, error) {
	row := r.db.QueryRowContext(ctx, functionSelect+` WHERE id = @p1`, id)
	f, err := scanFunction(row)
	if err != nil {
		return nil, err
	}
	if err := r.enrich(ctx, f); err != nil {
		return nil, err
	}
	return f, nil
}

func servicesJSON(services []string) interface{} {
	if len(services) == 0 {
		return nil
	}
	b, _ := json.Marshal(services)
	return string(b)
}

func (r *FunctionRepo) Create(ctx context.Context, f *models.Function) (int64, error) {
	var id int64
	var assigned interface{}
	if f.AssignedDate != nil {
		assigned = f.AssignedDate.Format("2006-01-02")
	}
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO dbo.functions (inquiry_id, customer_name, phone_number, address, function_type, function_date,
			total_amount, advance_paid, assigned_editor, assigned_date, album_status, video_status, delivery_status, overall_status,
			customer_booking_notes, services_json, complimentary_json, admin_notes, drive_links)
		OUTPUT INSERTED.id
		VALUES (@p1,@p2,@p3,@p4,@p5,@p6,@p7,@p8,@p9,@p10,@p11,@p12,@p13,@p14,@p15,@p16,@p17,@p18,@p19)`,
		f.InquiryID, f.CustomerName, f.PhoneNumber, f.Address, f.FunctionType, f.FunctionDate.Format("2006-01-02"),
		f.TotalAmount, f.AdvancePaid, f.AssignedEditor, assigned,
		f.AlbumStatus, f.VideoStatus, f.DeliveryStatus, f.OverallStatus,
		f.CustomerBookingNotes, servicesJSON(f.Services), servicesJSON(f.Complimentary),
		f.AdminNotes, f.DriveLinks,
	).Scan(&id)
	return id, err
}

func (r *FunctionRepo) Update(ctx context.Context, f *models.Function) error {
	var assigned interface{}
	if f.AssignedDate != nil {
		assigned = f.AssignedDate.Format("2006-01-02")
	}
	_, err := r.db.ExecContext(ctx, `
		UPDATE dbo.functions SET customer_name=@p1, phone_number=@p2, address=@p3, function_type=@p4, function_date=@p5,
			total_amount=@p6, advance_paid=@p7, assigned_editor=@p8, assigned_date=@p9,
			album_status=@p10, video_status=@p11, delivery_status=@p12, overall_status=@p13,
			customer_booking_notes=@p14, services_json=@p15, complimentary_json=@p16,
			admin_notes=@p17, drive_links=@p18, updated_at=SYSUTCDATETIME() WHERE id=@p19`,
		f.CustomerName, f.PhoneNumber, f.Address, f.FunctionType, f.FunctionDate.Format("2006-01-02"),
		f.TotalAmount, f.AdvancePaid, f.AssignedEditor, assigned,
		f.AlbumStatus, f.VideoStatus, f.DeliveryStatus, f.OverallStatus,
		f.CustomerBookingNotes, servicesJSON(f.Services), servicesJSON(f.Complimentary),
		f.AdminNotes, f.DriveLinks, f.ID)
	return err
}

func (r *FunctionRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM dbo.functions WHERE id = @p1`, id)
	return err
}

func (r *FunctionRepo) ListEventDates(ctx context.Context, functionID int64) ([]models.FunctionEventDate, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, function_id, event_date, day_label, sort_order
		FROM dbo.function_event_dates WHERE function_id = @p1 ORDER BY sort_order, event_date`, functionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.FunctionEventDate
	for rows.Next() {
		var d models.FunctionEventDate
		var label sql.NullString
		if err := rows.Scan(&d.ID, &d.FunctionID, &d.EventDate, &label, &d.SortOrder); err != nil {
			return nil, err
		}
		if label.Valid {
			d.DayLabel = &label.String
		}
		list = append(list, d)
	}
	if list == nil {
		list = []models.FunctionEventDate{}
	}
	return list, rows.Err()
}

func (r *FunctionRepo) ReplaceEventDates(ctx context.Context, functionID int64, dates []models.FunctionEventDate) error {
	if _, err := r.db.ExecContext(ctx, `DELETE FROM dbo.function_event_dates WHERE function_id = @p1`, functionID); err != nil {
		return err
	}
	for i, d := range dates {
		label := d.DayLabel
		_, err := r.db.ExecContext(ctx, `
			INSERT INTO dbo.function_event_dates (function_id, event_date, day_label, sort_order)
			VALUES (@p1, @p2, @p3, @p4)`,
			functionID, d.EventDate.Format("2006-01-02"), label, i)
		if err != nil {
			return err
		}
	}
	return nil
}

func ParseEventDateStrings(dateStrs []string) ([]models.FunctionEventDate, time.Time, error) {
	if len(dateStrs) == 0 {
		return nil, time.Time{}, sql.ErrNoRows
	}
	var dates []models.FunctionEventDate
	for i, s := range dateStrs {
		s = trimDate(s)
		if s == "" {
			continue
		}
		t, err := time.Parse("2006-01-02", s)
		if err != nil {
			return nil, time.Time{}, err
		}
		dates = append(dates, models.FunctionEventDate{EventDate: t, SortOrder: i})
	}
	if len(dates) == 0 {
		return nil, time.Time{}, sql.ErrNoRows
	}
	sort.Slice(dates, func(i, j int) bool {
		return dates[i].EventDate.Before(dates[j].EventDate)
	})
	return dates, dates[0].EventDate, nil
}

func trimDate(s string) string {
	if len(s) > 10 {
		return s[:10]
	}
	return s
}

func (r *FunctionRepo) CountPendingAlbums(ctx context.Context) (int, error) {
	var n int
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM dbo.functions
		WHERE album_status NOT IN ('printed','delivered') AND overall_status NOT IN ('delivered')`).Scan(&n)
	return n, err
}

func (r *FunctionRepo) CountPendingVideos(ctx context.Context) (int, error) {
	var n int
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM dbo.functions
		WHERE video_status NOT IN ('completed','delivered') AND overall_status NOT IN ('delivered')`).Scan(&n)
	return n, err
}

func (r *FunctionRepo) CountPendingDeliveries(ctx context.Context) (int, error) {
	var n int
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM dbo.functions WHERE overall_status IN ('album_ready','editing','completed') AND delivery_status != 'delivered'`).Scan(&n)
	return n, err
}

func (r *FunctionRepo) CountUpcoming(ctx context.Context) (int, error) {
	var n int
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM dbo.functions WHERE overall_status = 'upcoming' AND function_date >= CAST(GETUTCDATE() AS DATE)`).Scan(&n)
	return n, err
}
