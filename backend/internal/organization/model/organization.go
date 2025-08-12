package organizationmodel

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
	Name        string         `gorm:"not null" json:"name"`
	Description string         `json:"description"`
	Address     string         `json:"address"`
	Logo        string         `json:"logo"`
	Website     string         `json:"website"`
	Phone       string         `json:"phone"`
	Email       string         `json:"email"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	Settings    Settings       `gorm:"type:jsonb" json:"settings"`
}

type Settings struct {
	Language             string `json:"language"`
	Timezone             string `json:"timezone"`
	FeedbackNotification bool   `json:"feedback_notification"`
	LowRatingThreshold   int    `json:"low_rating_threshold"`
}

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
	Tags          pq.StringArray    `gorm:"type:text[]" json:"tags" swaggertype:"array,string"`
	IsAvailable   bool              `gorm:"default:true" json:"is_available"`
	IsActive      bool              `gorm:"default:true" json:"is_active"`
	DisplayOrder  int               `gorm:"default:0" json:"display_order"`
}
