package models

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"kyooar/internal/shared/models"
)

type Organization struct {
	models.BaseModel
	AccountID   uuid.UUID      `gorm:"not null" json:"account_id"`
	// Account     Account        `json:"account,omitempty"` // TODO: Add when cross-domain refs are ready
	Name        string         `gorm:"not null" json:"name"`
	Description string         `json:"description"`
	Logo        string         `json:"logo"`
	Website     string         `json:"website"`
	Phone       string         `json:"phone"`
	Email       string         `json:"email"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	Settings    Settings       `gorm:"type:jsonb" json:"settings"`
	Locations   []Location     `json:"locations,omitempty"`
	// Products      []Product         `json:"products,omitempty"`     // TODO: Add when menu domain is ready
	// QRCodes     []QRCode       `json:"qr_codes,omitempty"`  // TODO: Add when qrcode domain is ready
}

type Settings struct {
	Language             string `json:"language"`
	Timezone             string `json:"timezone"`
	FeedbackNotification bool   `json:"feedback_notification"`
	LowRatingThreshold   int    `json:"low_rating_threshold"`
}

// GORM Scanner/Valuer interfaces for JSONB
func (s Settings) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *Settings) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return json.Unmarshal([]byte("{}"), s)
	}
	return json.Unmarshal(bytes, s)
}

type Location struct {
	models.BaseModel
	OrganizationID uuid.UUID  `gorm:"not null" json:"organization_id"`
	Organization   Organization `json:"organization,omitempty"`
	Name         string     `gorm:"not null" json:"name"`
	Address      string     `json:"address"`
	City         string     `json:"city"`
	State        string     `json:"state"`
	Country      string     `json:"country"`
	PostalCode   string     `json:"postal_code"`
	Latitude     float64    `json:"latitude"`
	Longitude    float64    `json:"longitude"`
	IsActive     bool       `gorm:"default:true" json:"is_active"`
	// QRCodes      []QRCode   `json:"qr_codes,omitempty"` // TODO: Add when qrcode domain is ready
}

type Product struct {
	models.BaseModel
	OrganizationID  uuid.UUID         `gorm:"not null" json:"organization_id"`
	Organization    Organization        `json:"organization,omitempty"`
	Name          string            `gorm:"not null" json:"name"`
	Description   string            `json:"description"`
	Category      string            `json:"category"`
	Price         float64           `json:"price"`
	Currency      string            `gorm:"default:'USD'" json:"currency"`
	Image         string            `json:"image"`
	Tags          pq.StringArray    `gorm:"type:text[]" json:"tags"`
	IsAvailable   bool              `gorm:"default:true" json:"is_available"`
	IsActive      bool              `gorm:"default:true" json:"is_active"`
	DisplayOrder  int               `gorm:"default:0" json:"display_order"`
	// Questionnaire *Questionnaire    `json:"questionnaire,omitempty"` // TODO: Add when feedback domain is ready
	// Feedbacks     []Feedback        `json:"feedbacks,omitempty"`     // TODO: Add when feedback domain is ready
}
