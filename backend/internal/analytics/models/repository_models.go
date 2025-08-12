package models

import (
	"time"
	"github.com/google/uuid"
	feedbackModels "kyooar/internal/feedback/models"
)

type FeedbackCounts struct {
	Total        int64 `gorm:"column:total"`
	Today        int64 `gorm:"column:today"`
	Yesterday    int64 `gorm:"column:yesterday"`
	Recent30Days int64 `gorm:"column:recent_30_days"`
}

type QRCodeMetrics struct {
	TotalQRCodes int64 `gorm:"column:total_qr_codes"`
	ActiveCount  int64 `gorm:"column:active_count"`
	TotalScans   int64 `gorm:"column:total_scans"`
	ScansToday   int64 `gorm:"column:scans_today"`
}

type FeedbackWithQRCode struct {
	FeedbackID         uuid.UUID                   `gorm:"column:feedback_id"`
	FeedbackCreatedAt  time.Time                   `gorm:"column:feedback_created_at"`
	QRCodeID           uuid.UUID                   `gorm:"column:qr_code_id"`
	DeviceInfo         feedbackModels.DeviceInfo   `gorm:"column:device_info;type:jsonb"`
	QRLastScannedAt    *time.Time                  `gorm:"column:qr_last_scanned_at"`
}

type QRCodePerformanceData struct {
	ID             uuid.UUID  `gorm:"column:id"`
	Label          string     `gorm:"column:label"`
	OrganizationID uuid.UUID  `gorm:"column:organization_id"`
	ScansCount     int64      `gorm:"column:scans_count"`
	LastScannedAt  *time.Time `gorm:"column:last_scanned_at"`
	IsActive       bool       `gorm:"column:is_active"`
	Location       *string    `gorm:"column:location"`
	FeedbackCount  int64      `gorm:"column:feedback_count"`
}

type ChartFilters struct {
	DateFrom  *time.Time
	DateTo    *time.Time
	ProductID *uuid.UUID
}

type ChartDataResult struct {
	OrganizationID uuid.UUID
	TotalResponses int64
	QuestionData   interface{}
}

type TimeSeriesDataPoint struct {
	Date          time.Time `gorm:"column:date"`
	Count         int64     `gorm:"column:count"`
	AverageRating float64   `gorm:"column:average_rating"`
}

type ProductMetrics struct {
	AverageRating float64 `gorm:"column:average_rating"`
	FeedbackCount int64   `gorm:"column:feedback_count"`
}