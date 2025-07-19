package models

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/google/uuid"
	menuModels "lecritique/internal/menu/models"
	qrcodeModels "lecritique/internal/qrcode/models"
	restaurantModels "lecritique/internal/restaurant/models"
	sharedModels "lecritique/internal/shared/models"
)

type Feedback struct {
	sharedModels.BaseModel
	RestaurantID uuid.UUID  `gorm:"not null" json:"restaurant_id"`
	Restaurant   restaurantModels.Restaurant `json:"restaurant,omitempty"`
	DishID       uuid.UUID  `gorm:"not null" json:"dish_id"`
	Dish         menuModels.Dish       `json:"dish,omitempty"`
	QRCodeID     uuid.UUID  `gorm:"not null" json:"qr_code_id"`
	QRCode       qrcodeModels.QRCode     `json:"qr_code,omitempty"`
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
	QuestionID   uuid.UUID    `json:"question_id"`
	QuestionText string       `json:"question_text,omitempty"`
	QuestionType QuestionType `json:"question_type,omitempty"`
	Answer       any          `json:"answer"`
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

