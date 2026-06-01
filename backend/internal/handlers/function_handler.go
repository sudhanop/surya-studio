package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/suryaphotography/backend/internal/jsonutil"
	"github.com/suryaphotography/backend/internal/models"
	"github.com/suryaphotography/backend/internal/repositories"
)

type FunctionHandler struct {
	repo *repositories.FunctionRepo
}

func NewFunctionHandler(repo *repositories.FunctionRepo) *FunctionHandler {
	return &FunctionHandler{repo: repo}
}

type functionRequest struct {
	InquiryID            *int64   `json:"inquiry_id"`
	CustomerName         string   `json:"customer_name"`
	PhoneNumber          string   `json:"phone_number"`
	Address              *string  `json:"address"`
	FunctionType         string   `json:"function_type"`
	FunctionDate         string   `json:"function_date"`
	EventDates           []string `json:"event_dates"`
	TotalAmount          float64  `json:"total_amount"`
	AdvancePaid          float64  `json:"advance_paid"`
	AssignedEditor       *string  `json:"assigned_editor"`
	AssignedDate         string   `json:"assigned_date"`
	AlbumStatus          string   `json:"album_status"`
	VideoStatus          string   `json:"video_status"`
	DeliveryStatus       string   `json:"delivery_status"`
	OverallStatus        string   `json:"overall_status"`
	CustomerBookingNotes *string  `json:"customer_booking_notes"`
	Services             []string `json:"services"`
	Complimentary        []string `json:"complimentary"`
	AdminNotes           *string  `json:"admin_notes"`
	DriveLinks           *string  `json:"drive_links"`
}

func parseFunctionCreate(req functionRequest) (*models.Function, []models.FunctionEventDate, error) {
	if req.CustomerName == "" || req.PhoneNumber == "" || req.FunctionType == "" {
		return nil, nil, fiber.NewError(400, "missing required fields")
	}
	var fd time.Time
	var eventDates []models.FunctionEventDate
	var err error
	if len(req.EventDates) > 0 {
		eventDates, fd, err = repositories.ParseEventDateStrings(req.EventDates)
		if err != nil {
			return nil, nil, fiber.NewError(400, "invalid event_dates")
		}
	} else if req.FunctionDate != "" {
		fd, err = time.Parse("2006-01-02", req.FunctionDate)
		if err != nil {
			return nil, nil, fiber.NewError(400, "invalid function_date")
		}
		eventDates = []models.FunctionEventDate{{EventDate: fd}}
	} else {
		return nil, nil, fiber.NewError(400, "function_date or event_dates required")
	}
	f := &models.Function{
		InquiryID:            req.InquiryID,
		CustomerName:         req.CustomerName,
		PhoneNumber:          req.PhoneNumber,
		Address:              req.Address,
		FunctionType:         req.FunctionType,
		FunctionDate:         fd,
		TotalAmount:          req.TotalAmount,
		AdvancePaid:          req.AdvancePaid,
		AssignedEditor:       req.AssignedEditor,
		AlbumStatus:          defaultStr(req.AlbumStatus, "not_started"),
		VideoStatus:          defaultStr(req.VideoStatus, "not_started"),
		DeliveryStatus:       defaultStr(req.DeliveryStatus, "pending"),
		OverallStatus:        defaultStr(req.OverallStatus, "upcoming"),
		CustomerBookingNotes: req.CustomerBookingNotes,
		Services:             req.Services,
		Complimentary:        req.Complimentary,
		AdminNotes:           req.AdminNotes,
		DriveLinks:           req.DriveLinks,
	}
	if f.Services == nil {
		f.Services = []string{}
	}
	if f.Complimentary == nil {
		f.Complimentary = []string{}
	}
	if req.AssignedDate != "" {
		if t, e := time.Parse("2006-01-02", req.AssignedDate); e == nil {
			f.AssignedDate = &t
		}
	}
	return f, eventDates, nil
}

