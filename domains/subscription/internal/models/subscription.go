package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/lecritique/api/shared/models"
)

type Subscription struct {
	models.BaseModel
	AccountID          uuid.UUID          `gorm:"not null" json:"account_id"`
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
	models.BaseModel
	Name        string      `gorm:"not null" json:"name"`
	Code        string      `gorm:"uniqueIndex;not null" json:"code"`
	Description string      `json:"description"`
	Price       float64     `gorm:"not null" json:"price"`
	Currency    string      `gorm:"default:'USD'" json:"currency"`
	Interval    string      `gorm:"default:'month'" json:"interval"`
	Features    PlanFeatures `gorm:"type:jsonb" json:"features"`
	IsActive    bool        `gorm:"default:true" json:"is_active"`
	StripePriceID string    `json:"-"`
}

type PlanFeatures struct {
	MaxRestaurants          int  `json:"max_restaurants"`
	MaxLocationsPerRestaurant int  `json:"max_locations_per_restaurant"`
	MaxQRCodesPerLocation   int  `json:"max_qr_codes_per_location"`
	MaxFeedbacksPerMonth    int  `json:"max_feedbacks_per_month"`
	MaxTeamMembers          int  `json:"max_team_members"`
	AdvancedAnalytics       bool `json:"advanced_analytics"`
	CustomBranding          bool `json:"custom_branding"`
	APIAccess               bool `json:"api_access"`
	PrioritySupport         bool `json:"priority_support"`
}

func (s *Subscription) IsActive() bool {
	return s.Status == SubscriptionActive && time.Now().Before(s.CurrentPeriodEnd)
}

func (s *Subscription) CanAddRestaurant(currentCount int) bool {
	if s.Plan.Features.MaxRestaurants == -1 {
		return true
	}
	return currentCount < s.Plan.Features.MaxRestaurants
}

// GORM Scanner/Valuer interface for PlanFeatures JSONB
func (p PlanFeatures) Value() (driver.Value, error) { return json.Marshal(p) }
func (p *PlanFeatures) Scan(value interface{}) error { 
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return json.Unmarshal([]byte("{}"), p)
	}
	return json.Unmarshal(bytes, p) 
}