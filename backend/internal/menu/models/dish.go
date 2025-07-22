package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"kyooar/internal/shared/models"
)

type Product struct {
	models.BaseModel
	OrganizationID  uuid.UUID      `gorm:"not null" json:"organization_id"`
	Name          string         `gorm:"not null" json:"name"`
	Description   string         `json:"description"`
	Category      string         `json:"category"`
	Price         float64        `json:"price"`
	Currency      string         `gorm:"default:'USD'" json:"currency"`
	Image         string         `json:"image"`
	Tags          pq.StringArray `gorm:"type:text[]" json:"tags"`
	IsAvailable   bool           `gorm:"default:true" json:"is_available"`
	IsActive      bool           `gorm:"default:true" json:"is_active"`
	DisplayOrder  int            `gorm:"default:0" json:"display_order"`
}
