package repositories

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/suryaphotography/backend/internal/models"
)

type FunctionRepo struct {
	file string
	mu   sync.Mutex
}

func NewFunctionRepo(dataDir string) *FunctionRepo {
	return &FunctionRepo{file: filepath.Join(dataDir, "functions.json")}
}

func (r *FunctionRepo) load() ([]models.Function, error) {
	var list []models.Function
	err := readJSONFile(r.file, &list)
	if err != nil {
		if os.IsNotExist(err) {
			return []models.Function{}, nil
		}
		return nil, err
	}
	if list == nil {
		list = []models.Function{}
	}
	for i := range list {
		normalizeFunction(&list[i])
	}
	return list, nil
}

func (r *FunctionRepo) save(list []models.Function) error {
	return writeJSONFile(r.file, list)
}

func (r *FunctionRepo) List(ctx context.Context, status string, limit int) ([]models.Function, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	list, err := r.load()
	if err != nil {
		return nil, err
	}

	out := make([]models.Function, 0, len(list))
	for _, f := range list {
		if status != "" && f.OverallStatus != status {
			continue
		}
		out = append(out, f)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].FunctionDate.Before(out[j].FunctionDate)
	})
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out, nil
}

func (r *FunctionRepo) Upcoming(ctx context.Context, limit int) ([]models.Function, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	list, err := r.load()
	if err != nil {
		return nil, err
	}
	today := todayUTC()

	out := make([]models.Function, 0, len(list))
	for _, f := range list {
		if f.OverallStatus != "upcoming" && f.OverallStatus != "editing" {
			continue
		}
		if f.FunctionDate.Before(today) {
			continue
		}
		out = append(out, f)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].FunctionDate.Before(out[j].FunctionDate)
	})
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out, nil
}

func (r *FunctionRepo) GetByID(ctx context.Context, id int64) (*models.Function, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	list, err := r.load()
	if err != nil {
		return nil, err
	}
	for _, f := range list {
		if f.ID == id {
			cp := f
			normalizeFunction(&cp)
			return &cp, nil
		}
	}
	return nil, ErrNotFound
}

func (r *FunctionRepo) Create(ctx context.Context, f *models.Function) (int64, error) {
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
	n := *f
	n.ID = maxID + 1
	n.CreatedAt = now
	n.UpdatedAt = now
	normalizeFunction(&n)
	if len(n.EventDates) == 0 {
		n.EventDates = []models.FunctionEventDate{{EventDate: n.FunctionDate, SortOrder: 0}}
	}
	n.FunctionDate = primaryEventDate(n)
	recomputeBalance(&n)

	list = append(list, n)
	if err := r.save(list); err != nil {
		return 0, err
	}
	return n.ID, nil
}

func (r *FunctionRepo) Update(ctx context.Context, f *models.Function) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	list, err := r.load()
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	for i := range list {
		if list[i].ID == f.ID {
			createdAt := list[i].CreatedAt
			existingDates := list[i].EventDates

			n := *f
			n.CreatedAt = createdAt
			n.UpdatedAt = now
			if n.EventDates == nil {
				n.EventDates = existingDates
			}
			normalizeFunction(&n)
			n.FunctionDate = primaryEventDate(n)
			recomputeBalance(&n)

			list[i] = n
			return r.save(list)
		}
	}
	return ErrNotFound
}

func (r *FunctionRepo) Delete(ctx context.Context, id int64) error {
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

func (r *FunctionRepo) ListEventDates(ctx context.Context, functionID int64) ([]models.FunctionEventDate, error) {
	f, err := r.GetByID(ctx, functionID)
	if err != nil {
		return nil, err
	}
	out := make([]models.FunctionEventDate, len(f.EventDates))
	copy(out, f.EventDates)
	sort.Slice(out, func(i, j int) bool {
		if out[i].SortOrder != out[j].SortOrder {
			return out[i].SortOrder < out[j].SortOrder
		}
		return out[i].EventDate.Before(out[j].EventDate)
	})
	return out, nil
}

func (r *FunctionRepo) ReplaceEventDates(ctx context.Context, functionID int64, dates []models.FunctionEventDate) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	list, err := r.load()
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	for i := range list {
		if list[i].ID == functionID {
			next := make([]models.FunctionEventDate, 0, len(dates))
			for j := range dates {
				d := dates[j]
				d.SortOrder = j
				next = append(next, d)
			}
			list[i].EventDates = next
			normalizeFunction(&list[i])
			list[i].FunctionDate = primaryEventDate(list[i])
			list[i].UpdatedAt = now
			return r.save(list)
		}
	}
	return ErrNotFound
}

func ParseEventDateStrings(dateStrs []string) ([]models.FunctionEventDate, time.Time, error) {
	if len(dateStrs) == 0 {
		return nil, time.Time{}, errors.New("no dates")
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
		return nil, time.Time{}, errors.New("no dates")
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
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	list, err := r.load()
	if err != nil {
		return 0, err
	}
	n := 0
	for _, f := range list {
		if f.OverallStatus == "delivered" {
			continue
		}
		if f.AlbumStatus == "printed" || f.AlbumStatus == "delivered" {
			continue
		}
		n++
	}
	return n, nil
}

func (r *FunctionRepo) CountPendingVideos(ctx context.Context) (int, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	list, err := r.load()
	if err != nil {
		return 0, err
	}
	n := 0
	for _, f := range list {
		if f.OverallStatus == "delivered" {
			continue
		}
		if f.VideoStatus == "completed" || f.VideoStatus == "delivered" {
			continue
		}
		n++
	}
	return n, nil
}

func (r *FunctionRepo) CountPendingDeliveries(ctx context.Context) (int, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	list, err := r.load()
	if err != nil {
		return 0, err
	}
	n := 0
	for _, f := range list {
		if f.DeliveryStatus == "delivered" {
			continue
		}
		if f.OverallStatus == "album_ready" || f.OverallStatus == "editing" || f.OverallStatus == "completed" {
			n++
		}
	}
	return n, nil
}

func (r *FunctionRepo) CountUpcoming(ctx context.Context) (int, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	list, err := r.load()
	if err != nil {
		return 0, err
	}
	today := todayUTC()
	n := 0
	for _, f := range list {
		if f.OverallStatus != "upcoming" {
			continue
		}
		if f.FunctionDate.Before(today) {
			continue
		}
		n++
	}
	return n, nil
}

func todayUTC() time.Time {
	now := time.Now().UTC()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
}

func normalizeFunction(f *models.Function) {
	if f.Services == nil {
		f.Services = []string{}
	}
	if f.Complimentary == nil {
		f.Complimentary = []string{}
	}
	if f.EventDates == nil {
		f.EventDates = []models.FunctionEventDate{}
	}
	recomputeBalance(f)
}

func recomputeBalance(f *models.Function) {
	f.BalanceAmount = f.TotalAmount - f.AdvancePaid
}

func primaryEventDate(f models.Function) time.Time {
	if len(f.EventDates) == 0 {
		return f.FunctionDate
	}
	min := f.EventDates[0].EventDate
	for i := 1; i < len(f.EventDates); i++ {
		if f.EventDates[i].EventDate.Before(min) {
			min = f.EventDates[i].EventDate
		}
	}
	return min
}
