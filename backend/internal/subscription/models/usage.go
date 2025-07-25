package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	sharedModels "kyooar/internal/shared/models"
)

// SubscriptionUsage tracks usage metrics for a subscription billing period
type SubscriptionUsage struct {
	sharedModels.BaseModel
	SubscriptionID      uuid.UUID    `gorm:"not null;index" json:"subscription_id"`
	Subscription        Subscription `json:"subscription,omitempty"`
	PeriodStart         time.Time    `gorm:"not null;index" json:"period_start"`
	PeriodEnd           time.Time    `gorm:"not null;index" json:"period_end"`
	FeedbacksCount      int          `gorm:"default:0" json:"feedbacks_count"`
	OrganizationsCount    int          `gorm:"default:0" json:"organizations_count"`
	LocationsCount      int          `gorm:"default:0" json:"locations_count"`
	QRCodesCount        int          `gorm:"default:0" json:"qr_codes_count"`
	TeamMembersCount    int          `gorm:"default:0" json:"team_members_count"`
	LastUpdatedAt       time.Time    `json:"last_updated_at"`
}

// TableName overrides the table name used by this model
func (SubscriptionUsage) TableName() string {
	return "subscription_usage"
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
	ResourceTypeOrganization = "organization"
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
		OrganizationsCount: 0,
		LocationsCount:   0,
		QRCodesCount:     0,
		TeamMembersCount: 0,
	}
}

// CanAddResource checks if a resource can be added based on plan limits
func (u *SubscriptionUsage) CanAddResource(resourceType string, plan *SubscriptionPlan) (bool, string) {
	var limit int
	var currentUsage int
	
	switch resourceType {
	case ResourceTypeFeedback:
		limit = plan.MaxFeedbacksPerMonth
		currentUsage = u.FeedbacksCount
	case ResourceTypeOrganization:
		limit = plan.MaxOrganizations
		currentUsage = u.OrganizationsCount
	case ResourceTypeLocation:
		// Locations are no longer limited separately
		return true, ""
	case ResourceTypeQRCode:
		limit = plan.MaxQRCodes
		currentUsage = u.QRCodesCount
	case ResourceTypeTeamMember:
		limit = plan.MaxTeamMembers
		currentUsage = u.TeamMembersCount
	default:
		return true, ""
	}
	if limit == -1 {
		return true, ""
	}
	
	if currentUsage >= limit {
		var resourceName string
		switch resourceType {
		case ResourceTypeFeedback:
			resourceName = "Monthly feedback"
		case ResourceTypeOrganization:
			resourceName = "Organization"
		case ResourceTypeQRCode:
			resourceName = "QR code"
		case ResourceTypeTeamMember:
			resourceName = "Team member"
		default:
			resourceName = "Resource"
		}
		return false, fmt.Sprintf("%s limit reached (%d/%d)", resourceName, currentUsage, limit)
	}
	
	return true, ""
}
