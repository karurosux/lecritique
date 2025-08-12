package feedbackmodel

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	menuModels "kyooar/internal/menu/models"
	qrcodeModels "kyooar/internal/qrcode/models"
	organizationModels "kyooar/internal/organization/models"
	sharedModels "kyooar/internal/shared/models"
)

type Feedback struct {
	sharedModels.BaseModel
	OrganizationID uuid.UUID                           `gorm:"not null" json:"organization_id"`
	Organization   organizationModels.Organization     `json:"organization,omitempty"`
	ProductID      uuid.UUID                           `gorm:"not null" json:"product_id"`
	Product        menuModels.Product                  `json:"product,omitempty"`
	QRCodeID       uuid.UUID                           `gorm:"not null" json:"qr_code_id"`
	QRCode         qrcodeModels.QRCode                 `json:"qr_code,omitempty"`
	CustomerName   string                              `json:"customer_name"`
	CustomerEmail  string                              `json:"customer_email"`
	CustomerPhone  string                              `json:"customer_phone"`
	OverallRating  int                                 `json:"overall_rating"`
	Responses      Responses                           `gorm:"type:jsonb" json:"responses"`
	DeviceInfo     DeviceInfo                          `gorm:"type:jsonb" json:"device_info"`
	IsComplete     bool                                `gorm:"default:true" json:"is_complete"`
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

type FeedbackFilter struct {
	Search     string     `json:"search,omitempty"`
	RatingMin  *int       `json:"rating_min,omitempty"`
	RatingMax  *int       `json:"rating_max,omitempty"`
	DateFrom   *time.Time `json:"date_from,omitempty"`
	DateTo     *time.Time `json:"date_to,omitempty"`
	ProductID  *uuid.UUID `json:"product_id,omitempty"`
	IsComplete *bool      `json:"is_complete,omitempty"`
}

type FeedbackStats struct {
	TotalFeedbacks    int64   `json:"total_feedbacks"`
	AverageRating     float64 `json:"average_rating"`
	FeedbacksToday    int64   `json:"feedbacks_today"`
	FeedbacksThisWeek int64   `json:"feedbacks_this_week"`
}