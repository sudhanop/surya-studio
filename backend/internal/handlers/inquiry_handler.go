package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/suryaphotography/backend/internal/email"
	"github.com/suryaphotography/backend/internal/jsonutil"
	"github.com/suryaphotography/backend/internal/models"
	"github.com/suryaphotography/backend/internal/repositories"
)

type InquiryHandler struct {
	repo  *repositories.InquiryRepo
	funcRepo *repositories.FunctionRepo
	mail  *email.Service
}

func NewInquiryHandler(repo *repositories.InquiryRepo, funcRepo *repositories.FunctionRepo, mail *email.Service) *InquiryHandler {
	return &InquiryHandler{repo: repo, funcRepo: funcRepo, mail: mail}
}

type createInquiryRequest struct {
	CustomerName string `json:"customer_name"`
	PhoneNumber  string `json:"phone_number"`
	OccasionType string `json:"occasion_type"`
	WantedDate   string `json:"wanted_date"`
	Address      string `json:"address"`
	Message      string `json:"message"`
}

// Create godoc
// @Summary Submit contact inquiry
// @Tags inquiries
// @Accept json
// @Produce json
// @Param body body createInquiryRequest true "Inquiry"
// @Success 201 {object} models.MessageResponse
// @Router /api/inquiries [post]
func (h *InquiryHandler) Create(c *fiber.Ctx) error {
	var req createInquiryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}
	if req.CustomerName == "" || req.PhoneNumber == "" || req.OccasionType == "" {
		return c.Status(400).JSON(fiber.Map{"error": "name, phone and occasion are required"})
	}
	inquiry := &models.Inquiry{
		CustomerName: req.CustomerName,
		PhoneNumber:  req.PhoneNumber,
		OccasionType: req.OccasionType,
		Status:       "new",
	}
	if req.Address != "" {
		inquiry.Address = &req.Address
	}
	if req.Message != "" {
		inquiry.Message = &req.Message
	}
	if req.WantedDate != "" {
		if t, err := time.Parse("2006-01-02", req.WantedDate); err == nil {
			inquiry.WantedDate = &t
		}
	}
	id, err := h.repo.Create(c.Context(), inquiry)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "could not save inquiry"})
	}
	inquiry.ID = id
	_ = h.mail.SendInquiryNotification(inquiry)
	return c.Status(201).JSON(fiber.Map{"message": "thank you — we will contact you soon"})
}

// List godoc
// @Summary List inquiries
// @Tags admin-inquiries
// @Security BearerAuth
// @Router /api/admin/inquiries [get]
func (h *InquiryHandler) List(c *fiber.Ctx) error {
	status := c.Query("status")
	list, err := h.repo.List(c.Context(), status, 100)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed"})
	}
	return c.JSON(jsonutil.Slice(list))
}

type updateInquiryStatusRequest struct {
	Status string `json:"status"`
}

// UpdateStatus godoc
// @Summary Update inquiry status
// @Tags admin-inquiries
// @Security BearerAuth
// @Param id path int true "Inquiry ID"
// @Router /api/admin/inquiries/{id}/status [put]
func (h *InquiryHandler) UpdateStatus(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	var req updateInquiryStatusRequest
	if err := c.BodyParser(&req); err != nil || req.Status == "" {
		return c.Status(400).JSON(fiber.Map{"error": "status required"})
	}
	if err := h.repo.UpdateStatus(c.Context(), int64(id), req.Status); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "update failed"})
	}
	return c.JSON(fiber.Map{"message": "updated"})
}

// Delete godoc
// @Summary Delete inquiry
// @Tags admin-inquiries
// @Security BearerAuth
// @Router /api/admin/inquiries/{id} [delete]
func (h *InquiryHandler) Delete(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	if err := h.repo.Delete(c.Context(), int64(id)); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "delete failed"})
	}
	return c.JSON(fiber.Map{"message": "deleted"})
}

type convertInquiryRequest struct {
	FunctionDate string  `json:"function_date"`
	TotalAmount  float64 `json:"total_amount"`
	AdvancePaid  float64 `json:"advance_paid"`
}

// Convert godoc
// @Summary Convert inquiry to function record
// @Tags admin-inquiries
// @Security BearerAuth
// @Param id path int true "Inquiry ID"
// @Router /api/admin/inquiries/{id}/convert [post]
func (h *InquiryHandler) Convert(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	inquiry, err := h.repo.GetByID(c.Context(), int64(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	var req convertInquiryRequest
	_ = c.BodyParser(&req)
	fnDate := time.Now().AddDate(0, 1, 0)
	if req.FunctionDate != "" {
		if t, e := time.Parse("2006-01-02", req.FunctionDate); e == nil {
			fnDate = t
		}
	} else if inquiry.WantedDate != nil {
		fnDate = *inquiry.WantedDate
	}
	inquiryID := inquiry.ID
	f := &models.Function{
		InquiryID:      &inquiryID,
		CustomerName:   inquiry.CustomerName,
		PhoneNumber:    inquiry.PhoneNumber,
		Address:        inquiry.Address,
		FunctionType:   inquiry.OccasionType,
		FunctionDate:   fnDate,
		TotalAmount:    req.TotalAmount,
		AdvancePaid:    req.AdvancePaid,
		AlbumStatus:    "pending",
		VideoStatus:    "pending",
		DeliveryStatus: "pending",
		OverallStatus:  "upcoming",
	}
	fid, err := h.funcRepo.Create(c.Context(), f)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "convert failed"})
	}
	_ = h.repo.UpdateStatus(c.Context(), inquiry.ID, "converted")
	return c.Status(201).JSON(fiber.Map{"function_id": fid, "message": "converted"})
}
