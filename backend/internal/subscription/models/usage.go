package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	sharedModels "github.com/lecritique/api/internal/shared/models"
)

// SubscriptionUsage tracks usage metrics for a subscription billing period
type SubscriptionUsage struct {
	sharedModels.BaseModel
	SubscriptionID      uuid.UUID    `gorm:"not null;index" json:"subscription_id"`
	Subscription        Subscription `json:"subscription,omitempty"`
	PeriodStart         time.Time    `gorm:"not null;index" json:"period_start"`
	PeriodEnd           time.Time    `gorm:"not null;index" json:"period_end"`
	FeedbacksCount      int          `gorm:"default:0" json:"feedbacks_count"`
	RestaurantsCount    int          `gorm:"default:0" json:"restaurants_count"`
	LocationsCount      int          `gorm:"default:0" json:"locations_count"`
	QRCodesCount        int          `gorm:"default:0" json:"qr_codes_count"`
	TeamMembersCount    int          `gorm:"default:0" json:"team_members_count"`
	LastUpdatedAt       time.Time    `json:"last_updated_at"`
}

// UsageEvent tracks individual usage events for auditing
type UsageEvent struct {
	sharedModels.BaseModel
	SubscriptionID uuid.UUID  `gorm:"not null;index" json:"subscription_id"`
	EventType      string     `gorm:"not null" json:"event_type"`
	ResourceType   string     `gorm:"not null" json:"resource_type"`
	ResourceID     uuid.UUID  `json:"resource_id"`
	Metadata       string     `gorm:"type:jsonb" json:"metadata"`
	CreatedAt      time.Time  `json:"created_at"`
}

// Event types
const (
	EventTypeCreate = "create"
	EventTypeDelete = "delete"
	EventTypeUpdate = "update"
)

// Resource types
const (
	ResourceTypeFeedback   = "feedback"
	ResourceTypeRestaurant = "restaurant"
	ResourceTypeLocation   = "location"
	ResourceTypeQRCode     = "qr_code"
	ResourceTypeTeamMember = "team_member"
)

// GetCurrentUsage returns the usage for the current billing period
func (s *Subscription) GetCurrentUsage() *SubscriptionUsage {
	// This would be implemented to fetch from database
	// For now, returning a placeholder
	return &SubscriptionUsage{
		SubscriptionID:   s.ID,
		PeriodStart:      s.CurrentPeriodStart,
		PeriodEnd:        s.CurrentPeriodEnd,
		FeedbacksCount:   0,
		RestaurantsCount: 0,
		LocationsCount:   0,
		QRCodesCount:     0,
		TeamMembersCount: 0,
	}
}

// CanAddResource checks if a resource can be added based on plan limits
func (u *SubscriptionUsage) CanAddResource(resourceType string, plan PlanFeatures) (bool, string) {
	var limitKey string
	var currentUsage int64
	
	switch resourceType {
	case ResourceTypeFeedback:
		limitKey = LimitFeedbacksPerMonth
		currentUsage = int64(u.FeedbacksCount)
	case ResourceTypeRestaurant:
		limitKey = LimitRestaurants
		currentUsage = int64(u.RestaurantsCount)
	case ResourceTypeLocation:
		// Locations are no longer limited separately
		return true, ""
	case ResourceTypeQRCode:
		limitKey = LimitQRCodes
		currentUsage = int64(u.QRCodesCount)
	case ResourceTypeTeamMember:
		limitKey = LimitTeamMembers
		currentUsage = int64(u.TeamMembersCount)
	default:
		return true, ""
	}
	
	limit := plan.GetLimit(limitKey)
	if limit == -1 {
		return true, ""
	}
	
	if currentUsage >= limit {
		def, _ := GetFeatureDefinition(limitKey)
		return false, fmt.Sprintf("%s limit reached", def.DisplayName)
	}
	
	return true, ""
}