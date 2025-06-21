package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Restaurant struct {
	BaseModel
	AccountID   uuid.UUID      `gorm:"not null" json:"account_id"`
	Account     Account        `json:"account,omitempty"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `json:"description"`
	Logo        string         `json:"logo"`
	Website     string         `json:"website"`
	Phone       string         `json:"phone"`
	Email       string         `json:"email"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	Settings    Settings       `gorm:"type:jsonb" json:"settings"`
	Locations   []Location     `json:"locations,omitempty"`
	Dishes      []Dish         `json:"dishes,omitempty"`
	QRCodes     []QRCode       `json:"qr_codes,omitempty"`
}

type Settings struct {
	Language             string `json:"language"`
	Timezone             string `json:"timezone"`
	FeedbackNotification bool   `json:"feedback_notification"`
	LowRatingThreshold   int    `json:"low_rating_threshold"`
}

type Location struct {
	BaseModel
	RestaurantID uuid.UUID  `gorm:"not null" json:"restaurant_id"`
	Restaurant   Restaurant `json:"restaurant,omitempty"`
	Name         string     `gorm:"not null" json:"name"`
	Address      string     `json:"address"`
	City         string     `json:"city"`
	State        string     `json:"state"`
	Country      string     `json:"country"`
	PostalCode   string     `json:"postal_code"`
	Latitude     float64    `json:"latitude"`
	Longitude    float64    `json:"longitude"`
	IsActive     bool       `gorm:"default:true" json:"is_active"`
	QRCodes      []QRCode   `json:"qr_codes,omitempty"`
}

type Dish struct {
	BaseModel
	RestaurantID  uuid.UUID         `gorm:"not null" json:"restaurant_id"`
	Restaurant    Restaurant        `json:"restaurant,omitempty"`
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
	Questionnaire *Questionnaire    `json:"questionnaire,omitempty"`
	Feedbacks     []Feedback        `json:"feedbacks,omitempty"`
}
