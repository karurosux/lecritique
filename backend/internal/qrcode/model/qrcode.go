package qrcodemodel

import (
	"time"

	"github.com/google/uuid"
	organizationmodel "kyooar/internal/organization/model"
	qrcodeconstants "kyooar/internal/qrcode/constants"
	sharedModels "kyooar/internal/shared/models"
)

type QRCodeType = qrcodeconstants.QRCodeType

type QRCode struct {
	sharedModels.BaseModel
	OrganizationID uuid.UUID   `gorm:"not null" json:"organization_id"`
	Organization   organizationmodel.Organization  `json:"organization,omitempty"`
	Location     *string     `json:"location"`
	Code         string      `gorm:"uniqueIndex;not null" json:"code"`
	Label        string      `json:"label"`
	Type         QRCodeType  `gorm:"not null" json:"type"`
	IsActive     bool        `gorm:"default:true" json:"is_active"`
	ScansCount   int         `gorm:"default:0" json:"scans_count"`
	LastScannedAt *time.Time `json:"last_scanned_at"`
	ExpiresAt    *time.Time  `json:"expires_at"`
}

func (q *QRCode) IsValid() bool {
	return q.IsActive && (q.ExpiresAt == nil || time.Now().Before(*q.ExpiresAt))
}