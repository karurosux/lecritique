package subscriptionservice

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	subscriptioninterface "kyooar/internal/subscription/interface"
	subscriptionmodel "kyooar/internal/subscription/model"
)

type usageService struct {
	usageRepo        subscriptioninterface.UsageRepository
	subscriptionRepo subscriptioninterface.SubscriptionRepository
}

func NewUsageService(
	usageRepo subscriptioninterface.UsageRepository,
	subscriptionRepo subscriptioninterface.SubscriptionRepository,
) subscriptioninterface.UsageService {
	return &usageService{
		usageRepo:        usageRepo,
		subscriptionRepo: subscriptionRepo,
	}
}

func (s *usageService) TrackUsage(ctx context.Context, subscriptionID uuid.UUID, resourceType string, delta int) error {
	subscription, err := s.subscriptionRepo.FindByID(ctx, subscriptionID)
	if err != nil {
		return fmt.Errorf("subscription not found: %w", err)
	}

	usage, err := s.usageRepo.FindBySubscriptionAndPeriod(ctx, subscriptionID, subscription.CurrentPeriodStart, subscription.CurrentPeriodEnd)
	if err != nil {
		usage = &subscriptionmodel.SubscriptionUsage{
			SubscriptionID: subscriptionID,
			PeriodStart:    subscription.CurrentPeriodStart,
			PeriodEnd:      subscription.CurrentPeriodEnd,
		}
		if err := s.usageRepo.Create(ctx, usage); err != nil {
			return fmt.Errorf("failed to create usage record: %w", err)
		}
	}

	switch resourceType {
	case subscriptionmodel.ResourceTypeFeedback:
		usage.FeedbacksCount += delta
	case subscriptionmodel.ResourceTypeOrganization:
		usage.OrganizationsCount += delta
	case subscriptionmodel.ResourceTypeLocation:
		usage.LocationsCount += delta
	case subscriptionmodel.ResourceTypeQRCode:
		usage.QRCodesCount += delta
	case subscriptionmodel.ResourceTypeTeamMember:
		usage.TeamMembersCount += delta
	default:
		return fmt.Errorf("unknown resource type: %s", resourceType)
	}

	usage.LastUpdatedAt = time.Now()
	return s.usageRepo.Update(ctx, usage)
}

func (s *usageService) RecordUsageEvent(ctx context.Context, event *subscriptionmodel.UsageEvent) error {
	event.CreatedAt = time.Now()
	return s.usageRepo.CreateEvent(ctx, event)
}

func (s *usageService) CanAddResource(ctx context.Context, subscriptionID uuid.UUID, resourceType string) (bool, string, error) {
	subscription, err := s.subscriptionRepo.FindByID(ctx, subscriptionID)
	if err != nil {
		return false, "Subscription not found", err
	}

	if !subscription.IsActive() {
		return false, "Subscription is not active", nil
	}

	usage, err := s.GetCurrentUsage(ctx, subscriptionID)
	if err != nil {
		return false, "Failed to get usage data", err
	}

	canAdd, reason := usage.CanAddResource(resourceType, &subscription.Plan)
	return canAdd, reason, nil
}

func (s *usageService) GetCurrentUsage(ctx context.Context, subscriptionID uuid.UUID) (*subscriptionmodel.SubscriptionUsage, error) {
	subscription, err := s.subscriptionRepo.FindByID(ctx, subscriptionID)
	if err != nil {
		return nil, fmt.Errorf("subscription not found: %w", err)
	}

	usage, err := s.usageRepo.FindBySubscriptionAndPeriod(ctx, subscriptionID, subscription.CurrentPeriodStart, subscription.CurrentPeriodEnd)
	if err != nil {
		return &subscriptionmodel.SubscriptionUsage{
			SubscriptionID:   subscriptionID,
			PeriodStart:      subscription.CurrentPeriodStart,
			PeriodEnd:        subscription.CurrentPeriodEnd,
			FeedbacksCount:   0,
			OrganizationsCount: 0,
			LocationsCount:   0,
			QRCodesCount:     0,
			TeamMembersCount: 0,
		}, nil
	}

	return usage, nil
}

func (s *usageService) GetUsageForPeriod(ctx context.Context, subscriptionID uuid.UUID, start, end time.Time) (*subscriptionmodel.SubscriptionUsage, error) {
	return s.usageRepo.FindBySubscriptionAndPeriod(ctx, subscriptionID, start, end)
}

func (s *usageService) InitializeUsagePeriod(ctx context.Context, subscriptionID uuid.UUID, periodStart, periodEnd time.Time) error {
	existing, _ := s.usageRepo.FindBySubscriptionAndPeriod(ctx, subscriptionID, periodStart, periodEnd)
	if existing != nil {
		return nil
	}

	usage := &subscriptionmodel.SubscriptionUsage{
		SubscriptionID: subscriptionID,
		PeriodStart:    periodStart,
		PeriodEnd:      periodEnd,
		LastUpdatedAt:  time.Now(),
	}

	return s.usageRepo.Create(ctx, usage)
}

func (s *usageService) ResetMonthlyUsage(ctx context.Context) error {
	return s.usageRepo.ResetMonthlyUsage(ctx)
}