func mergeFunction(existing *models.Function, req functionRequest) (*models.Function, []models.FunctionEventDate, error) {
	f := *existing
	if req.CustomerName != "" {
		f.CustomerName = req.CustomerName
	}
	if req.PhoneNumber != "" {
		f.PhoneNumber = req.PhoneNumber
	}
	if req.Address != nil {
		f.Address = req.Address
	}
	if req.FunctionType != "" {
		f.FunctionType = req.FunctionType
	}
	var eventDates []models.FunctionEventDate
	if len(req.EventDates) > 0 {
		dates, primary, err := repositories.ParseEventDateStrings(req.EventDates)
		if err != nil {
			return nil, nil, fiber.NewError(400, "invalid event_dates")
		}
		f.FunctionDate = primary
		eventDates = dates
	} else if req.FunctionDate != "" {
		if t, e := time.Parse("2006-01-02", req.FunctionDate); e == nil {
			f.FunctionDate = t
		}
	}
	f.TotalAmount = req.TotalAmount
	f.AdvancePaid = req.AdvancePaid
	if req.AssignedEditor != nil {
		f.AssignedEditor = req.AssignedEditor
	}
	if req.AssignedDate != "" {
		if t, e := time.Parse("2006-01-02", req.AssignedDate); e == nil {
			f.AssignedDate = &t
		}
	}
	if req.AlbumStatus != "" {
		f.AlbumStatus = req.AlbumStatus
	}
	if req.VideoStatus != "" {
		f.VideoStatus = req.VideoStatus
	}
	if req.DeliveryStatus != "" {
		f.DeliveryStatus = req.DeliveryStatus
	}
	if req.OverallStatus != "" {
		f.OverallStatus = req.OverallStatus
	}
	if req.CustomerBookingNotes != nil {
		f.CustomerBookingNotes = req.CustomerBookingNotes
	}
	if req.Services != nil {
		f.Services = req.Services
	}
	if req.Complimentary != nil {
		f.Complimentary = req.Complimentary
	}
	if req.AdminNotes != nil {
		f.AdminNotes = req.AdminNotes
	}
	if req.DriveLinks != nil {
		f.DriveLinks = req.DriveLinks
	}
	if req.InquiryID != nil {
		f.InquiryID = req.InquiryID
	}
	return &f, eventDates, nil
}

func defaultStr(s, d string) string {
	if s == "" {
		return d
	}
	return s
}

func (h *FunctionHandler) List(c *fiber.Ctx) error {
	status := c.Query("status")
	list, err := h.repo.List(c.Context(), status, 200)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed"})
	}
	return c.JSON(jsonutil.Slice(list))
}

func (h *FunctionHandler) Upcoming(c *fiber.Ctx) error {
	list, err := h.repo.Upcoming(c.Context(), 20)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed"})
	}
	return c.JSON(jsonutil.Slice(list))
}

func (h *FunctionHandler) Create(c *fiber.Ctx) error {
	var req functionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}
	f, eventDates, err := parseFunctionCreate(req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	id, err := h.repo.Create(c.Context(), f)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "create failed"})
	}
	if len(eventDates) > 0 {
		_ = h.repo.ReplaceEventDates(c.Context(), id, eventDates)
	}
	f.ID = id
	created, _ := h.repo.GetByID(c.Context(), id)
	if created != nil {
		return c.Status(201).JSON(created)
	}
	return c.Status(201).JSON(f)
}

func (h *FunctionHandler) Update(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	existing, err := h.repo.GetByID(c.Context(), int64(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	var req functionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}
	f, eventDates, err := mergeFunction(existing, req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.repo.Update(c.Context(), f); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "update failed"})
	}
	if len(eventDates) > 0 {
		_ = h.repo.ReplaceEventDates(c.Context(), f.ID, eventDates)
	}
	updated, _ := h.repo.GetByID(c.Context(), f.ID)
	if updated != nil {
		return c.JSON(updated)
	}
	return c.JSON(f)
}

func (h *FunctionHandler) Delete(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	if err := h.repo.Delete(c.Context(), int64(id)); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "delete failed"})
	}
	return c.JSON(fiber.Map{"message": "deleted"})
}
