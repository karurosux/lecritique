package models

import (
	"time"

	"github.com/google/uuid"
)

type QRCode struct {
	BaseModel
	RestaurantID uuid.UUID   `gorm:"not null" json:"restaurant_id"`
	Restaurant   Restaurant  `json:"restaurant,omitempty"`
	LocationID   *uuid.UUID  `json:"location_id"`
	Location     *Location   `json:"location,omitempty"`
	Code         string      `gorm:"uniqueIndex;not null" json:"code"`
	Label        string      `json:"label"` // e.g., "Table 1", "Entrance", etc.
	Type         QRCodeType  `gorm:"not null" json:"type"`
	IsActive     bool        `gorm:"default:true" json:"is_active"`
	ScansCount   int         `gorm:"default:0" json:"scans_count"`
	LastScannedAt *time.Time `json:"last_scanned_at"`
	ExpiresAt    *time.Time  `json:"expires_at"`
}

type QRCodeType string

const (
	QRCodeTypeTable    QRCodeType = "table"
	QRCodeTypeLocation QRCodeType = "location"
	QRCodeTypeTakeaway QRCodeType = "takeaway"
	QRCodeTypeDelivery QRCodeType = "delivery"
	QRCodeTypeGeneral  QRCodeType = "general"
)

func (q *QRCode) IsValid() bool {
	return q.IsActive && (q.ExpiresAt == nil || time.Now().Before(*q.ExpiresAt))
}
