package email

import (
	"fmt"
	"time"

	"gopkg.in/gomail.v2"

	"github.com/suryaphotography/backend/internal/config"
	"github.com/suryaphotography/backend/internal/models"
)

type Service struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Service {
	return &Service{cfg: cfg}
}

func (s *Service) Enabled() bool {
	return s.cfg.SMTPHost != "" && s.cfg.AdminEmail != ""
}

func (s *Service) SendInquiryNotification(inquiry *models.Inquiry) error {
	if !s.Enabled() {
		return nil
	}
	subject := fmt.Sprintf("New Inquiry — %s (%s)", inquiry.CustomerName, inquiry.OccasionType)
	body := fmt.Sprintf(`New contact inquiry received:

Name: %s
Phone: %s
Occasion: %s
Wanted Date: %v
Address: %v
Message: %v

— Surya Photography Website`,
		inquiry.CustomerName,
		inquiry.PhoneNumber,
		inquiry.OccasionType,
		formatDate(inquiry.WantedDate),
		ptrStr(inquiry.Address),
		ptrStr(inquiry.Message),
	)

	m := gomail.NewMessage()
	m.SetHeader("From", s.cfg.SMTPFrom)
	m.SetHeader("To", s.cfg.AdminEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(s.cfg.SMTPHost, s.cfg.SMTPPort, s.cfg.SMTPUser, s.cfg.SMTPPassword)
	return d.DialAndSend(m)
}

func formatDate(t *time.Time) string {
	if t == nil {
		return "—"
	}
	return t.Format("2006-01-02")
}

func ptrStr(s *string) string {
	if s == nil {
		return "—"
	}
	return *s
}
