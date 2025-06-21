package models

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/lecritique/api/shared/models"
)

type Feedback struct {
	models.BaseModel
	RestaurantID uuid.UUID  `gorm:"not null" json:"restaurant_id"`
	DishID       uuid.UUID  `gorm:"not null" json:"dish_id"`
	QRCodeID     uuid.UUID  `gorm:"not null" json:"qr_code_id"`
	CustomerName string     `json:"customer_name"`
	CustomerEmail string    `json:"customer_email"`
	CustomerPhone string    `json:"customer_phone"`
	OverallRating int        `json:"overall_rating"`
	Responses    Responses  `gorm:"type:jsonb" json:"responses"`
	DeviceInfo   DeviceInfo `gorm:"type:jsonb" json:"device_info"`
	IsComplete   bool       `gorm:"default:true" json:"is_complete"`
}

type Responses []Response

type Response struct {
	QuestionID uuid.UUID   `json:"question_id"`
	Answer     interface{} `json:"answer"`
}

type DeviceInfo struct {
	UserAgent string `json:"user_agent"`
	IP        string `json:"ip"`
	Platform  string `json:"platform"`
	Browser   string `json:"browser"`
}

// GORM Scanner/Valuer interfaces for JSONB
func (r Responses) Value() (driver.Value, error) { return json.Marshal(r) }
func (r *Responses) Scan(value interface{}) error { 
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return json.Unmarshal([]byte("[]"), r)
	}
	return json.Unmarshal(bytes, r) 
}

func (d DeviceInfo) Value() (driver.Value, error) { return json.Marshal(d) }
func (d *DeviceInfo) Scan(value interface{}) error { 
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return json.Unmarshal([]byte("{}"), d)
	}
	return json.Unmarshal(bytes, d) 
}