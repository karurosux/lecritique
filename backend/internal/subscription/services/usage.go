package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lecritique/api/internal/subscription/models"
	subscriptionRepos "github.com/lecritique/api/internal/subscription/repositories"
)

type UsageService interface {
	// Track usage
	TrackUsage(ctx context.Context, subscriptionID uuid.UUID, resourceType string, delta int) error
	RecordUsageEvent(ctx context.Context, event *models.UsageEvent) error
	
	// Check limits
	CanAddResource(ctx context.Context, subscriptionID uuid.UUID, resourceType string) (bool, string, error)
	GetCurrentUsage(ctx context.Context, subscriptionID uuid.UUID) (*models.SubscriptionUsage, error)
	GetUsageForPeriod(ctx context.Context, subscriptionID uuid.UUID, start, end time.Time) (*models.SubscriptionUsage, error)
	
	// Period management
	InitializeUsagePeriod(ctx context.Context, subscriptionID uuid.UUID, periodStart, periodEnd time.Time) error
	ResetMonthlyUsage(ctx context.Context) error
}

type usageService struct {
	usageRepo        subscriptionRepos.UsageRepository
	subscriptionRepo subscriptionRepos.SubscriptionRepository
}

func NewUsageService(
	usageRepo subscriptionRepos.UsageRepository,
	subscriptionRepo subscriptionRepos.SubscriptionRepository,
) UsageService {
	return &usageService{
		usageRepo:        usageRepo,
		subscriptionRepo: subscriptionRepo,
	}
}

func (s *usageService) TrackUsage(ctx context.Context, subscriptionID uuid.UUID, resourceType string, delta int) error {
	// Get current subscription to verify it exists and get period
	subscription, err := s.subscriptionRepo.FindByID(ctx, subscriptionID)
	if err != nil {
		return fmt.Errorf("subscription not found: %w", err)
	}

	// Get or create usage record for current period
	usage, err := s.usageRepo.FindBySubscriptionAndPeriod(ctx, subscriptionID, subscription.CurrentPeriodStart, subscription.CurrentPeriodEnd)
	if err != nil {
		// Create new usage record
		usage = &models.SubscriptionUsage{
			SubscriptionID: subscriptionID,
			PeriodStart:    subscription.CurrentPeriodStart,
			PeriodEnd:      subscription.CurrentPeriodEnd,
		}
		if err := s.usageRepo.Create(ctx, usage); err != nil {
			return fmt.Errorf("failed to create usage record: %w", err)
		}
	}

	// Update the appropriate counter
	switch resourceType {
	case models.ResourceTypeFeedback:
		usage.FeedbacksCount += delta
	case models.ResourceTypeRestaurant:
		usage.RestaurantsCount += delta
	case models.ResourceTypeLocation:
		usage.LocationsCount += delta
	case models.ResourceTypeQRCode:
		usage.QRCodesCount += delta
	case models.ResourceTypeTeamMember:
		usage.TeamMembersCount += delta
	default:
		return fmt.Errorf("unknown resource type: %s", resourceType)
	}

	usage.LastUpdatedAt = time.Now()
	return s.usageRepo.Update(ctx, usage)
}

func (s *usageService) RecordUsageEvent(ctx context.Context, event *models.UsageEvent) error {
	event.CreatedAt = time.Now()
	return s.usageRepo.CreateEvent(ctx, event)
}

func (s *usageService) CanAddResource(ctx context.Context, subscriptionID uuid.UUID, resourceType string) (bool, string, error) {
	// Get subscription with plan
	subscription, err := s.subscriptionRepo.FindByID(ctx, subscriptionID)
	if err != nil {
		return false, "Subscription not found", err
	}

	if !subscription.IsActive() {
		return false, "Subscription is not active", nil
	}

	// Get current usage
	usage, err := s.GetCurrentUsage(ctx, subscriptionID)
	if err != nil {
		return false, "Failed to get usage data", err
	}

	// Check against plan limits
	canAdd, reason := usage.CanAddResource(resourceType, subscription.Plan.Features)
	return canAdd, reason, nil
}

func (s *usageService) GetCurrentUsage(ctx context.Context, subscriptionID uuid.UUID) (*models.SubscriptionUsage, error) {
	subscription, err := s.subscriptionRepo.FindByID(ctx, subscriptionID)
	if err != nil {
		return nil, fmt.Errorf("subscription not found: %w", err)
	}

	usage, err := s.usageRepo.FindBySubscriptionAndPeriod(ctx, subscriptionID, subscription.CurrentPeriodStart, subscription.CurrentPeriodEnd)
	if err != nil {
		// Return empty usage if not found
		return &models.SubscriptionUsage{
			SubscriptionID:   subscriptionID,
			PeriodStart:      subscription.CurrentPeriodStart,
			PeriodEnd:        subscription.CurrentPeriodEnd,
			FeedbacksCount:   0,
			RestaurantsCount: 0,
			LocationsCount:   0,
			QRCodesCount:     0,
			TeamMembersCount: 0,
		}, nil
	}

	return usage, nil
}

func (s *usageService) GetUsageForPeriod(ctx context.Context, subscriptionID uuid.UUID, start, end time.Time) (*models.SubscriptionUsage, error) {
	return s.usageRepo.FindBySubscriptionAndPeriod(ctx, subscriptionID, start, end)
}

func (s *usageService) InitializeUsagePeriod(ctx context.Context, subscriptionID uuid.UUID, periodStart, periodEnd time.Time) error {
	// Check if usage already exists for this period
	existing, _ := s.usageRepo.FindBySubscriptionAndPeriod(ctx, subscriptionID, periodStart, periodEnd)
	if existing != nil {
		return nil // Already initialized
	}

	// Create new usage record
	usage := &models.SubscriptionUsage{
		SubscriptionID: subscriptionID,
		PeriodStart:    periodStart,
		PeriodEnd:      periodEnd,
		LastUpdatedAt:  time.Now(),
	}

	return s.usageRepo.Create(ctx, usage)
}

func (s *usageService) ResetMonthlyUsage(ctx context.Context) error {
	// This would be called by a cron job to reset monthly counters
	// For now, we'll just return nil
	// In production, this would:
	// 1. Find all subscriptions where current period has ended
	// 2. Create new usage records for the new period
	// 3. Update subscription period dates
	return nil
}