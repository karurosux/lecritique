package models

import (
	"time"

	"github.com/google/uuid"
	organizationModels "kyooar/internal/organization/models"
	sharedModels "kyooar/internal/shared/models"
)

type QRCode struct {
	sharedModels.BaseModel
	OrganizationID uuid.UUID   `gorm:"not null" json:"organization_id"`
	Organization   organizationModels.Organization  `json:"organization,omitempty"`
	Location     *string     `json:"location"` // Free text location description
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
