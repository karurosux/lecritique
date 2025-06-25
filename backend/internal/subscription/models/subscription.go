package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	authModels "github.com/lecritique/api/internal/auth/models"
	sharedModels "github.com/lecritique/api/internal/shared/models"
)

type Subscription struct {
	sharedModels.BaseModel
	AccountID          uuid.UUID          `gorm:"not null" json:"account_id"`
	Account            authModels.Account `json:"account,omitempty"`
	PlanID             uuid.UUID          `gorm:"not null" json:"plan_id"`
	Plan               SubscriptionPlan   `json:"plan,omitempty"`
	Status             SubscriptionStatus `gorm:"not null" json:"status"`
	CurrentPeriodStart time.Time          `json:"current_period_start"`
	CurrentPeriodEnd   time.Time          `json:"current_period_end"`
	CancelAt           *time.Time         `json:"cancel_at"`
	CancelledAt        *time.Time         `json:"cancelled_at"`
	StripeCustomerID   string             `json:"-"`
	StripeSubscriptionID string           `json:"-"`
}

type SubscriptionStatus string

const (
	SubscriptionActive   SubscriptionStatus = "active"
	SubscriptionPending  SubscriptionStatus = "pending"
	SubscriptionCanceled SubscriptionStatus = "canceled"
	SubscriptionExpired  SubscriptionStatus = "expired"
)

type SubscriptionPlan struct {
	sharedModels.BaseModel
	Name        string      `gorm:"not null" json:"name"`
	Code        string      `gorm:"uniqueIndex;not null" json:"code"`
	Description string      `json:"description"`
	Price       float64     `gorm:"not null" json:"price"`
	Currency    string      `gorm:"default:'USD'" json:"currency"`
	Interval    string      `gorm:"default:'month'" json:"interval"`
	Features    PlanFeatures `gorm:"type:jsonb" json:"features"`
	IsActive    bool        `gorm:"default:true" json:"is_active"`
	IsVisible   bool        `gorm:"default:true" json:"is_visible"`
	TrialDays   int         `gorm:"default:0" json:"trial_days"`
	StripePriceID string    `json:"-"`
}

// PlanFeatures uses generic maps for flexibility
type PlanFeatures struct {
	Limits map[string]int64        `json:"limits"`
	Flags  map[string]bool         `json:"flags"`
	Custom map[string]interface{}  `json:"custom,omitempty"`
}

// Helper methods for type-safe access
func (pf PlanFeatures) GetLimit(key string) int64 {
	if pf.Limits == nil {
		return 0
	}
	return pf.Limits[key]
}

func (pf PlanFeatures) GetFlag(key string) bool {
	if pf.Flags == nil {
		return false
	}
	return pf.Flags[key]
}

func (pf PlanFeatures) IsUnlimited(key string) bool {
	return pf.GetLimit(key) == -1
}

// Common limit keys as constants for type safety
const (
	LimitRestaurants          = "max_restaurants"
	LimitLocationsPerRestaurant = "max_locations_per_restaurant"
	LimitQRCodesPerLocation   = "max_qr_codes_per_location"
	LimitFeedbacksPerMonth    = "max_feedbacks_per_month"
	LimitTeamMembers          = "max_team_members"
	LimitStorageGB            = "max_storage_gb"
	LimitAPICallsPerHour      = "max_api_calls_per_hour"
)

// Common feature flags
const (
	FlagAdvancedAnalytics = "advanced_analytics"
	FlagCustomBranding    = "custom_branding"
	FlagAPIAccess         = "api_access"
	FlagPrioritySupport   = "priority_support"
	FlagWhiteLabel        = "white_label"
	FlagCustomDomain      = "custom_domain"
)

// GORM Scanner/Valuer interfaces for JSONB
func (pf PlanFeatures) Value() (driver.Value, error) {
	return json.Marshal(pf)
}

func (pf *PlanFeatures) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return json.Unmarshal([]byte("{}"), pf)
	}
	return json.Unmarshal(bytes, pf)
}

func (s *Subscription) IsActive() bool {
	return s.Status == SubscriptionActive && time.Now().Before(s.CurrentPeriodEnd)
}

func (s *Subscription) CanAddRestaurant(currentCount int) bool {
	limit := s.Plan.Features.GetLimit(LimitRestaurants)
	if limit == -1 {
		return true
	}
	return int64(currentCount) < limit
}
