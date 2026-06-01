package models

import "time"

type Admin struct {
	ID           int64      `json:"id"`
	Username     string     `json:"username"`
	PasswordHash string     `json:"-"`
	Email        string     `json:"email"`
	DisplayName  *string    `json:"display_name,omitempty"`
	IsActive     bool       `json:"is_active"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
}

type Category struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	Description  *string   `json:"description,omitempty"`
	CoverImage   *string   `json:"cover_image,omitempty"`
	DisplayOrder int       `json:"display_order"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	PhotoCount   int       `json:"photo_count,omitempty"`
	VideoCount   int       `json:"video_count,omitempty"`
}

type PortfolioMedia struct {
	ID            int64     `json:"id"`
	CategoryID    int64     `json:"category_id"`
	Title         *string   `json:"title,omitempty"`
	Caption       *string   `json:"caption,omitempty"`
	MediaType     string    `json:"media_type"`
	FilePath      string    `json:"file_path"`
	ThumbnailPath *string   `json:"thumbnail_path,omitempty"`
	MimeType      *string   `json:"mime_type,omitempty"`
	FileSizeBytes *int64    `json:"file_size_bytes,omitempty"`
	DurationSec   *int      `json:"duration_sec,omitempty"`
	IsFeatured    bool      `json:"is_featured"`
	DisplayOrder  int       `json:"display_order"`
	IsPublished   bool      `json:"is_published"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	URL           string    `json:"url,omitempty"`
	ThumbnailURL  string    `json:"thumbnail_url,omitempty"`
	CategorySlug  string    `json:"category_slug,omitempty"`
	CategoryName  string    `json:"category_name,omitempty"`
}

type Inquiry struct {
	ID           int64      `json:"id"`
	CustomerName string     `json:"customer_name"`
	PhoneNumber  string     `json:"phone_number"`
	OccasionType string     `json:"occasion_type"`
	WantedDate   *time.Time `json:"wanted_date,omitempty"`
	Address      *string    `json:"address,omitempty"`
	Message      *string    `json:"message,omitempty"`
	Status       string     `json:"status"`
	ContactedAt  *time.Time `json:"contacted_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type FunctionEventDate struct {
	ID         int64     `json:"id,omitempty"`
	FunctionID int64     `json:"function_id,omitempty"`
	EventDate  time.Time `json:"event_date"`
	DayLabel   *string   `json:"day_label,omitempty"`
	SortOrder  int       `json:"sort_order"`
}

type Function struct {
	ID                   int64               `json:"id"`
	InquiryID            *int64              `json:"inquiry_id,omitempty"`
	CustomerName         string              `json:"customer_name"`
	PhoneNumber          string              `json:"phone_number"`
	Address              *string             `json:"address,omitempty"`
	FunctionType         string              `json:"function_type"`
	FunctionDate         time.Time           `json:"function_date"`
	EventDates           []FunctionEventDate `json:"event_dates,omitempty"`
	TotalAmount          float64             `json:"total_amount"`
	AdvancePaid          float64             `json:"advance_paid"`
	BalanceAmount        float64             `json:"balance_amount"`
	AssignedEditor       *string             `json:"assigned_editor,omitempty"`
	AssignedDate         *time.Time          `json:"assigned_date,omitempty"`
	AlbumStatus          string              `json:"album_status"`
	VideoStatus          string              `json:"video_status"`
	DeliveryStatus       string              `json:"delivery_status"`
	OverallStatus        string              `json:"overall_status"`
	CustomerBookingNotes *string             `json:"customer_booking_notes,omitempty"`
	Services             []string            `json:"services,omitempty"`
	Complimentary        []string            `json:"complimentary,omitempty"`
	AdminNotes           *string             `json:"admin_notes,omitempty"`
	DriveLinks           *string             `json:"drive_links,omitempty"`
	CreatedAt            time.Time           `json:"created_at"`
	UpdatedAt            time.Time           `json:"updated_at"`
}

type DashboardStats struct {
	UpcomingFunctions   int `json:"upcoming_functions"`
	RecentInquiries     int `json:"recent_inquiries"`
	PendingAlbums       int `json:"pending_albums"`
	PendingVideoEdits   int `json:"pending_video_edits"`
	PendingDeliveries   int `json:"pending_deliveries"`
	TotalUploads        int `json:"total_uploads"`
	LatestPortfolioCount int `json:"latest_portfolio_updates"`
}

type ContactInfo struct {
	WhatsApp         string `json:"whatsapp"`
	InstagramURL     string `json:"instagram_url"`
	FacebookURL      string `json:"facebook_url"`
	YouTubeURL       string `json:"youtube_url"`
	ContactEmail     string `json:"contact_email"`
	PhoneNumber      string `json:"phone_number"`
	PhoneSecondary   string `json:"phone_secondary,omitempty"`
	Address          string `json:"address,omitempty"`
	Pincode          string `json:"pincode,omitempty"`
	GoogleMapsEmbed  string `json:"google_maps_embed"`
}

type Testimonial struct {
	Name string `json:"name"`
	Text string `json:"text"`
}

type PublicSiteData struct {
	EventsCovered    string        `json:"events_covered"`
	YearsOfCraft     string        `json:"years_of_craft"`
	HappyFamilies    string        `json:"happy_families"`
	OwnerPortraitURL string        `json:"owner_portrait_url"`
	LogoURL          string        `json:"logo_url"`
	Testimonials     []Testimonial `json:"testimonials"`
	Contact          ContactInfo   `json:"contact"`
}

type APIError struct {
	Error string `json:"error"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
