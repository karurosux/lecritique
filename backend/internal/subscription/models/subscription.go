package models

import (
	"time"

	"github.com/google/uuid"
	authModels "lecritique/internal/auth/models"
	sharedModels "lecritique/internal/shared/models"
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
	
	// Limits (as columns)
	MaxOrganizations      int `gorm:"not null;default:1;check:max_organizations >= -1" json:"max_organizations"`
	MaxQRCodes         int `gorm:"column:max_qr_codes;not null;default:5;check:max_qr_codes >= -1" json:"max_qr_codes"`
	MaxFeedbacksPerMonth int `gorm:"column:max_feedbacks_per_month;not null;default:50;check:max_feedbacks_per_month >= -1" json:"max_feedbacks_per_month"`
	MaxTeamMembers     int `gorm:"column:max_team_members;not null;default:2;check:max_team_members >= -1" json:"max_team_members"`
	
	// Feature flags (as columns)
	HasBasicAnalytics    bool `gorm:"column:has_basic_analytics;not null;default:false" json:"has_basic_analytics"`
	HasAdvancedAnalytics bool `gorm:"column:has_advanced_analytics;not null;default:false" json:"has_advanced_analytics"`
	HasFeedbackExplorer  bool `gorm:"column:has_feedback_explorer;not null;default:false" json:"has_feedback_explorer"`
	HasCustomBranding    bool `gorm:"column:has_custom_branding;not null;default:false" json:"has_custom_branding"`
	HasPrioritySupport   bool `gorm:"column:has_priority_support;not null;default:false" json:"has_priority_support"`
	
	IsActive    bool        `gorm:"default:true" json:"is_active"`
	IsVisible   bool        `gorm:"default:true" json:"is_visible"`
	TrialDays   int         `gorm:"default:0" json:"trial_days"`
	StripePriceID string    `json:"-"`
}

// Helper methods for SubscriptionPlan
func (sp *SubscriptionPlan) GetLimit(key string) int {
	switch key {
	case LimitOrganizations:
		return sp.MaxOrganizations
	case LimitQRCodes:
		return sp.MaxQRCodes
	case LimitFeedbacksPerMonth:
		return sp.MaxFeedbacksPerMonth
	case LimitTeamMembers:
		return sp.MaxTeamMembers
	default:
		return 0
	}
}

func (sp *SubscriptionPlan) GetFlag(key string) bool {
	switch key {
	case FlagBasicAnalytics:
		return sp.HasBasicAnalytics
	case FlagAdvancedAnalytics:
		return sp.HasAdvancedAnalytics
	case FlagFeedbackExplorer:
		return sp.HasFeedbackExplorer
	case FlagCustomBranding:
		return sp.HasCustomBranding
	case FlagPrioritySupport:
		return sp.HasPrioritySupport
	default:
		return false
	}
}

func (sp *SubscriptionPlan) IsUnlimited(key string) bool {
	return sp.GetLimit(key) == -1
}

// Common limit keys as constants for type safety
const (
	LimitOrganizations       = "max_organizations"
	LimitQRCodes           = "max_qr_codes"
	LimitFeedbacksPerMonth = "max_feedbacks_per_month"
	LimitTeamMembers       = "max_team_members"
)

// Common feature flags
const (
	FlagBasicAnalytics    = "basic_analytics"
	FlagAdvancedAnalytics = "advanced_analytics"
	FlagFeedbackExplorer  = "feedback_explorer"
	FlagCustomBranding    = "custom_branding"
	FlagPrioritySupport   = "priority_support"
)


func (s *Subscription) IsActive() bool {
	return s.Status == SubscriptionActive && time.Now().Before(s.CurrentPeriodEnd)
}

func (s *Subscription) CanAddOrganization(currentCount int) bool {
	limit := s.Plan.GetLimit(LimitOrganizations)
	if limit == -1 {
		return true
	}
	return currentCount < limit
}
