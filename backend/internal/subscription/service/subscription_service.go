package subscriptionservice

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	organizationinterface "kyooar/internal/organization/interface"
	subscriptionconstants "kyooar/internal/subscription/constants"
	subscriptioninterface "kyooar/internal/subscription/interface"
	subscriptionmodel "kyooar/internal/subscription/model"
)

type subscriptionService struct {
	subscriptionRepo     subscriptioninterface.SubscriptionRepository
	planRepo            subscriptioninterface.SubscriptionPlanRepository
	organizationRepo      organizationinterface.OrganizationRepository
}

func NewSubscriptionService(
	subscriptionRepo subscriptioninterface.SubscriptionRepository,
	planRepo subscriptioninterface.SubscriptionPlanRepository,
	organizationRepo organizationinterface.OrganizationRepository,
) subscriptioninterface.SubscriptionService {
	return &subscriptionService{
		subscriptionRepo: subscriptionRepo,
		planRepo:        planRepo,
		organizationRepo:  organizationRepo,
	}
}

func (s *subscriptionService) GetUserSubscription(ctx context.Context, accountID uuid.UUID) (*subscriptionmodel.Subscription, error) {
	return s.subscriptionRepo.FindByAccountID(ctx, accountID)
}

func (s *subscriptionService) GetAvailablePlans(ctx context.Context) ([]subscriptionmodel.SubscriptionPlan, error) {
	return s.planRepo.FindAll(ctx)
}

func (s *subscriptionService) CanUserCreateOrganization(ctx context.Context, accountID uuid.UUID) (*subscriptioninterface.PermissionResponse, error) {
	subscription, err := s.subscriptionRepo.FindByAccountID(ctx, accountID)
	if err != nil {
		return &subscriptioninterface.PermissionResponse{
			CanCreate:          false,
			Reason:            "No active subscription found",
			SubscriptionStatus: "none",
		}, nil
	}

	if !subscription.IsActive() {
		return &subscriptioninterface.PermissionResponse{
			CanCreate:          false,
			Reason:            "Subscription is not active or has expired",
			SubscriptionStatus: string(subscription.Status),
		}, nil
	}

	organizations, err := s.organizationRepo.FindByAccountID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization count: %w", err)
	}

	currentCount := len(organizations)
	maxOrganizations := subscription.Plan.MaxOrganizations

	canCreate := subscription.CanAddOrganization(currentCount)
	
	var reason string
	if maxOrganizations == -1 {
		reason = "Unlimited organizations allowed"
	} else {
		reason = fmt.Sprintf("Current: %d/%d organizations", currentCount, maxOrganizations)
	}

	return &subscriptioninterface.PermissionResponse{
		CanCreate:          canCreate,
		Reason:            reason,
		CurrentCount:      currentCount,
		MaxAllowed:        maxOrganizations,
		SubscriptionStatus: string(subscription.Status),
	}, nil
}

func (s *subscriptionService) CreateSubscription(ctx context.Context, req *subscriptioninterface.CreateSubscriptionRequest) (*subscriptionmodel.Subscription, error) {
	plan, err := s.planRepo.FindByID(ctx, req.PlanID)
	if err != nil {
		return nil, fmt.Errorf("plan not found: %w", err)
	}

	now := time.Now()
	subscription := &subscriptionmodel.Subscription{
		AccountID:          req.AccountID,
		PlanID:             req.PlanID,
		Status:             subscriptionconstants.SubscriptionActive,
		Plan:               *plan,
		CurrentPeriodStart: now,
		CurrentPeriodEnd:   now.AddDate(0, 1, 0),
	}

	err = s.subscriptionRepo.Create(ctx, subscription)
	if err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	return subscription, nil
}

func (s *subscriptionService) GetAllPlans(ctx context.Context) ([]subscriptionmodel.SubscriptionPlan, error) {
	return s.planRepo.FindAllIncludingHidden(ctx)
}

func (s *subscriptionService) AssignCustomPlan(ctx context.Context, accountID uuid.UUID, planCode string) error {
	plan, err := s.planRepo.FindByCode(ctx, planCode)
	if err != nil {
		return fmt.Errorf("plan not found: %w", err)
	}

	existingSubscription, err := s.subscriptionRepo.FindByAccountID(ctx, accountID)
	if err == nil && existingSubscription != nil {
		existingSubscription.PlanID = plan.ID
		existingSubscription.Plan = *plan
		existingSubscription.Status = subscriptionconstants.SubscriptionActive
		return s.subscriptionRepo.Update(ctx, existingSubscription)
	}

	now := time.Now()
	subscription := &subscriptionmodel.Subscription{
		AccountID:          accountID,
		PlanID:             plan.ID,
		Status:             subscriptionconstants.SubscriptionActive,
		Plan:               *plan,
		CurrentPeriodStart: now,
		CurrentPeriodEnd:   now.AddDate(0, 1, 0),
	}

	return s.subscriptionRepo.Create(ctx, subscription)
}

func (s *subscriptionService) CancelSubscription(ctx context.Context, accountID uuid.UUID) error {
	subscription, err := s.subscriptionRepo.FindByAccountID(ctx, accountID)
	if err != nil {
		return fmt.Errorf("subscription not found: %w", err)
	}

	subscription.Status = subscriptionconstants.SubscriptionCanceled
	return s.subscriptionRepo.Update(ctx, subscription)
}