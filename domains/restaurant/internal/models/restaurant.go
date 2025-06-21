package models

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/lecritique/api/shared/models"
)

type Restaurant struct {
	models.BaseModel
	AccountID   uuid.UUID `gorm:"not null" json:"account_id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	Logo        string    `json:"logo"`
	Website     string    `json:"website"`
	Phone       string    `json:"phone"`
	Email       string    `json:"email"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	Settings    Settings  `gorm:"type:jsonb" json:"settings"`
	Locations   []Location `json:"locations,omitempty"`
}

type Settings struct {
	Language             string `json:"language"`
	Timezone             string `json:"timezone"`
	FeedbackNotification bool   `json:"feedback_notification"`
	LowRatingThreshold   int    `json:"low_rating_threshold"`
}

type Location struct {
	models.BaseModel
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
}

// GORM Scanner/Valuer interface for Settings JSONB
func (s Settings) Value() (driver.Value, error) { return json.Marshal(s) }
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